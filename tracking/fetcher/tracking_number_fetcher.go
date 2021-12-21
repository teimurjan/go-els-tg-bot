package fetcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	regexpUtil "github.com/teimurjan/go-els-tg-bot/utils/regexp"
)

const findTrackingURL = "https://dash.els.kg/find_tracking"
const loginURL = "https://dash.els.kg/login"
const weightSelector = "div>div>span:first-of-type"
const statusSelector = "div>div>span:last-of-type"

const EmptyStatusValue = "Not arrived at ELS yet. ⚠️"
const EmptyWeightValue = "Unknown"

type trackingNumberFetcher struct {
	conf *config.Config
}

type Ship24ResponseData struct {
	Data struct {
		ID              string `json:"_id"`
		TrackingNumber  string `json:"tracking_number"`
		TrackingNumbers []struct {
			IsMain         bool     `json:"is_main"`
			TrackingNumber string   `json:"tracking_number"`
			CrawlerCodes   []string `json:"crawler_codes"`
		} `json:"tracking_numbers"`
		ParcelIdentifier       string  `json:"parcel_identifier"`
		DestinationCountryCode string  `json:"destination_country_code"`
		Weight                 float64 `json:"weight"`
		RecipientCity          string  `json:"recipient_city"`
		Events                 []struct {
			Datetime       time.Time `json:"datetime"`
			DispatchCodeID int       `json:"dispatch_code_id,omitempty"`
			Status         string    `json:"status"`
			Location       string    `json:"location,omitempty"`
			Courier        struct {
				Slug        string      `json:"slug"`
				CrawlerCode string      `json:"crawler_code"`
				CountryCode interface{} `json:"country_code"`
				LogoImage   struct {
					Path string `json:"path"`
				} `json:"logo_image"`
				Translation struct {
					LangCode string      `json:"lang_code"`
					Name     string      `json:"name"`
					FullName string      `json:"full_name"`
					Website  string      `json:"website"`
					Phone    interface{} `json:"phone"`
				} `json:"translation"`
			} `json:"courier"`
		} `json:"events"`
		DispatchCode struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Step string `json:"step"`
			Desc string `json:"desc"`
		} `json:"dispatch_code"`
		Couriers []struct {
			Slug        string      `json:"slug"`
			CrawlerCode string      `json:"crawler_code"`
			CountryCode interface{} `json:"country_code"`
			LogoImage   struct {
				Path string `json:"path"`
			} `json:"logo_image"`
			Translation struct {
				LangCode string      `json:"lang_code"`
				Name     string      `json:"name"`
				FullName string      `json:"full_name"`
				Website  string      `json:"website"`
				Phone    interface{} `json:"phone"`
			} `json:"translation"`
		} `json:"couriers"`
	} `json:"data"`
}

// NewTrackingNumberFetcher creates a new instance of tracking status fetcher
func NewTrackingNumberFetcher(conf *config.Config) tracking.TrackingNumberFetcher {
	return &trackingNumberFetcher{
		conf,
	}
}

var cachedRememberUserToken = ""

// Fetch fetches order status by tracking
func (t *trackingNumberFetcher) Fetch(trackingNumber string) (*models.Tracking, error) {
	res, err := http.Get(loginURL)
	if err != nil {
		return nil, err
	}

	appSessionCookie, err := getCookie(res, "_els_app_session")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	authencityToken, err := getAuthencityToken(body)
	if err != nil {
		return nil, err
	}

	csrfToken, err := getCSRFToken(body)
	if err != nil {
		return nil, err
	}

	if len(cachedRememberUserToken) == 0 {
		authorizationRes, err := authorize(t.conf.ElsUserEmail, t.conf.ElsUserPassword, appSessionCookie, authencityToken)
		if err != nil {
			return nil, err
		}

		appSessionCookie, err = getCookie(authorizationRes, "_els_app_session")
		if err != nil {
			return nil, err
		}

		cachedRememberUserToken, err = getCookie(authorizationRes, "remember_user_token")
		if err != nil {
			return nil, err
		}
	}

	elsResponse, err := fetchTrackingAtEls(trackingNumber, csrfToken, appSessionCookie, cachedRememberUserToken)
	if err != nil {
		return nil, err
	}

	tracking, err := extractTrackingFromElsResponse(elsResponse)
	// There is no tracking extracted from ELS response
	if err != nil {
		// Try to get info from ship24
		ship24Response, err := fetchTrackingAtShip24(trackingNumber)
		// No tracking info anywhere
		if err != nil || ship24Response.StatusCode > 201 {
			return &models.Tracking{
				Status: EmptyStatusValue,
				Weight: EmptyWeightValue,
			}, nil
		}

		ship24Body, err := ioutil.ReadAll(ship24Response.Body)
		if err != nil {
			return nil, err
		}
		var ship24ResponseData Ship24ResponseData
		err = json.Unmarshal(ship24Body, &ship24ResponseData)
		if err != nil {
			return nil, err
		}
		// No events at ship24
		if len(ship24ResponseData.Data.Events) == 0 {
			return &models.Tracking{
				Status: EmptyStatusValue,
				Weight: EmptyWeightValue,
			}, nil
		}

		lastEvent := ship24ResponseData.Data.Events[0]
		y, m, d := lastEvent.Datetime.Date()
		return &models.Tracking{
			Status: lastEvent.Status + " at " + lastEvent.Location + " (" + fmt.Sprint(d, "-", int(m), "-", y) + ")",
			Weight: strconv.Itoa(int(ship24ResponseData.Data.Weight)),
		}, nil
	}

	return tracking, nil
}

func getCookie(res *http.Response, name string) (string, error) {
	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}

	return "", errors.New("App session cookie is not found")
}

func getAuthencityToken(body []byte) (string, error) {
	regexpGroups := regexpUtil.GetGroups(
		`<input type="hidden" name="authenticity_token" value="(?P<token>.*)" />`,
		string(body),
	)
	token, ok := regexpGroups["token"]

	if ok {
		return token, nil
	}

	return "", fmt.Errorf("Can't get authencity token")
}

func getCSRFToken(body []byte) (string, error) {
	regexpGroups := regexpUtil.GetGroups(
		`<meta name="csrf-token" content="(?P<token>.*)" />`,
		string(body),
	)
	token, ok := regexpGroups["token"]

	if ok {
		return token, nil
	}

	return "", fmt.Errorf("Can't get CSRF token")
}

func authorize(email string, password string, appSessionCookie string, authencityToken string) (*http.Response, error) {
	values := url.Values{
		"user[email]":        {email},
		"user[password]":     {password},
		"authenticity_token": {authencityToken},
		"commit":             {"Войти"},
	}
	payload := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", loginURL, payload)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "_els_app_session", Value: appSessionCookie})

	client := &http.Client{
		// Do not follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return client.Do(req)
}

func fetchTrackingAtEls(trackingNumber string, csrfToken string, appSessionCookie string, rememberUserToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", findTrackingURL+getQuery(trackingNumber), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "_els_app_session", Value: appSessionCookie})
	req.AddCookie(&http.Cookie{Name: "signed_in", Value: "true"})
	req.AddCookie(&http.Cookie{Name: "remember_user_token", Value: rememberUserToken})
	req.Header.Add("X-CSRF-Token", csrfToken)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	return http.DefaultClient.Do(req)
}

func fetchTrackingAtShip24(trackingNumber string) (*http.Response, error) {
	payload := []byte("{\"userAgent\":\"\",\"os\":\"Mac\",\"browser\":\"Chrome\",\"device\":\"Unknown\",\"os_version\":\"mac-os-x-15\",\"browser_version\":\"96.0.4664.55\",\"uL\":\"en-US\"}")
	body := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", "https://api.ship24.com/api/parcels/update/"+trackingNumber+"?lang=en", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authority", "api.ship24.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Origin", "https://www.ship24.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.ship24.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8,nl;q=0.7")
	req.Header.Set("X-Ship24-Token", "225,150,128,44,49,54,52,48,48,54,54,52,49,54,54,48,48,44,225,150,130")

	return http.DefaultClient.Do(req)
}

func extractTrackingFromElsResponse(res *http.Response) (*models.Tracking, error) {
	doc, err := parseHTMLResponse(res)
	if err != nil {
		return nil, err
	}

	status := doc.Find(statusSelector).Text()
	if status == "" {
		if len(doc.Nodes) == 1 && doc.Find("h3") != nil {
			return nil, errors.New("no tracking in the response")
		}

		return nil, errors.New("can't extract any valuable infromation from the response")
	}

	weight := doc.Find(weightSelector).Text()
	if weight == "" {
		weight = EmptyWeightValue
	}

	return &models.Tracking{Status: strip(status), Weight: strip(weight)}, nil
}

func strip(str string) string {
	return strings.Join(strings.Fields(str), " ")
}

func getQuery(trackingNumber string) string {
	return "?" + url.PathEscape(
		fmt.Sprintf(
			"q=%s&commit=Поиск",
			trackingNumber,
		),
	)
}

func parseHTMLResponse(res *http.Response) (*goquery.Document, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	regexpGroups := regexpUtil.GetGroups(`^.*\.innerHTML\s*=\s*"(?P<html>.*)"`, string(body))

	html, ok := regexpGroups["html"]
	if ok {
		replacer := strings.NewReplacer("\\n", "", "\\", "")
		return goquery.NewDocumentFromReader(
			strings.NewReader(replacer.Replace(html)),
		)
	}

	return nil, fmt.Errorf("Can't get HTML from response %s", body)
}
