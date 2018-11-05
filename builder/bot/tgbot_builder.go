package builder

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teimurjan/go-els-tg-bot/config"
)

func MakeTelegramBot(c *config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(c.TelegramBotToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = c.Debug

	return bot, nil
}

func MakeTelegramBotUpdateConfig() *tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5
	return &u
}
