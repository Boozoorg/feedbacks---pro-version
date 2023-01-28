package handlers

import (
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 5,
	WriteBufferSize: 1024 * 1024 * 5,
}

type Chat struct {
	wsServer *WsServer
	conn     *websocket.Conn
	send     chan []byte
	name     string
	room     int64
	support  bool
}

// @Summary Agent uses this route to create chat room
// @ID      create-chat-room
// @Produce json
// @Tags    Chat
// @Param   user_id path     string true "id of user" Format(userID)
// @Param   fio     path     string true "users name" Format(fio)
// @Success 200     {object} string
// @Failure 404     {object} string
// @Router  /createRoom [get]
// func OpenChat(c *gin.Context) {
// 	id := c.Query("userID")
// 	name := c.Query("fio")

// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
// 		c.Abort()
// 		return
// 	}

// 	client := newClient(name, id, conn, Ws)

// 	go client.readMessage()
// 	go client.writeMessage()

// Ws.register <- client
// }

// func newClient(name, ID string, conn *websocket.Conn, wsServer *WsServer) *Chat {
// 	id, err := strconv.ParseInt(ID, 10, 64)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	return &Chat{
// 		name:     name,
// 		conn:     conn,
// 		wsServer: wsServer,
// 		send:     make(chan []byte, 256),
// 		room:     id,
// 		support:  false,
// 	}
// }
