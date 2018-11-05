package containers

import (
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type HandlersContainer struct {
	UserHandler     user.UserHandler
	TrackingHandler tracking.TrackingHandler
}

func NewHandlersContainer(
	userHandler user.UserHandler,
	trackingHandler tracking.TrackingHandler,
) *HandlersContainer {
	return &HandlersContainer{
		userHandler,
		trackingHandler,
	}
}
