package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/teimurjan/go-els-tg-bot/errs"
	"github.com/teimurjan/go-els-tg-bot/texts"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	utils "github.com/teimurjan/go-els-tg-bot/utils/regexp"
)

const elsURL = "https://els.kg/en/find_tracking"
const weightSelector = "div>div>span:first-of-type"
const statusSelector = "div>div>span:last-of-type"
const CSRFTokenSelector = "meta[name=\"csrf-token\"]"

type trackingDataFetcher struct{}

// NewTrackingDataFetcher creates a new instance of tracking status fetcher
func NewTrackingDataFetcher() tracking.TrackingDataFetcher {
	return &trackingDataFetcher{}
}

func strip(str string) string {
	return strings.Join(strings.Fields(str), " ")
}

// Fetch fetches oreder status by tracking
func (t *trackingDataFetcher) Fetch(trackingNumber string) (*tracking.TrackingData, error) {
	pageResponse, err := http.Get(elsURL)
	if err != nil {
		return nil, err
	}

	cookie, err := getCookie(pageResponse)
	if err != nil {
		return nil, err
	}

	CSRFToken, err := getCSRFToken(pageResponse)
	if err != nil {
		return nil, err
	}

	request, err := NewRequestBuilder(elsURL+getQuery(trackingNumber), CSRFToken, cookie).Build()
	if err != nil {
		return nil, err
	}

	body, err := sendRequest(request)
	if err != nil {
		return nil, err
	}

	doc, err := getHTMLFromResponse(string(body))
	if err != nil {
		return nil, err
	}

	status := doc.Find(statusSelector).Text()
	if status == "" {
		return nil, errs.NewErr(
			errs.NoSuchTrackingErrCode,
			texts.GetTrackingNotExistsMessage(trackingNumber),
		)
	}

	weight := doc.Find(weightSelector).Text()
	if weight == "" {
		weight = "Unknown"
	}

	return &tracking.TrackingData{Status: strip(status), Weight: strip(weight)}, nil
}

func sendRequest(request *http.Request) (string, error) {
	client := &http.Client{}
	statusResponse, err := client.Do(request)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(statusResponse.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getCookie(res *http.Response) (string, error) {
	for k, v := range res.Header {
		if k == "Set-Cookie" {
			return v[0], nil
		}
	}
	return "", fmt.Errorf("Can't get cookie")
}

func getCSRFToken(res *http.Response) (string, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	regexpGroups := utils.GetGroups(
		`<meta name="csrf-token" content="(?P<CSRFToken>.*)" />`,
		string(body),
	)
	CSRFToken, ok := regexpGroups["CSRFToken"]

	if ok {
		return CSRFToken, nil
	}

	return "", fmt.Errorf("Can't get CSRFToken")
}

func getQuery(trackingNumber string) string {
	return "?" + url.PathEscape(
		fmt.Sprintf(
			"utf8=✓&q=%s&commit=Search",
			trackingNumber,
		),
	)
}

func getHTMLFromResponse(response string) (*goquery.Document, error) {
	regexpGroups := utils.GetGroups(`^.*\.innerHTML\s*=\s*"(?P<html>.*)"`, response)

	html, ok := regexpGroups["html"]
	if ok {
		replacer := strings.NewReplacer("\\n", "", "\\", "")
		return goquery.NewDocumentFromReader(
			strings.NewReader(replacer.Replace(html)),
		)
	}

	return nil, fmt.Errorf("Can't get HTML from response %s", response)
}
