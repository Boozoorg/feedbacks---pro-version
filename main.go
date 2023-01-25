package main

import (
	"feedbacks/db"
	"feedbacks/handlers"
	// "feedbacks/logging"
	"feedbacks/models"
	"log"
	_ "feedbacks/docs"
)

// @title       Swagger Example API
// @version     1.0
// @description this is swagger

// @host     localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Description for what is this security definition being used
func main() {
	log.Println("[FEEDBACKS service is starting . . .]")
	defer deferFunc()
	models.Setup("config/config.json")
	// logging.InitLogger()
	// db.SetupDB()
	handlers.LaunchRoutes()
}

func deferFunc() {
	log.Println("[FEEDBACKS Service is shutting down...]")
	db.CloseDB()
}