package tibia

import (
	"net/http"
	"time"
)

type ApiClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewApiClient() *ApiClient {
	return &ApiClient{
		httpClient: newHttpClient(),
		baseUrl:    "https://api.tibiadata.com",
	}
}

func newHttpClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 20
	transport.MaxConnsPerHost = 20
	transport.MaxIdleConnsPerHost = 20

	return &http.Client{
		Timeout:   time.Second * 20,
		Transport: transport,
	}
}
