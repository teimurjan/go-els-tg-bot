package fetcher

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/teimurjan/go-els-tg-bot/tracking"
	errsUtil "github.com/teimurjan/go-els-tg-bot/utils/errs"
)

const apiURL = "https://els.kg/api/search"

type anonymouseTrackingDataFetcher struct{}

// NewAnonymouseFetcher creates a new instance of anonymous tracking status fetcher
func NewAnonymouseFetcher() tracking.TrackingDataFetcher {
	return &anonymouseTrackingDataFetcher{}
}

// Fetch fetches order status by tracking
func (t *anonymouseTrackingDataFetcher) Fetch(trackingNumber string) (*tracking.TrackingData, error) {
	var data = []byte(`{"tracking":"` + trackingNumber + `"}`)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 422 {
		return nil, errsUtil.NewI18NErr("invalidTrackingNumberFormat")
	}
	if resp.StatusCode == 404 {
		return &tracking.TrackingData{
			Status: EmptyStatusValue,
			Weight: EmptyWeightValue,
		}, nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	trackingData := &tracking.TrackingData{}

	json.Unmarshal([]byte(string(body)), &trackingData)

	return trackingData, nil
}
