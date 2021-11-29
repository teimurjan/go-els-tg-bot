package fetcher

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/teimurjan/go-els-tg-bot/config"
	"github.com/teimurjan/go-els-tg-bot/models"
	usaAddress "github.com/teimurjan/go-els-tg-bot/usa-address"
)

const usaAddressURL = "https://els.kg/address-usa"

type usaAddressFetcher struct {
	conf *config.Config
}

// NewUsaAddressFetcher creates a new instance of usaAddress status fetcher
func NewUsaAddressFetcher(conf *config.Config) usaAddress.UsaAddressFetcher {
	return &usaAddressFetcher{
		conf,
	}
}

// Fetch fetches USA address
func (t *usaAddressFetcher) Fetch() (*models.UsaAddress, error) {
	usaAddressRes, err := fetchUSAAdress()
	if err != nil {
		return nil, err
	}

	return responseToUsaAddress(usaAddressRes)
}

func fetchUSAAdress() (*http.Response, error) {
	req, err := http.NewRequest("GET", usaAddressURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}

	return client.Do(req)
}

func responseToUsaAddress(res *http.Response) (*models.UsaAddress, error) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	address := &models.UsaAddress{}

	labels := doc.Find("dl > dt").Nodes
	values := doc.Find("dl > dd").Nodes

	for i, label := range labels {
		value := values[i]
		labelText := strings.ToLower(label.FirstChild.Data)
		valueText := value.FirstChild.Data

		if strings.Contains(labelText, "street address") {
			address.Street = valueText
		} else if strings.Contains(labelText, "city") {
			address.City = valueText
		} else if strings.Contains(labelText, "state") {
			address.State = valueText
		} else if strings.Contains(labelText, "zip") {
			address.Zip = valueText
		} else if strings.Contains(labelText, "phone") {
			address.PhoneNumber = valueText
		}
	}

	return address, nil
}
