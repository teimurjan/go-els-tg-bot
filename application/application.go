package application

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"

	builder "github.com/teimurjan/go-els-tg-bot/builder/bot"
	"github.com/teimurjan/go-els-tg-bot/commands"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	"github.com/teimurjan/go-els-tg-bot/storage"
	"github.com/teimurjan/go-els-tg-bot/texts"
	trackingFetcher "github.com/teimurjan/go-els-tg-bot/tracking/fetcher"
	trackingHandler "github.com/teimurjan/go-els-tg-bot/tracking/handler"
	trackingRepository "github.com/teimurjan/go-els-tg-bot/tracking/repository"
	trackingService "github.com/teimurjan/go-els-tg-bot/tracking/service"
	userHandler "github.com/teimurjan/go-els-tg-bot/user/handler"
	userRepository "github.com/teimurjan/go-els-tg-bot/user/repository"
	userService "github.com/teimurjan/go-els-tg-bot/user/service"
	"github.com/teimurjan/go-els-tg-bot/utils/callbacks"
)

type app struct {
	conf              *config.Config
	db                *sqlx.DB
	bot               *tgbotapi.BotAPI
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

func NewApp(conf *config.Config) *app {
	db, err := storage.NewPostgreSQL(conf)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := builder.MakeTelegramBot(conf)
	if err != nil {
		log.Fatal(err)
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
		),
	)

	handlersContainer := containers.NewHandlersContainer(
		userHandler.NewTgbotUserHandler(servicesContainer.UserService, bot),
		trackingHandler.NewTgbotTrackingHandler(servicesContainer.TrackingService, bot),
	)

	return &app{
		conf,
		db,
		bot,
		reposContainer,
		servicesContainer,
		handlersContainer,
	}
}

func (app *app) Start() {
	updates, err := app.getBotUpdates()
	if err != nil {
		log.Fatal(err)
	}
	app.startCheckingUpdates()
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			app.handleCommand(&update)
		} else if update.CallbackQuery != nil {
			app.handleCallback(&update)
		}
	}
}

func (app *app) getBotUpdates() (tgbotapi.UpdatesChannel, error) {
	if !app.conf.UseWebhook {
		updateConfig := builder.MakeTelegramBotUpdateConfig()
		return app.bot.GetUpdatesChan(*updateConfig)
	}
	_, err := app.bot.SetWebhook(
		tgbotapi.NewWebhookWithCert(app.conf.HerokuBaseUrl+app.bot.Token, "cert.pem"),
	)
	if err != nil {
		return nil, err
	}

	updates := app.bot.ListenForWebhook("/" + app.bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
	return updates, nil
}

func (app *app) handleCommand(update *tgbotapi.Update) {
	command := update.Message.Command()
	if command == commands.Start {
		go app.handlersContainer.UserHandler.Join(update.Message.Chat.ID)
	} else if command == commands.AddTracking {
		commandArgs := update.Message.CommandArguments()
		go app.handlersContainer.TrackingHandler.AddTracking(commandArgs, update.Message.Chat.ID)
	} else if command == commands.GetAll {
		go app.handlersContainer.TrackingHandler.GetAll(update.Message.Chat.ID)
	}
}

func (app *app) handleCallback(update *tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	if trackingID, err := callbacks.ParseDeleteTrackingCallback(callbackData); err == nil {
		go app.handlersContainer.TrackingHandler.DeleteTracking(
			trackingID,
			update.CallbackQuery.Message.Chat.ID,
			int64(update.CallbackQuery.Message.MessageID),
		)
	}
}

func (app *app) startCheckingUpdates() {
	ticker := time.NewTicker(2 * time.Hour)
	go func() {
		for range ticker.C {
			updates, err := app.servicesContainer.TrackingService.CheckUpdates()
			if err == nil {
				for _, update := range updates {
					msg := tgbotapi.NewMessage(update.User.ChatID, fmt.Sprintf(
						texts.TrackingInfoUpdatedTempl,
						update.Tracking.Name,
						update.Tracking.Status,
						update.Tracking.Value,
					))
					msg.ParseMode = tgbotapi.ModeMarkdown
					app.bot.Send(msg)
				}
			}
		}
	}()
}
