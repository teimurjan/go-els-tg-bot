package tracking

type TrackingStatusFetcher interface {
	Fetch(trackingNumber string) (string, error)
}
