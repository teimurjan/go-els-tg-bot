package worker

import (
	"github.com/sirupsen/logrus"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"

	builder "github.com/teimurjan/go-els-tg-bot/builder/bot"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	trackingFetcher "github.com/teimurjan/go-els-tg-bot/tracking/fetcher"
	trackingHandler "github.com/teimurjan/go-els-tg-bot/tracking/handler"
	trackingRepository "github.com/teimurjan/go-els-tg-bot/tracking/repository"
	trackingService "github.com/teimurjan/go-els-tg-bot/tracking/service"
	userHandler "github.com/teimurjan/go-els-tg-bot/user/handler"
	userRepository "github.com/teimurjan/go-els-tg-bot/user/repository"
	userService "github.com/teimurjan/go-els-tg-bot/user/service"
)

type tgBotWorker struct {
	conf              *config.Config
	logger            *logrus.Logger
	db                *sqlx.DB
	bot               *tgbotapi.BotAPI
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

func NewTgBotWorker(conf *config.Config, db *sqlx.DB, logger *logrus.Logger) *tgBotWorker {
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

	return &tgBotWorker{
		conf,
		logger,
		db,
		bot,
		reposContainer,
		servicesContainer,
		handlersContainer,
	}
}

func (tgBotWorker *tgBotWorker) Do() {
	tgBotWorker.handlersContainer.TrackingHandler.CheckUpdates()
}
