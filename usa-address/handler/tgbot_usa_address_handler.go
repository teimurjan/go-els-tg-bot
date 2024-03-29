package handler

import (
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	helper "github.com/teimurjan/go-els-tg-bot/helper/i18n"
	"github.com/teimurjan/go-els-tg-bot/tgbot"
	usaAddress "github.com/teimurjan/go-els-tg-bot/usa-address"
	"github.com/teimurjan/go-els-tg-bot/user"
	errsUtil "github.com/teimurjan/go-els-tg-bot/utils/errs"
)

type tgbotUsaAddressHandler struct {
	service     usaAddress.UsaAddressService
	userService user.UserService
	bot         tgbot.TgBot
	i18nHelper  helper.I18nHelper
}

// NewTgbotUsaAddressHandler creates new tgbotUsaAddressHandler instance
func NewTgbotUsaAddressHandler(
	service usaAddress.UsaAddressService,
	userService user.UserService,
	bot tgbot.TgBot,
	i18nHelper helper.I18nHelper,
) usaAddress.UsaAddressHandler {
	return &tgbotUsaAddressHandler{
		service,
		userService,
		bot,
		i18nHelper,
	}
}

func (h *tgbotUsaAddressHandler) GetAddress(chatID int64) {
	usaAddress, err := h.service.GetFirst()
	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, errsUtil.GetErrorMessage(err, localizer)))
		return
	}

	text := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    "usaAddressInfo",
		TemplateData: usaAddress,
	})
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	h.bot.Send(msg)
}

func (h *tgbotUsaAddressHandler) CheckDiff() {
	users, err := h.userService.GetAll()
	if err != nil {
		return
	}

	address, newAddress, changelog, err := h.service.CheckAddressUpdates()
	if err != nil || len(changelog) == 0 {
		return
	}

	for _, change := range changelog {
		value := reflect.ValueOf(address)

		fieldName := change.Path[len(change.Path)-1]
		fieldValue := value.Elem().FieldByName(fieldName)

		if fieldValue.Kind() == reflect.String {
			fieldValue.SetString(change.To.(string) + " 🆕")
		}
	}

	for _, user := range users {
		localizer := h.i18nHelper.MustGetLocalizer(user.ChatID)
		infoText := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:    "usaAddressInfo",
			TemplateData: newAddress,
		})
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "usaAddressUpdated",
				Other: "❗️❗️❗️ The following fields in the USA address have been changed ❗️❗️❗️",
			},
		}) + "\n\n" + infoText

		msg := tgbotapi.NewMessage(
			user.ChatID,
			text,
		)
		msg.ParseMode = tgbotapi.ModeMarkdown
		h.bot.Send(msg)
	}
}
