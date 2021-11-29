package containers

import (
	addTrackingDialog "github.com/teimurjan/go-els-tg-bot/add-tracking-dialog"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	usaAddress "github.com/teimurjan/go-els-tg-bot/usa-address"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type HandlersContainer struct {
	UserHandler              user.UserHandler
	TrackingHandler          tracking.TrackingHandler
	AddTrackingDialogHandler addTrackingDialog.AddTrackingDialogHandler
	UsaAddressHandler        usaAddress.UsaAddressHandler
}

func NewHandlersContainer(
	userHandler user.UserHandler,
	trackingHandler tracking.TrackingHandler,
	addTrackingDialogHandler addTrackingDialog.AddTrackingDialogHandler,
	usaAddressHandler usaAddress.UsaAddressHandler,
) *HandlersContainer {
	return &HandlersContainer{
		userHandler,
		trackingHandler,
		addTrackingDialogHandler,
		usaAddressHandler,
	}
}
