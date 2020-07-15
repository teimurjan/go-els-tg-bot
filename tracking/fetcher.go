package tracking

type TrackingStatus struct {
	Status string
	Weight string
}

type TrackingStatusFetcher interface {
	Fetch(trackingNumber string) (*TrackingStatus, error)
}
