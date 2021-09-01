package utils

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type I18nErr struct {
	MessageID string
}

func (e *I18nErr) Error() string {
	return fmt.Sprintf("There is an error with the message id %v. Look at the translation files if you need an explanation.", e.MessageID)
}

func NewI18NErr(messageID string) *I18nErr {
	return &I18nErr{MessageID: messageID}
}

func GetErrorMessage(err error, localizer *i18n.Localizer) string {
	if err == nil {
		return ""
	}

	if i18NErr, ok := err.(*I18nErr); ok {
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: i18NErr.MessageID,
			},
		})
	}

	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "error",
		},
	})
}
