package legacyFetcher

import "net/http"

// RequestBuilder is a request builder interface
type RequestBuilder interface {
	Build() (*http.Request, error)
}

type requestBuilder struct {
	URL       string
	CSRFToken string
	cookie    string
}

// NewRequestBuilder creates requestBuilder instance
func NewRequestBuilder(
	URL string,
	CSRFToken string,
	cookie string,
) RequestBuilder {
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
