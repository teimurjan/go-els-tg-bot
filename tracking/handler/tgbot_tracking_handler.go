package handler

import (
	"fmt"

	"github.com/teimurjan/go-els-tg-bot/tgbot"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teimurjan/go-els-tg-bot/texts"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	utils "github.com/teimurjan/go-els-tg-bot/utils/arguments"
)

type tgbotTrackingHandler struct {
	service tracking.TrackingService
	bot     tgbot.TgBot
}

// NewTgbotTrackingHandler creates new tgbotTrackingHandler instance
func NewTgbotTrackingHandler(
	service tracking.TrackingService,
	bot tgbot.TgBot,
) tracking.TrackingHandler {
	return &tgbotTrackingHandler{
		service,
		bot,
	}
}

func (h *tgbotTrackingHandler) AddTracking(arguments string, chatID int64) {
	parsedArguments := utils.ParseArguments(arguments)
	trackingNumber, trackingOk := parsedArguments["v"]
	name, nameOk := parsedArguments["n"]
	if !trackingOk || !nameOk {
		h.bot.Send(tgbotapi.NewMessage(chatID, texts.NotEnoughArgumentsForTracking))
		return
	}

	tracking, err := h.service.Create(trackingNumber, name, chatID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err)))
		return
	}
	msgAdded := tgbotapi.NewMessage(chatID, texts.GetTrackingAddedMessage())
	msgInfo := tgbotapi.NewMessage(chatID, texts.GetTrackingInfoMessage(tracking))
	msgInfo.ParseMode = tgbotapi.ModeMarkdown

	h.bot.Send(msgAdded)
	h.bot.Send(msgInfo)
}

func (h *tgbotTrackingHandler) GetAll(chatID int64) {
	trackings, err := h.service.GetAll(chatID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err)))
		return
	}

	if len(trackings) == 0 {
		h.bot.Send(tgbotapi.NewMessage(chatID, texts.GetNoTrackingsMessage()))
		return
	}

	for _, tracking := range trackings {
		inlineBtns := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(
					texts.Delete,
					fmt.Sprintf("/delete_tracking %d", tracking.ID),
				),
			},
		)

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
			texts.TrackingInfoTempl,
			tracking.Name,
			tracking.Status,
			tracking.Value,
		))
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = inlineBtns
		h.bot.Send(msg)
	}
}

func (h *tgbotTrackingHandler) DeleteTracking(trackingID int64, chatID int64, messageID int64) {
	err := h.service.Delete(trackingID)
	var msg tgbotapi.Chattable
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err))
	} else {
		msg = tgbotapi.NewDeleteMessage(chatID, int(messageID))
	}
	h.bot.Send(msg)
}

func (h *tgbotTrackingHandler) CheckUpdates() {
	updates, err := h.service.GetUpdates()
	if err == nil {
		for _, update := range updates {
			msg := tgbotapi.NewMessage(update.User.ChatID, fmt.Sprintf(
				texts.GetTrackingUpdatedMessage(),
				update.Tracking.Name,
				update.Tracking.Status,
				update.Tracking.Value,
			))
			msg.ParseMode = tgbotapi.ModeMarkdown
			h.bot.Send(msg)
		}
	}
}
