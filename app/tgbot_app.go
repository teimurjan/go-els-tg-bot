package app

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"

	"github.com/robfig/cron"
	builder "github.com/teimurjan/go-els-tg-bot/builder/bot"
	"github.com/teimurjan/go-els-tg-bot/commands"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	trackingFetcher "github.com/teimurjan/go-els-tg-bot/tracking/fetcher"
	trackingHandler "github.com/teimurjan/go-els-tg-bot/tracking/handler"
	trackingRepository "github.com/teimurjan/go-els-tg-bot/tracking/repository"
	trackingService "github.com/teimurjan/go-els-tg-bot/tracking/service"
	userHandler "github.com/teimurjan/go-els-tg-bot/user/handler"
	userRepository "github.com/teimurjan/go-els-tg-bot/user/repository"
	userService "github.com/teimurjan/go-els-tg-bot/user/service"
	"github.com/teimurjan/go-els-tg-bot/utils/callbacks"
)

type tgBotApp struct {
	conf              *config.Config
	logger            *logrus.Logger
	db                *sqlx.DB
	bot               *tgbotapi.BotAPI
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

func NewTgBotApp(conf *config.Config, db *sqlx.DB, logger *logrus.Logger) *tgBotApp {
	bot, err := builder.MakeTelegramBot(conf)
	if err != nil {
		logger.Fatal("Can't create a telegram bot.", err)
	}

	reposContainer := containers.NewRepositoriesContainer(
		userRepository.NewPostgresqlUserRepository(db),
		trackingRepository.NewPostgresqlTrackingRepository(db),
	)

	servicesContainer := containers.NewServicesContainer(
		userService.NewUserService(reposContainer.UserRepo),
		trackingService.NewTrackingService(
			reposContainer.TrackingRepo,
			reposContainer.UserRepo,
			trackingFetcher.NewTrackingStatusFetcher(),
			logger,
		),
	)

	handlersContainer := containers.NewHandlersContainer(
		userHandler.NewTgbotUserHandler(servicesContainer.UserService, bot),
		trackingHandler.NewTgbotTrackingHandler(servicesContainer.TrackingService, bot),
	)

	return &tgBotApp{
		conf,
		logger,
		db,
		bot,
		reposContainer,
		servicesContainer,
		handlersContainer,
	}
}

func (tgBotApp *tgBotApp) Start() {
	updates, err := tgBotApp.getBotUpdates()
	if err != nil {
		tgBotApp.logger.Fatal("Can't get bot updates.", err)
	}

	c := cron.New()
	c.AddFunc("@every 1m", tgBotApp.handlersContainer.TrackingHandler.CheckUpdates)
	c.Start()

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			tgBotApp.handleCommand(&update)
		} else if update.CallbackQuery != nil {
			tgBotApp.handleCallback(&update)
		}
	}
}

func (tgBotApp *tgBotApp) getBotUpdates() (tgbotapi.UpdatesChannel, error) {
	if !tgBotApp.conf.UseWebhook {
		updateConfig := builder.MakeTelegramBotUpdateConfig()
		tgBotApp.logger.Info("Start polling.")
		return tgBotApp.bot.GetUpdatesChan(*updateConfig)
	}

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
