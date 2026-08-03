package main

import (
	"log"

	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/database"
)

func main() {
	cfg := configs.LoadConfig()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")
	//router := gin.Default()//gin router
	log.Println(" Server running")

	/*if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}*/

}
