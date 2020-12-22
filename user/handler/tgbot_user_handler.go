package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	helper "github.com/teimurjan/go-els-tg-bot/helper/i18n"
	"github.com/teimurjan/go-els-tg-bot/tgbot"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type tgbotUserHandler struct {
	service    user.UserService
	bot        tgbot.TgBot
	i18nHelper helper.I18nHelper
}

// NewTgbotUserHandler creates a new instance of user handler for telegram bot
func NewTgbotUserHandler(service user.UserService, bot tgbot.TgBot, i18nHelper helper.I18nHelper) user.UserHandler {
	return &tgbotUserHandler{
		service,
		bot,
		i18nHelper,
	}
}

// Join adds a new user
func (h *tgbotUserHandler) Join(chatID int64) {
	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	_, err := h.service.Create(chatID)
	var msg tgbotapi.MessageConfig
	if err != nil {
		msg = tgbotapi.NewMessage(chatID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "error"}))
	} else {
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "greetings",
				Other: "Hi there! 👋\nStart monitoring your orders by typing:\n/add_tracking",
			},
		})
		msg = tgbotapi.NewMessage(chatID, text)
	}
	h.bot.Send(msg)
}

// ChangeLanguage changes user language
func (h *tgbotUserHandler) RequestLanguageChange(chatID int64) {
	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	text := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "requestLanguageChange",
			Other: "Which language do you prefer?",
		},
	})
	inlineBtns := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"English",
				"/change_language en",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"Русский",
				"/change_language ru",
			),
		},
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = inlineBtns
	h.bot.Send(msg)
}

// ChangeLanguage changes user language
func (h *tgbotUserHandler) ChangeLanguage(language string, chatID int64, messageID int64) {
	h.service.Update(chatID, language)

	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	text := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "changeLanguage",
			Other: "Language is successfully changed.",
		},
	})

	textEditMsg := tgbotapi.NewEditMessageText(chatID, int(messageID), text)
	replyMarkupEditMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, int(messageID), tgbotapi.NewInlineKeyboardMarkup())

	h.bot.Send(replyMarkupEditMsg)
	h.bot.Send(textEditMsg)
}
