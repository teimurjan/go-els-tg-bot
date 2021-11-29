package containers

import (
	addTrackingDialog "github.com/teimurjan/go-els-tg-bot/add-tracking-dialog"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	usaAddress "github.com/teimurjan/go-els-tg-bot/usa-address"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type RepositoriesContainer struct {
	UserRepo              user.UserRepository
	TrackingRepo          tracking.TrackingRepository
	AddTrackingDialogRepo addTrackingDialog.AddTrackingDialogRepository
	UsaAddressRepo        usaAddress.UsaAddressRepository
}

func NewRepositoriesContainer(
	userRepository user.UserRepository,
	trackingRepository tracking.TrackingRepository,
	addTrackingDialogRepository addTrackingDialog.AddTrackingDialogRepository,
	usaAddressRepository usaAddress.UsaAddressRepository,
) *RepositoriesContainer {
	return &RepositoriesContainer{
		userRepository,
		trackingRepository,
		addTrackingDialogRepository,
		usaAddressRepository,
	}
}
