package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teimurjan/go-els-tg-bot/config"
)

func TestNewPostgreSQL(t *testing.T) {
	c := config.Config{
		DatabaseURL:      "",
		TelegramBotToken: "bot_token",
		UseWebhook:       false,
		HerokuBaseUrl:    "heroku_base_url",
		Debug:            true,
		LogFile:          "path/to/log",
		Port:             "8080",
	}
	db, err := NewPostgreSQL(&c)
	assert.NoError(t, err)
	assert.Equal(t, "postgres", db.DriverName())
}
