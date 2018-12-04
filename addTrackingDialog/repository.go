package addTrackingDialog

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type AddTrackingDialogRepository interface {
	GetForUser(userId int64) (*models.AddTrackingDialog, error)
	Store(t *models.AddTrackingDialog) (int64, error)
	Update(t *models.AddTrackingDialog) error
}
