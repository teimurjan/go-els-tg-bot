package usaAddress

import "github.com/teimurjan/go-els-tg-bot/models"

type UsaAddressFetcher interface {
	Fetch() (*models.UsaAddress, error)
}
