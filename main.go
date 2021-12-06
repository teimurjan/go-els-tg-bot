package main

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/tracking/fetcher"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	f := fetcher.NewTrackingNumberFetcher(c)
	fmt.Println(f.Fetch("1ZV501970443106874"))
}
