package handler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teimurjan/go-els-tg-bot/texts"
	"github.com/teimurjan/go-els-tg-bot/tgbot"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type tgbotUserHandler struct {
	service user.UserService
	bot     tgbot.TgBot
}

// NewTgbotUserHandler creates a new instance of user handler for telegram bot
func NewTgbotUserHandler(service user.UserService, bot tgbot.TgBot) user.UserHandler {
	return &tgbotUserHandler{
		service,
		bot,
	}
}

// Join adds a new user
func (h *tgbotUserHandler) Join(chatID int64) {
	_, err := h.service.Create(chatID)
	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err))
	} else {
		msg = tgbotapi.NewMessage(chatID, texts.GetWelcomeMessage())
	}
	h.bot.Send(msg)
}
