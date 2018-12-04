package handler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teimurjan/go-els-tg-bot/addTrackingDialog"
	"github.com/teimurjan/go-els-tg-bot/errs"
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
		msg = texts.GetErrorMessage()
	}
	h.bot.Send(tgbotapi.NewMessage(chatID, msg))
}

func (h *addTrackingDialogHandler) UpdateDialogIfActive(text string, chatID int64) {
	var msg string
	dialog, err := h.service.GetDialogForChat(chatID)
	if dialog != nil && dialog.Step == 1 {
		err = h.service.UpdateDialogName(dialog, text)
		msg = texts.GetEnterTrackingMessage()
	} else if dialog != nil && dialog.Step == 2 {
		err = h.service.UpdateDialogTracking(dialog, text)
		msg = texts.GetTrackingAddedMessage()
	}
	if err != nil {
		switch err.(type) {
		case *errs.Err:
			msg = err.Error()
		default:
			msg = texts.GetErrorMessage()
		}
	}
	h.bot.Send(tgbotapi.NewMessage(chatID, msg))
}

func (h *addTrackingDialogHandler) ResetDialog(chatID int64) {
	dialog, err := h.service.GetDialogForChat(chatID)
	if err == nil {
		h.service.ResetDialog(dialog)
	}
}
