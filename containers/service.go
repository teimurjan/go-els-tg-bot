package containers

import (
	addTrackingDialog "github.com/teimurjan/go-els-tg-bot/add-tracking-dialog"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type ServicesContainer struct {
	UserService              user.UserService
	TrackingService          tracking.TrackingService
	AddTrackingDialogService addTrackingDialog.AddTrackingDialogService
}

func NewServicesContainer(
	userService user.UserService,
	trackingService tracking.TrackingService,
	addTrackingDialogService addTrackingDialog.AddTrackingDialogService,
) *ServicesContainer {
	return &ServicesContainer{
		userService,
		trackingService,
		addTrackingDialogService,
	}
}
