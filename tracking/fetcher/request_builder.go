package fetcher

import "net/http"

type requestBuilder struct {
	URL       string
	CSRFToken string
	cookie    string
}

func NewRequestBuilder(
	URL string,
	CSRFToken string,
	cookie string,
) *requestBuilder {
	return &requestBuilder{
		URL,
		CSRFToken,
		cookie,
	}
}

func (rb *requestBuilder) Build() (*http.Request, error) {
	request, err := http.NewRequest("GET", rb.URL, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-CSRF-Token", rb.CSRFToken)
	request.Header.Add("Cookie", rb.cookie)
	request.Header.Add("Host", "els.kg")
	request.Header.Add("Referer", "https://els.kg/en/find_tracking")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")

	return request, nil
}
