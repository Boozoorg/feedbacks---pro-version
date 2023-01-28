package handlers

import (
	// handlers "feedbacks/handlers/middleware"
	"encoding/json"
	"feedbacks/db"
	"feedbacks/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var CORSTrue = func(r *http.Request) bool {
	return true
}
var SIo *SocketsService

func LaunchRoutes() {
	SIo = NewWebsocketServer()

	router := gin.Default()
	{
		//############################################## Работа с клиентами ####################################################\\

		//################################################### Чатинг ##########################################################\\
		server := socketio.NewServer(&engineio.Options{
			Transports: []transport.Transport{
				&polling.Transport{
					CheckOrigin: CORSTrue,
				},
				&websocket.Transport{
					CheckOrigin: CORSTrue,
				},
			},
		})
		SIo.SocketsEvents(server)
		go func() {
			if err := server.Serve(); err != nil {
				log.Fatalf("socketio listen error: %s\n", err)
			}
		}()
		defer server.Close()
		router.GET("/socket.io/*any", func(c *gin.Context) {
			log.Print(c.Request.Header.Get("userId"))
		}, gin.WrapH(server))

		router.GET("/chatHistory", func(c *gin.Context) {
			var message models.Message
			var chatHistory []models.Message
			value, err := db.Rdb.LRange(c.Request.Header.Get("userId"), 0, -1).Result()
			if err != nil {
				log.Println(err)
				return
			}
			a, _ := strconv.ParseInt(c.Query("userID"), 10, 64)
			switch a {
			case SUPERVISOR:
				for _, v := range value {
					err = json.Unmarshal([]byte(v), &chatHistory)
					log.Println(v)
					if err != nil {
						log.Printf("Error while marshaling: %v", err)
						return
					}
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
			}
			log.Println(chatHistory)
			c.JSON(200, chatHistory)
		})
		//################################################### F.A.Q ##########################################################\\
		router.GET("/getFaqs", GetFaqs)
		router.POST("/getFaqs", CreateFaqs)
		router.PUT("/getFaqs", UpdateFaqs)
		router.DELETE("/getFaqs", DeleteFaqs)
		//################################################# Telegram Bot ##########################################################\\

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "FEEDBACKS_PAGE_NOT_FOUND", "message": "FDBK app page not found"})
	})
	err := router.Run(":8999")
	if err != nil {
		log.Println("failed to run port")
		return
	}
}
