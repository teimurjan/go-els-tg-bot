package tracking

type TrackingData struct {
	Status string `json:"status_warehouse"`
	Weight string `json:"weight"`
}

type TrackingDataFetcher interface {
	Fetch(trackingNumber string) (*TrackingData, error)
}
