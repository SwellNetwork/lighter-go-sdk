package lighter

import "time"

type WSClientConfig struct {
	BaseURL                  string
	PingInterval             time.Duration
	ReconnectAttempts        int
	AuthToken                string
	SubscribeRetryAttempts   int
	SubscribeRetryMinBackoff time.Duration
	SubscribeRetryMaxBackoff time.Duration
}

func DefaultMainnetWSClientConfig() WSClientConfig {
	return WSClientConfig{
		BaseURL:                  "wss://mainnet.zklighter.elliot.ai/stream",
		PingInterval:             30 * time.Second,
		ReconnectAttempts:        5,
		SubscribeRetryAttempts:   0,
		SubscribeRetryMinBackoff: time.Second,
		SubscribeRetryMaxBackoff: 30 * time.Second,
	}
}

func DefaultTestnetWSClientConfig() WSClientConfig {
	return WSClientConfig{
		BaseURL:                  "wss://testnet.zklighter.elliot.ai/stream",
		PingInterval:             30 * time.Second,
		ReconnectAttempts:        5,
		SubscribeRetryAttempts:   0,
		SubscribeRetryMinBackoff: time.Second,
		SubscribeRetryMaxBackoff: 30 * time.Second,
	}
}
