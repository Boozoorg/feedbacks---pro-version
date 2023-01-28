package handlers

import (
	"encoding/json"
	"feedbacks/db"
	"feedbacks/models"
	"fmt"
	"log"
	"strconv"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

var ch = make(chan string, 100)
var sender = make(chan string, 100)

type SocketsService struct {
	users map[models.User]bool
}

func NewWebsocketServer() *SocketsService {
	return &SocketsService{
		users: make(map[models.User]bool),
	}
}

func (socket *SocketsService) SocketsEvents(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		user, err := ParseUsers(s)
		if err != nil {
			log.Println("error while parsing user: ", err)
			return err
		}
		socket.users[*user] = true
		switch user.Role {
		case SUPERVISOR, SUPPORT:
		default:

		}

		socket.RoomListr(server)
		return nil
	})

	server.OnEvent("/", "join", func(s socketio.Conn, data map[string]string) {
		for _, v := range s.Rooms() {
			if v == "0" || v == s.ID() {
				continue
			}
			s.Leave(v)
		}

		if server.RoomLen("/", data["room_id"]) > 1 && data["room_id"] != "0" {
			if ok := server.BroadcastToRoom("/", s.ID(), "error", "Эта комнота уже занята другим саппортом."); !ok {
				log.Println("error while sending data to support")
				return
			}
			return
		}
		if ok := server.JoinRoom("/", data["room_id"], s); !ok {
			log.Println("error while connecting to supports room")
			return
		}

		socket.RoomListr(server)
		log.Println("user:", data["room_id"], ", join:", s.ID())
	})

	server.OnEvent("/", "send", func(s socketio.Conn, msg interface{}) {
		var resp models.ReceiveMessage
		var message = models.Message{}
		data, _ := json.Marshal(msg)
		json.Unmarshal(data, &resp)
		if resp.Room == "" {
			if ok := server.BroadcastToRoom("/", s.ID(), "message", message); !ok {
				log.Println("error while sending data to supports")
				return
			}
		}
		if ok := server.BroadcastToRoom("/", resp.Room, "message", message); !ok {
			log.Println("error while sending data to supports")
			return
		}
		// socket.ChatHistory(room, 0)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		user, err := ParseUsers(s)
		if err != nil {
			log.Println("error while parsing user: ", err)
			return
		}

		delete(socket.users, *user)
		socket.RoomListr(server)
	})
}

func (socket *SocketsService) SaveMessages(sender, message, room string, messageType int64) {
	redisText := fmt.Sprintf(`{"sender":"%v", "message":"%v", "type":"%v", "time":"%v"}`, sender, message, messageType, time.Now().Format("2006/02/01 15:04:05"))
	_, err := db.Rdb.LPush(fmt.Sprint(room), redisText).Result()
	if err != nil {
		log.Println(err)
		return
	}
}

func (socket *SocketsService) RoomListr(server *socketio.Server) {
	var room = models.Rooms{
		RoomsConf: nil,
		Max:       0,
	}
	for user := range socket.users {
		if user.Role != SUPERVISOR && user.Role != SUPPORT {
			room.Max++
			room.RoomsConf = append(room.RoomsConf, models.RoomsConf{
				RoomID:         user.Room,
				Name:           user.Name,
				Time:           time.Now().Format("2006/02/01 15:04"),
				ImportantLevel: user.ImportantLevel,
				Counter:        server.RoomLen("/", user.Room),
			})
		}
	}

	if ok := server.BroadcastToRoom("/", "0", "conf", room); !ok {
		log.Println("Error while sending client coute")
		return
	}
}

func (socket *SocketsService) SupportIsAvailable(user *models.User) {
	var countUser int64 = 0
	var countSup int64 = 0

	for user := range socket.users {
		switch user.Role {
		case SUPPORT:
			countSup++
		default:
			countUser++
		}
	}

	if countSup < countUser {
		sender <- fmt.Sprintf(`{"sender":"server", "message":"Сейчас все консультанты занять, на ваш вопрос ответять позже...", "time":"%v"}`, time.Now().Format("2006/02/01 15:04:05"))
		go WaitTime(5, user)
		return
	}

	sender <- fmt.Sprintf(`{"sender":"server", "message":"Консультант скоро ответит на ваш воопрос.", "time":"%v"}`, time.Now().Format("2006/02/01 15:04:05"))
	go WaitTime(2, user)
}

func WaitTime(t int64, user *models.User) {
	t2 := time.NewTimer(time.Duration(t) * time.Minute)
	for {
		select {
		case <-t2.C:
			// SenderBot(client)
			return
		case c := <-ch:
			if user.Room == c {
				return
			}
		}
	}
}

func ParseUsers(s socketio.Conn) (*models.User, error) {
	role, err := strconv.ParseInt(s.RemoteHeader().Get("role"), 10, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	IL, err := strconv.ParseInt(s.RemoteHeader().Get("importantLevel"), 10, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.User{
		ID:             s.RemoteHeader().Get("userID"),
		Name:           s.RemoteHeader().Get("FIO"),
		Role:           role,
		Room:           s.ID(),
		ImportantLevel: IL,
	}, nil
}
