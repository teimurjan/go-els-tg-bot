package tracking

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type TrackingService interface {
	Create(tracking string, name string, chatID int64) (*models.Tracking, error)
	GetForChat(chatID int64) ([]*models.Tracking, error)
	Update(tracking *models.Tracking) (bool, error)
	Delete(trackingID int64) error
	SyncAll(trackings []*models.Tracking) (chan *models.Tracking, chan error)
	SyncOnlyUpdated(trackings []*models.Tracking) (chan *models.Tracking, chan error)
	GetAllGroupedByUser() (map[*models.User][]*models.Tracking, error)
}
