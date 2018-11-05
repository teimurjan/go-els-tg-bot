package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/teimurjan/go-els-tg-bot/application"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
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

	app := application.NewApp(c)
	app.Start()
}
