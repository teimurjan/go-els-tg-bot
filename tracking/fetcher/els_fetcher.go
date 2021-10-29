package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	regexpUtil "github.com/teimurjan/go-els-tg-bot/utils/regexp"
)

const findTrackingURL = "https://dash.els.kg/find_tracking"
const loginURL = "https://dash.els.kg/login"
const weightSelector = "div>div>span:first-of-type"
const statusSelector = "div>div>span:last-of-type"

const EmptyStatusValue = "Not arrived at ELS yet. ⚠️"
const EmptyWeightValue = "Unknown"

type trackingDataFetcher struct {
	conf *config.Config
}

// NewTrackingDataFetcher creates a new instance of tracking status fetcher
func NewTrackingDataFetcher(conf *config.Config) tracking.TrackingDataFetcher {
	return &trackingDataFetcher{
		conf,
	}
}

var cachedRememberUserToken = ""

// Fetch fetches order status by tracking
func (t *trackingDataFetcher) Fetch(trackingNumber string) (*tracking.TrackingData, error) {
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

	trackingRes, err := fetchTrackingInfo(trackingNumber, csrfToken, appSessionCookie, cachedRememberUserToken)
	if err != nil {
		return nil, err
	}

	return responseToTrackingData(trackingRes)
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

func fetchTrackingInfo(trackingNumber string, csrfToken string, appSessionCookie string, rememberUserToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", findTrackingURL+getQuery(trackingNumber), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "_els_app_session", Value: appSessionCookie})
	req.AddCookie(&http.Cookie{Name: "signed_in", Value: "true"})
	req.AddCookie(&http.Cookie{Name: "remember_user_token", Value: rememberUserToken})
	req.Header.Add("X-CSRF-Token", csrfToken)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}

	return client.Do(req)
}

func responseToTrackingData(res *http.Response) (*tracking.TrackingData, error) {
	doc, err := parseHTMLResponse(res)
	if err != nil {
		return nil, err
	}

	status := doc.Find(statusSelector).Text()
	if status == "" {
		if len(doc.Nodes) == 1 && doc.Find("h3") != nil {
			return &tracking.TrackingData{
				Status: EmptyStatusValue,
				Weight: EmptyWeightValue,
			}, nil
		}

		return nil, errors.New("can't extract any valuable infromation from the response")
	}

	weight := doc.Find(weightSelector).Text()
	if weight == "" {
		weight = EmptyWeightValue
	}

	return &tracking.TrackingData{Status: strip(status), Weight: strip(weight)}, nil
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
