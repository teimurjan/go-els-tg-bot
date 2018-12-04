package handler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teimurjan/go-els-tg-bot/texts"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type tgbotUserHandler struct {
	service user.UserService
	bot     *tgbotapi.BotAPI
}

func NewTgbotUserHandler(service user.UserService, bot *tgbotapi.BotAPI) *tgbotUserHandler {
	return &tgbotUserHandler{
		service,
		bot,
	}
}

func (h *tgbotUserHandler) Join(chatID int64) {
	_, err := h.service.Create(chatID)
	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, texts.GetErrorMessage())
	} else {
		msg = tgbotapi.NewMessage(chatID, texts.GetWelcomeMessage())
	}
	h.bot.Send(msg)
}
