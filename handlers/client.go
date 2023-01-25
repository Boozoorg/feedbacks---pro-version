package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
	role     string
	room     int64
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
func OpenChat(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("userID"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	var name string
	var role string
	log.Println("role:", c.Request.Header.Get("role"))
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		c.Abort()
		return
	}
	switch id {
	case 3:
		name = "AMIN"
		role = "TERMINAL"
	case 4:
		name = "SHER"
		role = "DEALER"
	case 5:
		name = "UMED"
		role = "MANAGER"
	default:
		name = "ABROR"
		role = "DEVELOPER"
	}
	log.Println("name:", name, "role:", role, "id:", id)

	client := newClient(name, role, id, conn, Ws)

	go client.readMessage()
	go client.writeMessage()

	Ws.register <- client
}

func newClient(name, role string, id int64, conn *websocket.Conn, ws *WsServer) *Chat {
	return &Chat{
		name:     name,
		conn:     conn,
		wsServer: ws,
		send:     make(chan []byte, 256),
		room:     id,
		role:     role,
	}
}
