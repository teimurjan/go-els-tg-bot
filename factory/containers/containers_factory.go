package factory

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-els-tg-bot/containers"
	trackingFetcher "github.com/teimurjan/go-els-tg-bot/tracking/fetcher"
	trackingHandler "github.com/teimurjan/go-els-tg-bot/tracking/handler"
	trackingRepository "github.com/teimurjan/go-els-tg-bot/tracking/repository"
	trackingService "github.com/teimurjan/go-els-tg-bot/tracking/service"
	userHandler "github.com/teimurjan/go-els-tg-bot/user/handler"
	userRepository "github.com/teimurjan/go-els-tg-bot/user/repository"
	userService "github.com/teimurjan/go-els-tg-bot/user/service"
)

// MakeReposContainer creates container with repositories
func MakeReposContainer(db *sqlx.DB) *containers.RepositoriesContainer {
	return containers.NewRepositoriesContainer(
		userRepository.NewPostgresqlUserRepository(db),
		trackingRepository.NewPostgresqlTrackingRepository(db),
	)
}

// MakeServicesContainer creates container with services
func MakeServicesContainer(
	repos *containers.RepositoriesContainer,
	logger *logrus.Logger,
) *containers.ServicesContainer {
	return containers.NewServicesContainer(
		userService.NewUserService(
			repos.UserRepo,
			logger,
		),
		trackingService.NewTrackingService(
			repos.TrackingRepo,
			repos.UserRepo,
			trackingFetcher.NewTrackingStatusFetcher(),
			logger,
		),
	)
}

// MakeHandlersContainer creates container with handlers
func MakeHandlersContainer(
	services *containers.ServicesContainer,
	bot *tgbotapi.BotAPI,
) *containers.HandlersContainer {
	return containers.NewHandlersContainer(
		userHandler.NewTgbotUserHandler(services.UserService, bot),
		trackingHandler.NewTgbotTrackingHandler(services.TrackingService, bot),
	)
}
