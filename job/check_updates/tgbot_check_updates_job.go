package job

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	botFactory "github.com/teimurjan/go-els-tg-bot/factory/bot"
	containersFactory "github.com/teimurjan/go-els-tg-bot/factory/containers"
)

type tgBotJob struct {
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

// NewTgBotJob creates a job for telegram bot
func NewTgBotJob(conf *config.Config, db *sqlx.DB, logger *logrus.Logger) *tgBotJob {
	bot, err := botFactory.MakeTelegramBot(conf)
	if err != nil {
		logger.Fatal("Can't create a telegram bot.", err)
	}

	reposContainer := containersFactory.MakeReposContainer(db)
	servicesContainer := containersFactory.MakeServicesContainer(reposContainer, logger)
	handlersContainer := containersFactory.MakeHandlersContainer(servicesContainer, bot)
	return &tgBotJob{
		reposContainer,
		servicesContainer,
		handlersContainer,
	}
}

func (tgBotJob *tgBotJob) Do() {
	tgBotJob.handlersContainer.TrackingHandler.CheckUpdates()
}
