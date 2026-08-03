package main

import (

	"github.com/pri-laxmi/MailFlow/configs"
	"github.com/pri-laxmi/MailFlow/internal/database"
)

func main() {
	config := configs.LoadConfig()
	db := database.ConnectDB(config)

	defer db.Close()
}
