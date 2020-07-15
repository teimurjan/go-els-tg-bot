package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teimurjan/go-els-tg-bot/config"
)

func TestNewMySQL(t *testing.T) {
	c := config.Config{
		DatabaseURL:      "jdbc:mysql://localhost:3306/test",
		TelegramBotToken: "bot_token",
		UseWebhook:       false,
		HerokuBaseUrl:    "heroku_base_url",
		Debug:            true,
		LogFile:          "path/to/log",
		Port:             "8080",
	}
	db, err := NewMySQL(&c)
	assert.NoError(t, err)
	assert.Equal(t, "mysql", db.DriverName())
}
