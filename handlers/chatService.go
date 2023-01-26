package handlers

import (
	"encoding/json"
	"feedbacks/db"
	"feedbacks/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var ch = make(chan int64, 100)

type WsServer struct {
	clients    map[*Chat]bool
	register   chan *Chat
	unregister chan *Chat
	receiver   chan models.ReceiveMessage
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Chat]bool),
		register:   make(chan *Chat),
		unregister: make(chan *Chat),
		receiver:   make(chan models.ReceiveMessage),
	}
}

func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case receive := <-server.receiver:
			server.ReceiveMessage(receive)
		}
	}
}

func (client *Chat) readMessage() {
	defer func() {
		client.disconnect()
	}()
	var n = 1
	var message models.ReceiveMessage

	for {
		err := client.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if err.Error() == "websocket: close 1000 (normal)" {
					break
				}
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		if message.MessageType == "text" {
			client.wsServer.receiver <- models.ReceiveMessage{User: models.User{
				ID:      client.user.ID,
				Role:    client.user.Role,
				Name:    client.user.Name,
				Room_id: client.user.Room_id,
			}, Message: message.Message, MessageType: "message"}
		} else if message.MessageType == "file" {
			fileType := http.DetectContentType([]byte(message.Message))
			switch fileType {
			case "image/jpeg":
				fileType = ".jpg"
			case "image/png":
				fileType = ".png"
			default:
				log.Println("error while detecting file type")
				return
			}

			err = ioutil.WriteFile(filepath.Join("files/temp", "photo_"+strconv.Itoa(n)+fileType), []byte(message.Message), 0700)
			if err != nil {
				log.Println("error while creating photo in files/temp: ", err)
				return
			}

			client.wsServer.receiver <- models.ReceiveMessage{User: models.User{
				ID:      client.user.ID,
				Role:    client.user.Role,
				Name:    client.user.Name,
				Room_id: client.user.Room_id,
			}, Message: "files/temp/photo_" + strconv.Itoa(n) + fileType, MessageType: "file"}
			n++
		}
	}
}

func (client *Chat) writeMessage() {
	defer client.conn.Close()
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			resp, err := json.Marshal(message)
			if err != nil {
				log.Println("error while marshaling resp:", err)
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(resp)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (client *Chat) disconnect() {
	client.wsServer.unregister <- client
	close(client.send)
	client.conn.Close()
}

func (server *WsServer) registerClient(client *Chat) {
	server.clients[client] = true
	server.roomListr()
	if (client.user.Role == SUPERVISOR || client.user.Role == SUPPORT) && client.user.Room_id != 0 {
		ch <- client.user.Room_id
		chatHistory(client)
	} else if client.user.Role != SUPERVISOR && client.user.Role != SUPPORT {
		server.supportIsAvailable(client)
		chatHistory(client)
	}
}

func chatHistory(chat *Chat) {
	var message models.Message
	var chatHistory []models.Message
	value, err := db.Rdb.LRange(fmt.Sprint(chat.user.Room_id), 0, -1).Result()
	if err != nil {
		log.Println(err)
		return
	}
	switch chat.user.Role {
	case SUPERVISOR:
		for _, v := range value {
			err = json.Unmarshal([]byte(v), &chatHistory)
			log.Println(v)
			if err != nil {
				log.Printf("Error while marshaling: %v", err)
				return
			}
		}
		chat.send <- models.SendMessage{
			Type: "message",
			Data: chatHistory,
		}
	case SUPPORT:
		for _, v := range value {
			err = json.Unmarshal([]byte(v), &message)
			log.Println(v)
			if err != nil {
				log.Printf("Error while marshaling: %v", err)
				return
			}
			if message.Time > time.Now().Add(-1*time.Hour).Format("2006/02/01 15:04:05") {
				chatHistory = append(chatHistory, message)
			}
		}
		chat.send <- models.SendMessage{
			Type: "message",
			Data: chatHistory,
		}
	default:
		for _, v := range value {
			err = json.Unmarshal([]byte(v), &message)
			if err != nil {
				log.Printf("Error while marshaling: %v", err)
				return
			}
			if message.Time > time.Now().Add(-20*time.Minute).Format("2006/02/01 15:04:05") {
				chatHistory = append(chatHistory, message)
			}
		}
		chat.send <- models.SendMessage{
			Type: "message",
			Data: chatHistory,
		}
	}
}

func (server *WsServer) unregisterClient(client *Chat) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
	server.roomListr()
	if client.user.Role != SUPERVISOR && client.user.Role != SUPPORT {

		server.ReceiveMessage(models.ReceiveMessage{
			User: models.User{
				ID:      0,
				Room_id: client.user.Room_id,
				Name:    "server",
				Role:    0,
			},
			Message:     fmt.Sprintf("Клиент: %v, покинул чат.", client.user.Name),
			MessageType: "message",
		})
	}
}

func (server *WsServer) ReceiveMessage(data models.ReceiveMessage) {
	redisText := fmt.Sprintf(`{"id":%v, "role":%v, "sender":"%v", "message":"%v", "type":"%v", "time":"%v"}`, data.User.ID, data.User.Role, data.User.Name, data.Message, data.MessageType, time.Now().Format("2006/02/01 15:04:05"))
	_, err := db.Rdb.LPush(fmt.Sprint(data.User.Room_id), redisText).Result()
	if err != nil {
		log.Println(err)
		return
	}

	for client := range server.clients {
		if client.user.Room_id == data.User.Room_id {
			client.send <- models.SendMessage{
				Type: data.MessageType,
				Data: &models.Message{
					ID:      data.User.ID,
					Role:    data.User.Role,
					Sender:  data.User.Name,
					Message: data.Message,
					Time:    time.Now().Format("2006/02/01 15:04"),
				},
			}
		}
	}
}

func (server *WsServer) roomListr() {
	var room = &models.Rooms{
		Rooms_conf: nil,
		Max:        0,
	}
	for client := range server.clients {
		if client.user.Role != SUPERVISOR && client.user.Role != SUPPORT {
			room.Max++
			room.Rooms_conf = append(room.Rooms_conf, models.RoomsConf{
				ID:   client.user.Role,
				Name: client.user.Name,
				Time: time.Now().Format("2006/02/01 15:04"),
			})
		}
	}
	message := fmt.Sprintf(`%+v`, room)
	log.Println(string(message))
	for client := range server.clients {
		if client.user.Role == SUPERVISOR || client.user.Role == SUPPORT {
			client.send <- models.SendMessage{
				Type: "conf",
				Data: room,
			}
		}
	}
}

func (server *WsServer) supportIsAvailable(client *Chat) {
	var countUser int64 = 0
	var countSup int64 = 0

	for client := range server.clients {
		if client.user.Role != SUPERVISOR && client.user.Role != SUPPORT {
			countUser++
		} else {
			countSup++
		}
	}

	if countSup < countUser {
		client.send <- models.SendMessage{
			Type: "message",
			Data: &models.Message{
				ID:      0,
				Role:    0,
				Sender:  "server",
				Message: "Сейчас все консультанты занять, на ваш вопрос ответять позже...",
				Time:    time.Now().Format("2006/02/01 15:04"),
			},
		}
		go waitTime(5, client)
		return
	}
	client.send <- models.SendMessage{
		Type: "message",
		Data: &models.Message{
			ID:      0,
			Role:    0,
			Sender:  "server",
			Message: "Консультант скоро ответит на ваш воопрос.",
			Time:    time.Now().Format("2006/02/01 15:04"),
		},
	}
	go waitTime(2, client)
}
func waitTime(t int64, client *Chat) {
	t2 := time.NewTimer(time.Duration(t) * time.Minute)
	for {
		select {
		case <-t2.C:
			// SenderBot(client)
			return
		case c := <-ch:
			if client.user.Room_id == c {
				return
			}
		}
	}
}
