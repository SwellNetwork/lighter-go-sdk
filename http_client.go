package lighter

import "resty.dev/v3"

type HTTPClient struct {
	config HTTPClientConfig

	client *resty.Client
}

func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		config: config,
		client: resty.New().SetBaseURL(config.BaseURL).SetTimeout(config.Timeout),
	}
}
