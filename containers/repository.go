package containers

import (
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type RepositoriesContainer struct {
	UserRepo     user.UserRepository
	TrackingRepo tracking.TrackingRepository
}

func NewRepositoriesContainer(
	userRepository user.UserRepository,
	trackingRepository tracking.TrackingRepository,
) *RepositoriesContainer {
	return &RepositoriesContainer{
		userRepository,
		trackingRepository,
	}
}
