package handler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teimurjan/go-els-tg-bot/addTrackingDialog"
	"github.com/teimurjan/go-els-tg-bot/texts"
)

type addTrackingDialogHandler struct {
	service addTrackingDialog.AddTrackingDialogService
	bot     *tgbotapi.BotAPI
}

func NewTgbotAddTrackingDialogHandler(
	service addTrackingDialog.AddTrackingDialogService,
	bot *tgbotapi.BotAPI,
) *addTrackingDialogHandler {
	return &addTrackingDialogHandler{
		service,
		bot,
	}
}

func (h *addTrackingDialogHandler) StartDialog(chatID int64) {
	var msg string
	_, err := h.service.StartDialog(chatID)
	msg = texts.GetEnterOrderNameMessage()
	if err != nil {
		msg = texts.GetErrorMessage(err)
	}
	h.bot.Send(tgbotapi.NewMessage(chatID, msg))
}

func (h *addTrackingDialogHandler) UpdateDialogIfActive(text string, chatID int64) {
	dialog, err := h.service.GetDialogForChat(chatID)
	if dialog != nil && dialog.Step == 1 {
		err = h.service.UpdateDialogName(dialog, text)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err)))
			return
		}
		h.bot.Send(tgbotapi.NewMessage(chatID, texts.GetEnterTrackingMessage()))
		return
	}

	if dialog != nil && dialog.Step == 2 {
		tracking, err := h.service.UpdateDialogTracking(dialog, text)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err)))
			return
		}
		msgAdded := tgbotapi.NewMessage(chatID, texts.GetTrackingAddedMessage())

		msgInfo := tgbotapi.NewMessage(chatID, texts.GetTrackingInfoMessage(tracking))
		msgInfo.ParseMode = tgbotapi.ModeMarkdown

		h.bot.Send(msgAdded)
		h.bot.Send(msgInfo)
		return
	}
}

func (h *addTrackingDialogHandler) ResetDialog(chatID int64) {
	dialog, err := h.service.GetDialogForChat(chatID)
	if err == nil {
		h.service.ResetDialog(dialog)
	}
}
