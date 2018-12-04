package containers

import (
	"github.com/teimurjan/go-els-tg-bot/addTrackingDialog"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type HandlersContainer struct {
	UserHandler              user.UserHandler
	TrackingHandler          tracking.TrackingHandler
	AddTrackingDialogHandler addTrackingDialog.AddTrackingDialogHandler
}

func NewHandlersContainer(
	userHandler user.UserHandler,
	trackingHandler tracking.TrackingHandler,
	addTrackingDialogHandler addTrackingDialog.AddTrackingDialogHandler,
) *HandlersContainer {
	return &HandlersContainer{
		userHandler,
		trackingHandler,
		addTrackingDialogHandler,
	}
}
