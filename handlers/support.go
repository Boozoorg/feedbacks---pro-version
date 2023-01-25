package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

// @Summary Supports choose chat to join
// @ID      join-chat-room
// @Produce json
// @Tags    Chat
// @Param   fio     path     string true  "supports name"                                                                      Format(fio)
// @Param   room_id path     string false "id of room which support must join or if want get max users and there id set it 0)" Format(roomID)
// @Success 200     {object} string
// @Failure 404     {object} string
// @Router  /joinRoom [get]
func SupportJoinChat(c *gin.Context) {
	var id int64
	var name string
	var role string
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomID := c.Query("roomID")
	if roomID != "" {
		id, err = strconv.ParseInt(roomID, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		if id < 0 {
			conn.WriteMessage(1, []byte("сэр вы абасрались"))
			conn.Close()
			log.Println("id of room cann't be low than 0")
			return
		}
	} else {
		id = 0
	}
	switch c.Query("userID") {
	case "1":
		name = "BUZURG"
		role = "SUPERVISOR"
	default:
		name = "KOMRON"
		role = "SUPPORT"
	}
	log.Println("name:", name, "role:", role, "id:", id)

	client := newSupport(name, role, id, conn, Ws)
	go client.readMessage()
	go client.writeMessage()

	Ws.register <- client
}

func newSupport(name, role string, roomID int64, conn *websocket.Conn, wsServer *WsServer) *Chat {
	return &Chat{
		name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		room:     roomID,
		role:     role,
	}
}
