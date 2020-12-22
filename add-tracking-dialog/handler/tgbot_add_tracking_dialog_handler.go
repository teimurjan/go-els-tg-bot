package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	addTrackingDialog "github.com/teimurjan/go-els-tg-bot/add-tracking-dialog"
	helper "github.com/teimurjan/go-els-tg-bot/helper/i18n"
	"github.com/teimurjan/go-els-tg-bot/tgbot"
)

type addTrackingDialogHandler struct {
	service    addTrackingDialog.AddTrackingDialogService
	bot        tgbot.TgBot
	i18nHelper helper.I18nHelper
}

// NewTgbotAddTrackingDialogHandler creates new addTrackingDialogHandler instance
func NewTgbotAddTrackingDialogHandler(
	service addTrackingDialog.AddTrackingDialogService,
	bot tgbot.TgBot,
	i18nHelper helper.I18nHelper,
) addTrackingDialog.AddTrackingDialogHandler {
	return &addTrackingDialogHandler{
		service,
		bot,
		i18nHelper,
	}
}

func (h *addTrackingDialogHandler) StartDialog(chatID int64) {
	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	var text string
	_, err := h.service.StartDialog(chatID)
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "enterOrderName",
			Other: "What is the name of your order?",
		},
	})
	if err != nil {
		text = localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "error",
				Other: "Something went wrong. 🆘",
			},
		})
	}
	h.bot.Send(tgbotapi.NewMessage(chatID, text))
}

func (h *addTrackingDialogHandler) UpdateDialogIfActive(text string, chatID int64) {
	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	dialog, err := h.service.GetDialogForChat(chatID)
	if dialog != nil && dialog.Step == 1 {
		err = h.service.UpdateDialogName(dialog, text)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "error"})))
			return
		}
		text = localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "enterTrackingNumber",
				Other: "What is the tracking number of your order?",
			},
		})
		h.bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	if dialog != nil && dialog.Step == 2 {
		tracking, err := h.service.UpdateDialogTracking(dialog, text)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "error"})))
			return
		}

		addedText := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "trackingAdded",
		})
		infoText := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:    "trackingInfo",
			TemplateData: tracking,
		})
		msgAdded := tgbotapi.NewMessage(chatID, addedText)
		msgInfo := tgbotapi.NewMessage(chatID, infoText)
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
