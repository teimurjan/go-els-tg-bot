package tracking

type TrackingData struct {
	Status string
	Weight string
}

type TrackingDataFetcher interface {
	Fetch(trackingNumber string) (*TrackingData, error)
}
