package lighter

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sonirico/vago/lol"
)

type WSClient struct {
	config WSClientConfig

	conn *websocket.Conn

	sharedSubscriptionByChannel map[string]*sharedSubscription
	dispatcherByChannelType     map[wsChannelType]dispatcher

	done        chan struct{}
	reconnectCh chan struct{}
	mu          sync.RWMutex
	writeMu     sync.Mutex
	closeOnce   sync.Once
	debug       bool

	logger lol.Logger
}

type WSClientOption func(*WSClient)

func WithWSClientLogger(logger lol.Logger) WSClientOption {
	return func(c *WSClient) {
		l := logger
		if l == nil {
			l = lol.NewZerolog()
		}
		c.logger = l
	}
}

func WithWSClientDebug(debug bool) WSClientOption {
	return func(c *WSClient) {
		c.debug = debug
	}
}

func WithWsClientAuthToken(authToken string) WSClientOption {
	return func(c *WSClient) {
		c.config.AuthToken = authToken
	}
}

func NewWSClient(config WSClientConfig, opts ...WSClientOption) *WSClient {
	client := &WSClient{
		config: config,

		sharedSubscriptionByChannel: make(map[string]*sharedSubscription),
		dispatcherByChannelType: map[wsChannelType]dispatcher{
			wsChannelTypeMarketStats: newMarketStatsDispatcher(),
		},

		done:        make(chan struct{}),
		reconnectCh: make(chan struct{}, 1),
		logger:      lol.NewZerolog(),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *WSClient) Connect(ctx context.Context) error {
	return c.connect(ctx)
}

func (c *WSClient) Close() error {
	var err error
	c.closeOnce.Do(func() {
		err = c.close()
	})
	return err
}

func (c *WSClient) close() error {
	close(c.done)

	err := c.dropActiveConnection()

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, s := range c.sharedSubscriptionByChannel {
		s.clear()
	}

	return err
}
