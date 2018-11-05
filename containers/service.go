package containers

import (
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type ServicesContainer struct {
	UserService     user.UserService
	TrackingService tracking.TrackingService
}

func NewServicesContainer(
	userService user.UserService,
	trackingService tracking.TrackingService,
) *ServicesContainer {
	return &ServicesContainer{
		userService,
		trackingService,
	}
}
