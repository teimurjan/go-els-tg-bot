package helper

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-els-tg-bot/user"
	"golang.org/x/text/language"
)

type I18nHelper interface {
	MustGetLocalizer(chatID int64) *i18n.Localizer
}

type i18nHelper struct {
	bundle   *i18n.Bundle
	userRepo user.UserRepository
	logger   *logrus.Logger
}

func NewI18nHelper(userRepo user.UserRepository, logger *logrus.Logger) I18nHelper {
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("active.en.toml")
	bundle.MustLoadMessageFile("active.ru.toml")

	return &i18nHelper{bundle, userRepo, logger}
}

func (helper *i18nHelper) MustGetLocalizer(chatID int64) *i18n.Localizer {
	user, err := helper.userRepo.GetByChatID(chatID)
	if err != nil {
		helper.logger.Info(fmt.Sprintf("There is no user with chat_id=%d. Using default locale (ru).", chatID))
		return i18n.NewLocalizer(helper.bundle, "ru")
	}

	helper.logger.Info(fmt.Sprintf("Using locale %s for user %d.", user.Language, user.ID))

	return i18n.NewLocalizer(helper.bundle, user.Language)
}
