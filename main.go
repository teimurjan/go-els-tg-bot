package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-els-tg-bot/application"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	"github.com/teimurjan/go-els-tg-bot/logging"
)

type app struct {
	reposContainer    containers.RepositoriesContainer
	servicesContainer containers.ServicesContainer
	handlersContainer containers.HandlersContainer
}

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(c)
	app := application.NewApp(c, logger)
	app.Start()
}
