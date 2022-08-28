package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-els-tg-bot/config"
	checkUpdates "github.com/teimurjan/go-els-tg-bot/job/check_address_updates"
	"github.com/teimurjan/go-els-tg-bot/logging"
	"github.com/teimurjan/go-els-tg-bot/storage"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(c)

	logger.Info("Connecting to the database.")
	db, err := storage.NewPostgreSQL(c)
	if err != nil {
		logger.Fatal("Can't create a database connection.", err)
	}
	logger.Info("Database is successfully connected.")

	j := checkUpdates.NewTgBotJob(c, db, logger)

	logger.Info("Starting the check updates job.")
	j.Do()
	logger.Info("Check updates job is finished.")
}
