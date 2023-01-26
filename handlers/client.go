package handlers

import (
	"feedbacks/models"
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
	send     chan models.SendMessage
	user     models.User
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
	var user models.User
	id, err := strconv.ParseInt(c.Query("userID"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	user.ID = id
	user.Room_id = id

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
		user.Name = "AMIN"
		user.Role = TERMINAL
	case 4:
		user.Name = "SHER"
		user.Role = DEALER
	case 5:
		user.Name = "UMED"
		user.Role = MANAGER
	default:
		user.Name = "ABROR"
		user.Role = DEVELOPER
	}
	log.Println("user:", user)

	client := newClient(user, conn, Ws)

	go client.readMessage()
	go client.writeMessage()

	Ws.register <- client
}

func newClient(user models.User, conn *websocket.Conn, ws *WsServer) *Chat {
	return &Chat{
		conn:     conn,
		wsServer: ws,
		send:     make(chan models.SendMessage, 256),
		user: models.User{
			ID:      user.ID,
			Name:    user.Name,
			Role:    user.Role,
			Room_id: user.ID,
		},
	}
}
