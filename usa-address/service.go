package usaAddress

import (
	"github.com/r3labs/diff"
	"github.com/teimurjan/go-els-tg-bot/models"
)

type AddressUpdate struct {
	FieldName  string
	FieldValue string
}

type UsaAddressService interface {
	GetAddressWithDiff() (*models.UsaAddress, diff.Changelog, error)
	GetFirst() (*models.UsaAddress, error)
	Delete() error
}
