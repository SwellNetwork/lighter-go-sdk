package lighter

import "time"

type HTTPClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		BaseURL: "https://mainnet.zklighter.elliot.ai",
		Timeout: 3 * time.Second,
	}
}
