package fetcher

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPFetcher struct {
	client *http.Client
}

func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (h *HTTPFetcher) Fetch(word string) (*http.Response, error) {
	url := fmt.Sprintf("https://dictionary.cambridge.org/dictionary/english/%s", word)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	return h.client.Do(req)
}
