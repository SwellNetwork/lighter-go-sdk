package lighter

import "time"

type HTTPClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

func DefaultMainnetHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		BaseURL: "https://mainnet.zklighter.elliot.ai",
		Timeout: 3 * time.Second,
	}
}

func DefaultTestnetHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		BaseURL: "https://testnet.zklighter.elliot.ai",
		Timeout: 3 * time.Second,
	}
}
