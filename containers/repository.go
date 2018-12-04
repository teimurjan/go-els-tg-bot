package containers

import (
	"github.com/teimurjan/go-els-tg-bot/addTrackingDialog"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type RepositoriesContainer struct {
	UserRepo              user.UserRepository
	TrackingRepo          tracking.TrackingRepository
	AddTrackingDialogRepo addTrackingDialog.AddTrackingDialogRepository
}

func NewRepositoriesContainer(
	userRepository user.UserRepository,
	trackingRepository tracking.TrackingRepository,
	addTrackingDialog addTrackingDialog.AddTrackingDialogRepository,
) *RepositoriesContainer {
	return &RepositoriesContainer{
		userRepository,
		trackingRepository,
		addTrackingDialog,
	}
}
