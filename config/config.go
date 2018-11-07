package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DatabaseURL      string `envconfig:"DB_URL" required:"true"`
	TelegramBotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	UseWebhook       bool   `envconfig:"USE_WEBHOOK"`
	HerokuBaseUrl    string `envconfig:"HEROKU_BASE_URL"`
	Debug            bool   `envconfig:"DEBUG"`
	LogFile          string `envconfig:"LOG_FILE"`
	Port             string `envconfig:"PORT"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return &c, err
	}
	return &c, nil
}
