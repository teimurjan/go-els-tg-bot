package usaAddress

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type UsaAddressService interface {
	GetUpdates() ([]string, *models.UsaAddress, error)
	GetFirst() (*models.UsaAddress, error)
	Delete() error
}
