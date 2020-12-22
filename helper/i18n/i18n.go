package helper

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/teimurjan/go-els-tg-bot/user"
	"golang.org/x/text/language"
)

type I18nHelper interface {
	MustGetLocalizer(chatID int64) *i18n.Localizer
}

type i18nHelper struct {
	bundle   *i18n.Bundle
	userRepo user.UserRepository
}

func NewI18nHelper(userRepo user.UserRepository) I18nHelper {
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("active.en.toml")
	bundle.MustLoadMessageFile("active.ru.toml")

	return &i18nHelper{bundle, userRepo}
}

func (helper *i18nHelper) MustGetLocalizer(chatID int64) *i18n.Localizer {
	user, err := helper.userRepo.GetByChatID(chatID)
	if err != nil {
		panic(err)
	}

	return i18n.NewLocalizer(helper.bundle, user.Language)
}
