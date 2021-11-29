package containers

import (
	addTrackingDialog "github.com/teimurjan/go-els-tg-bot/add-tracking-dialog"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	usaAddress "github.com/teimurjan/go-els-tg-bot/usa-address"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type ServicesContainer struct {
	UserService              user.UserService
	TrackingService          tracking.TrackingService
	AddTrackingDialogService addTrackingDialog.AddTrackingDialogService
	UsaAddressService        usaAddress.UsaAddressService
}

func NewServicesContainer(
	userService user.UserService,
	trackingService tracking.TrackingService,
	addTrackingDialogService addTrackingDialog.AddTrackingDialogService,
	usaAddressService usaAddress.UsaAddressService,
) *ServicesContainer {
	return &ServicesContainer{
		userService,
		trackingService,
		addTrackingDialogService,
		usaAddressService,
	}
}
