package handlers

import (
	"feedbacks/models"
)

type WsServer struct {
	clients    map[*Chat]bool
	register   chan *Chat
	unregister chan *Chat
	receiver   chan models.ReceiveMessage
}

// func (server *WsServer) Run() {
// 	for {
// 		select {
// 		case client := <-server.register:
// 			server.registerClient(client)

// 		case client := <-server.unregister:
// 			server.unregisterClient(client)

// 		case receive := <-server.receiver:
// 			server.receiveMessage(receive.Name, receive.Message, receive.Room_id, receive.MessageType)
// 		}
// 	}
// }

// func (client *Chat) readMessage() {
// 	defer func() {
// 		client.disconnect()
// 	}()
// 	var n = 1

// 	for {
// 		messageType, messsage, err := client.conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				if err.Error() == "websocket: close 1000 (normal)" {
// 					break
// 				}
// 				log.Printf("unexpected close error: %v", err)
// 			}
// 			break
// 		}
// 		if messageType == 1 {
// 			client.wsServer.receiver <- models.Receiver{Name: client.name, Message: string(messsage), MessageType: 1, Room_id: client.room}
// 		} else if messageType == 2 {
// 			fileType := http.DetectContentType(messsage)
// 			switch fileType {
// 			case "image/jpeg":
// 				fileType = ".jpg"
// 			case "image/png":
// 				fileType = ".png"
// 			default:
// 				log.Println("error while detecting file type")
// 				return
// 			}

// 			err = ioutil.WriteFile(filepath.Join("files/temp", "photo_"+strconv.Itoa(n)+fileType), messsage, 0700)
// 			if err != nil {
// 				log.Println("error while creating photo in files/temp: ", err)
// 				return
// 			}

// 			client.wsServer.receiver <- models.Receiver{Name: client.name, Message: "files/temp/photo_" + strconv.Itoa(n) + fileType, MessageType: 2, Room_id: client.room}
// 			n++
// 		}
// 	}
// }

// func (client *Chat) writeMessage() {
// 	defer client.conn.Close()

// 	for {
// 		select {
// 		case message, ok := <-client.send:
// 			if !ok {
// 				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := client.conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

// func (client *Chat) disconnect() {
// 	client.wsServer.unregister <- client
// 	close(client.send)
// 	client.conn.Close()
// }

// func (server *WsServer) registerClient(client *Chat) {
// 	server.clients[client] = true
// 	server.roomListr()

// 	if client.support && client.room != 0 {
// 		ch <- client.room
// 		chatHistory(client)
// 		server.receiveMessage(client.name, fmt.Sprintf("Здравствуйте, меня зовут %v и я кансультант.", client.name), client.room, models.Text)
// 	} else if !client.support {
// 		server.supportIsAvailable(client)
// 		chatHistory(client)
// 	}
// }

// func chatHistory(chat *Chat) {
// 	var lastChats models.Message
// 	var History []string
// 	value, err := db.Rdb.LRange(fmt.Sprint(chat.room), 0, -1).Result()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	if chat.support {
// 		for _, v := range value {
// 			err = json.Unmarshal([]byte(v), &lastChats)
// 			if err != nil {
// 				log.Printf("Error while marshaling: %v", err)
// 				return
// 			}
// 			if lastChats.Time > time.Now().Add(-1*time.Hour).Format("2006/02/01 15:04:05") {
// 				History = append(History, fmt.Sprintf(`{"sender":"%v", "message":"%v", "time":"%v"}`, lastChats.Sender, lastChats.Message, lastChats.Time))
// 			}
// 		}
// 		chat.send <- []byte(fmt.Sprintf("%v", History))
// 		return
// 	}

// 	for _, v := range value {
// 		err = json.Unmarshal([]byte(v[:len(v)-1]), &lastChats)
// 		if err != nil {
// 			log.Printf("Error while marshaling: %v", err)
// 			return
// 		}
// 		if lastChats.Time > time.Now().Add(-20*time.Minute).Format("2006/02/01 15:04:05") {
// 			chat.send <- []byte(fmt.Sprintf(`{"sender":"%v", "message":"%v", "time":"%v"}`, lastChats.Sender, lastChats.Message, lastChats.Time))
// 		}
// 	}
// }

// func (server *WsServer) unregisterClient(client *Chat) {
// 	if _, ok := server.clients[client]; ok {
// 		delete(server.clients, client)
// 	}
// 	server.roomListr()
// 	if client.support {
// 		server.receiveMessage("server", fmt.Sprintf("Клиент: %v, покинул чат.", client.name), client.room, models.Text)
// 	}
// }

// func (server *WsServer) receiveMessage(sender, message string, room, messageType int64) {
// 	redisText := fmt.Sprintf(`{"sender":"%v", "message":"%v", "type":"%v", "time":"%v"}`, sender, message, messageType, time.Now().Format("2006/02/01 15:04:05"))
// 	_, err := db.Rdb.LPush(fmt.Sprint(room), redisText).Result()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	for client := range server.clients {
// 		if client.room == room {
// 			client.send <- []byte(fmt.Sprintf(`{"sender":"%v", "message":"%v", "time":"%v"}`, sender, message, time.Now().Format("2006/02/01 15:04:05")))
// 		}
// 	}
// }

// func (server *WsServer) roomListr() {
// 	var room models.Rooms
// 	var count int64 = 0
// 	for client := range server.clients {
// 		if !client.support {
// 			count++
// 			room.RoomsID = append(room.RoomsID, client.room)
// 		}
// 	}
// 	room.Max = count

// 	message := []byte(fmt.Sprintf(`{"max_users":"%v", "ID":"%v"}`, room.Max, room.RoomsID))
// 	for client := range server.clients {
// 		if client.room == 0 {
// 			client.send <- message
// 		}
// 	}
// }

// func (server *WsServer) supportIsAvailable(client *Chat) {
// 	var countUser int64 = 0
// 	var countSup int64 = 0

// 	for client := range server.clients {
// 		if !client.support {
// 			countUser++
// 		} else {
// 			countSup++
// 		}
// 	}

// 	if countSup < countUser {
// 		client.send <- []byte(fmt.Sprintf(`{"sender":"server", "message":"Сейчас все консультанты занять, на ваш вопрос ответять позже...", "time":"%v"}`, time.Now().Format("2006/02/01 15:04:05")))
// 		go waitTime(5, client)
// 		return
// 	}

// 	client.send <- []byte(fmt.Sprintf(`{"sender":"server", "message":"Консультант скоро ответит на ваш воопрос.", "time":"%v"}`, time.Now().Format("2006/02/01 15:04:05")))
// 	go waitTime(2, client)
// }

// func waitTime(t int64, client *Chat) {
// 	t2 := time.NewTimer(time.Duration(t) * time.Minute)
// 	for {
// 		select {
// 		case <-t2.C:
// 			// SenderBot(client)
// 			return
// 		case c := <-ch:
// 			if client.room == c {
// 				return
// 			}
// 		}
// 	}
// }
