package app

import (
	"net/http"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"

	"github.com/teimurjan/go-els-tg-bot/commands"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	botFactory "github.com/teimurjan/go-els-tg-bot/factory/bot"
	containersFactory "github.com/teimurjan/go-els-tg-bot/factory/containers"
	helper "github.com/teimurjan/go-els-tg-bot/helper/i18n"
	callbacksUtil "github.com/teimurjan/go-els-tg-bot/utils/callbacks"
)

type TgBotApp interface {
	Start()
}

type tgBotApp struct {
	conf              *config.Config
	logger            *logrus.Logger
	bot               *tgbotapi.BotAPI
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

// NewTgBotApp creates new tg bot application
func NewTgBotApp(conf *config.Config, db *sqlx.DB, logger *logrus.Logger) TgBotApp {
	bot, err := botFactory.MakeTelegramBot(conf)
	if err != nil {
		logger.Fatal("Can't create a telegram bot.", err)
	}

	reposContainer := containersFactory.MakeReposContainer(db)
	i18nHelper := helper.NewI18nHelper(reposContainer.UserRepo)
	servicesContainer := containersFactory.MakeServicesContainer(reposContainer, logger, conf)
	handlersContainer := containersFactory.MakeHandlersContainer(servicesContainer, bot, i18nHelper)

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
			go func() {
				tgBotApp.resetAllDialogs(&update)
				tgBotApp.handleCommand(&update)
			}()
		} else if update.CallbackQuery != nil {
			go tgBotApp.handleCallback(&update)
		} else if update.Message != nil && len(update.Message.Text) > 0 {
			go tgBotApp.handleText(&update)
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
		tgBotApp.handlersContainer.UserHandler.Join(update.Message.Chat.ID)
	} else if command == commands.AddTracking {
		commandArgs := update.Message.CommandArguments()
		if len(commandArgs) > 0 {
			tgBotApp.handlersContainer.TrackingHandler.AddTracking(commandArgs, update.Message.Chat.ID)
		} else {
			tgBotApp.handlersContainer.AddTrackingDialogHandler.StartDialog(update.Message.Chat.ID)
		}
	} else if command == commands.GetAll {
		tgBotApp.handlersContainer.TrackingHandler.GetAll(update.Message.Chat.ID)
	} else if command == commands.ChangeLanguage {
		tgBotApp.handlersContainer.UserHandler.RequestLanguageChange(update.Message.Chat.ID)
	} else if command == commands.GetUsaAddress {
		tgBotApp.handlersContainer.UsaAddressHandler.GetAddress(update.Message.Chat.ID)
	}
}

func (tgBotApp *tgBotApp) handleCallback(update *tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	if trackingID, err := callbacksUtil.ParseDeleteTrackingCallback(callbackData); err == nil {
		tgBotApp.handlersContainer.TrackingHandler.DeleteTracking(
			trackingID,
			update.CallbackQuery.Message.Chat.ID,
			int64(update.CallbackQuery.Message.MessageID),
		)
	} else if language, err := callbacksUtil.ParseChangeLanguageCallback(callbackData); err == nil {
		tgBotApp.handlersContainer.UserHandler.ChangeLanguage(
			language,
			update.CallbackQuery.Message.Chat.ID,
			int64(update.CallbackQuery.Message.MessageID),
		)
	}
}

func (tgBotApp *tgBotApp) handleText(update *tgbotapi.Update) {
	text := update.Message.Text
	tgBotApp.handlersContainer.AddTrackingDialogHandler.UpdateDialogIfActive(
		text,
		update.Message.Chat.ID,
	)
}

func (tgBotApp *tgBotApp) resetAllDialogs(update *tgbotapi.Update) {
	tgBotApp.handlersContainer.AddTrackingDialogHandler.ResetDialog(
		update.Message.Chat.ID,
	)
}
