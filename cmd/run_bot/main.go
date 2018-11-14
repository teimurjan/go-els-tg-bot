package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-els-tg-bot/app"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/logging"
	"github.com/teimurjan/go-els-tg-bot/storage"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(c)

	db, err := storage.NewPostgreSQL(c)
	if err != nil {
		logger.Fatal("Can't create a database connection.", err)
	}

	a := app.NewTgBotApp(c, db, logger)
	a.Start()
}
