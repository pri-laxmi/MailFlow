package main

import (
	"log"

	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/apis"
	"github.com/pri-laxmi/MailFlow/internal/database"
	"github.com/pri-laxmi/MailFlow/internal/models"
	"github.com/pri-laxmi/MailFlow/internal/utils"
)

func main() {
	//load config
	cfg := configs.LoadConfig()
	jwtManager := utils.NewJWTManager(cfg.JWTSecret)
	//connect to database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Database connected successfully")
	if err := db.AutoMigrate(&models.User{}, &models.Contact{}, &models.Template{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	//setup router
	router := apis.SetRoutes(cfg, db, jwtManager)
	port := getPort(cfg)
	log.Printf("server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

func getPort(cfg *configs.Config) string {
	if cfg.ServerPort == "" {
		return "8080"
	}
	return cfg.ServerPort
}
