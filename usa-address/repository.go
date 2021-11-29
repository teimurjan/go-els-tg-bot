package usaAddress

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type UsaAddressRepository interface {
	GetFirst() (*models.UsaAddress, error)
	Store(address *models.UsaAddress) (int64, error)
	Update(address *models.UsaAddress) error
	Delete() error
}
