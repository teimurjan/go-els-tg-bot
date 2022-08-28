package job

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/containers"
	botFactory "github.com/teimurjan/go-els-tg-bot/factory/bot"
	containersFactory "github.com/teimurjan/go-els-tg-bot/factory/containers"
	helper "github.com/teimurjan/go-els-tg-bot/helper/i18n"
	"github.com/teimurjan/go-els-tg-bot/job"
)

type tgBotJob struct {
	reposContainer    *containers.RepositoriesContainer
	servicesContainer *containers.ServicesContainer
	handlersContainer *containers.HandlersContainer
}

// NewTgBotJob creates a job for telegram bot
func NewTgBotJob(conf *config.Config, db *sqlx.DB, logger *logrus.Logger) job.Job {
	bot, err := botFactory.MakeTelegramBot(conf)
	if err != nil {
		logger.Fatal("Can't create a telegram bot.", err)
	}

	reposContainer := containersFactory.MakeReposContainer(db)
	i18nHelper := helper.NewI18nHelper(reposContainer.UserRepo, logger)
	servicesContainer := containersFactory.MakeServicesContainer(reposContainer, logger, conf)
	handlersContainer := containersFactory.MakeHandlersContainer(servicesContainer, bot, i18nHelper)

	return &tgBotJob{
		reposContainer,
		servicesContainer,
		handlersContainer,
	}
}

// Do executes TgBotJob
func (tgBotJob *tgBotJob) Do() {
	tgBotJob.handlersContainer.TrackingHandler.CheckUpdates()
}
