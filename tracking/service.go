package tracking

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type TrackingUpdate struct {
	User     *models.User
	Tracking *models.Tracking
}

func NewTrackingUpdate(user *models.User, tracking *models.Tracking) *TrackingUpdate {
	return &TrackingUpdate{
		user,
		tracking,
	}
}

type TrackingService interface {
	Create(tracking string, name string, chatID int64) (*models.Tracking, error)
	GetAll(chatID int64) ([]*models.Tracking, error)
	GetUpdates() ([]*TrackingUpdate, error)
	Delete(trackingID int64) error
}
