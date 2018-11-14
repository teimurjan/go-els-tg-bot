package tracking

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type TrackingRepository interface {
	GetByID(id int64) (*models.Tracking, error)
	GetForUser(userId int64) ([]*models.Tracking, error)
	Store(t *models.Tracking) (int64, error)
	Update(t *models.Tracking) error
	Delete(trackingID int64) error
}
