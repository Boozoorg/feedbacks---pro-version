package handlers

import (
	// handlers "feedbacks/handlers/middleware"
	"feedbacks/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var Ws *WsServer

func LaunchRoutes() {
	Ws = NewWebsocketServer()
	go Ws.Run()
	
	router := gin.Default()
	router.GET("/ping", ping)
	// router.Use(handlers.ApiMiddleware())
	{
		//############################################## Работа с клиентами ####################################################\\

		//################################################### Чатинг ##########################################################\\
		router.GET("/createRoom", OpenChat)      // Содание комнаты пользователем )
		router.GET("/joinRoom", SupportJoinChat) // Саппорт заходит в комнату с пользователем (По дефолту комната с сапортами имеет id = 0 или roomID не записан)
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
	err := router.Run(fmt.Sprint(":", models.Conf.AppConf.Port))
	if err != nil {
		log.Println("failed to run port")
		return
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, "ping success! Gip-gip Urrraaaa!")
	u, err := getUser(c)
	if err != nil {
		return
	}
	log.Println("ping success from: ", u.Login)
}
