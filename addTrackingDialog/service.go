package addTrackingDialog

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type AddTrackingDialogService interface {
	GetDialogForChat(chatID int64) (*models.AddTrackingDialog, error)
	StartDialog(userID int64) (*models.AddTrackingDialog, error)
	UpdateDialogName(dialog *models.AddTrackingDialog, name string) error
	UpdateDialogTracking(dialog *models.AddTrackingDialog, tracking string) error
	ResetDialog(dialog *models.AddTrackingDialog) error
}
