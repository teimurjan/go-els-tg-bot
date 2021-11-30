package tracking

import "github.com/teimurjan/go-els-tg-bot/models"

type TrackingNumberFetcher interface {
	Fetch(trackingNumber string) (*models.Tracking, error)
}
