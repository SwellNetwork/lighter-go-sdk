package lighter

import "time"

type WSClientConfig struct {
	BaseURL           string
	PingInterval      time.Duration
	ReconnectAttempts int
	AuthToken         string
}

func DefaultMainnetWSClientConfig() WSClientConfig {
	return WSClientConfig{
		BaseURL:           "wss://mainnet.zklighter.elliot.ai/stream",
		PingInterval:      30 * time.Second,
		ReconnectAttempts: 5,
	}
}

func DefaultTestnetWSClientConfig() WSClientConfig {
	return WSClientConfig{
		BaseURL:           "wss://testnet.zklighter.elliot.ai/stream",
		PingInterval:      30 * time.Second,
		ReconnectAttempts: 5,
	}
}
