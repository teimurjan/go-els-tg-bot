package app

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"

	"github.com/teimurjan/go-els-tg-bot/commands"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	botFactory "github.com/teimurjan/go-els-tg-bot/factory/bot"
	containersFactory "github.com/teimurjan/go-els-tg-bot/factory/containers"
	"github.com/teimurjan/go-els-tg-bot/utils/callbacks"
)

type tgBotApp struct {
	conf              *config.Config
	logger            *logrus.Logger
	bot               *tgbotapi.BotAPI
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

// NewTgBotApp creates new tg bot application
func NewTgBotApp(conf *config.Config, db *sqlx.DB, logger *logrus.Logger) *tgBotApp {
	bot, err := botFactory.MakeTelegramBot(conf)
	if err != nil {
		logger.Fatal("Can't create a telegram bot.", err)
	}

	reposContainer := containersFactory.MakeReposContainer(db)
	servicesContainer := containersFactory.MakeServicesContainer(reposContainer, logger)
	handlersContainer := containersFactory.MakeHandlersContainer(servicesContainer, bot)

	return &tgBotApp{
		conf,
		logger,
		bot,
		reposContainer,
		servicesContainer,
		handlersContainer,
	}
}

func (tgBotApp *tgBotApp) Start() {
	updates, err := tgBotApp.getUpdates()
	if err != nil {
		tgBotApp.logger.Fatal("Can't get bot updates.", err)
	}

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			tgBotApp.handleCommand(&update)
		} else if update.CallbackQuery != nil {
			tgBotApp.handleCallback(&update)
		}
	}
}

func (tgBotApp *tgBotApp) getUpdates() (tgbotapi.UpdatesChannel, error) {
	if !tgBotApp.conf.UseWebhook {
		return tgBotApp.setupPolling()
	}
	return tgBotApp.setupWebhook()
}

func (tgBotApp *tgBotApp) setupPolling() (tgbotapi.UpdatesChannel, error) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 5
	tgBotApp.logger.Info("Start polling.")
	return tgBotApp.bot.GetUpdatesChan(updateConfig)
}

func (tgBotApp *tgBotApp) setupWebhook() (tgbotapi.UpdatesChannel, error) {
	webhookURL := tgBotApp.conf.HerokuBaseUrl + "/" + tgBotApp.bot.Token
	_, err := tgBotApp.bot.SetWebhook(
		tgbotapi.NewWebhook(webhookURL),
	)
	if err != nil {
		tgBotApp.logger.Fatal("There is a problem in setting webhook.", err)
		return nil, err
	}
	updates := tgBotApp.bot.ListenForWebhook("/" + tgBotApp.bot.Token)
	go http.ListenAndServe(":"+tgBotApp.conf.Port, nil)

	tgBotApp.logger.Info("Listening port " + tgBotApp.conf.Port + ". Webhook url is " + webhookURL + ".")
	return updates, nil
}

func (tgBotApp *tgBotApp) handleCommand(update *tgbotapi.Update) {
	command := update.Message.Command()
	if command == commands.Start {
		go tgBotApp.handlersContainer.UserHandler.Join(update.Message.Chat.ID)
	} else if command == commands.AddTracking {
		commandArgs := update.Message.CommandArguments()
		go tgBotApp.handlersContainer.TrackingHandler.AddTracking(commandArgs, update.Message.Chat.ID)
	} else if command == commands.GetAll {
		go tgBotApp.handlersContainer.TrackingHandler.GetAll(update.Message.Chat.ID)
	}
}

func (tgBotApp *tgBotApp) handleCallback(update *tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	if trackingID, err := callbacks.ParseDeleteTrackingCallback(callbackData); err == nil {
		go tgBotApp.handlersContainer.TrackingHandler.DeleteTracking(
			trackingID,
			update.CallbackQuery.Message.Chat.ID,
			int64(update.CallbackQuery.Message.MessageID),
		)
	}
}
