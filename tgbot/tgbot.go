package tgbot

import "github.com/go-telegram-bot-api/telegram-bot-api"

// TgBot is a telegram bot interface
type TgBot interface {
	Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
}
