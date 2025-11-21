package lighter

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

const (
	minReconnectBackoff             = time.Second
	maxReconnectBackoff             = 30 * time.Second
	defaultSubscribeRetryMinBackoff = time.Second
	defaultSubscribeRetryMaxBackoff = 30 * time.Second
)

func (c *WSClient) connect(ctx context.Context) error {
	select {
	case <-c.done:
		return fmt.Errorf("ws client closed")
	default:
	}

	c.mu.RLock()
	if c.conn != nil {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	conn, _, err := new(websocket.Dialer).DialContext(ctx, c.config.BaseURL, nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if c.conn != nil {
		c.mu.Unlock()
		_ = conn.Close()
		return nil
	}
	c.conn = conn
	c.mu.Unlock()

	if err = c.resubscribeAll(); err != nil {
		c.logger.Errorf("failed to replay subscriptions: %v", err)

		if dropActiveConnErr := c.dropActiveConnection(); dropActiveConnErr != nil {
			c.logger.Errorf("failed to drop active connections: %v", dropActiveConnErr)
		}
		return err
	}

	go c.messageLoop(ctx, conn)
	go c.heartbeatLoop(ctx)

	return nil
}

func (c *WSClient) messageLoop(ctx context.Context, conn *websocket.Conn) {
	defer c.clearConnection(conn)

	readTimeout := c.config.ReadTimeout
	if readTimeout <= 0 {
		if c.config.PingInterval > 0 {
			readTimeout = c.config.PingInterval * 2
		} else {
			readTimeout = 0
		}
	}

	for {
		select {
		case <-c.done:
			return
		case <-ctx.Done():
			return
		default:
			if readTimeout > 0 {
				if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
					c.logger.Errorf("failed to set websocket read deadline: %v", err)
					c.triggerReconnect()
					return
				}
			}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.logger.Errorf("websocket read error: %v", err)
				}
				c.triggerReconnect()
				return
			}

			if c.debug {
				c.logger.Debugf("[<] %s", string(msg))
			}

			var wsMsg wsMessage
			if err = json.Unmarshal(msg, &wsMsg); err != nil {
				c.logger.Errorf("websocket message parse error: %v", err)
				continue
			}

			if err = c.dispatch(wsMsg); err != nil {
				c.logger.Errorf("failed to dispatch websocket message: %v", err)
			}
		}
	}
}

func (c *WSClient) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ctx.Done():
			return
		case <-c.reconnectCh:
			c.reconnect(ctx)
			return
		case <-ticker.C:
			if err := c.writeJSON(wsCommandPing); err != nil {
				c.reconnect(ctx)
				return
			}
		}
	}
}

func (c *WSClient) reconnect(ctx context.Context) {
	if err := c.dropActiveConnection(); err != nil {
		c.logger.Errorf("failed to close websocket before reconnect: %v", err)
	}

	attempt := 0
	maxAttempts := c.config.ReconnectAttempts

	for {
		select {
		case <-c.done:
			return
		case <-ctx.Done():
			return
		default:
		}

		if err := c.connect(ctx); err == nil {
			return
		} else {
			attempt++
			c.logger.Errorf("reconnect attempt %d failed: %v", attempt, err)
		}

		if maxAttempts > 0 && attempt >= maxAttempts {
			c.logger.Errorf("exhausted reconnect attempts (%d)", maxAttempts)
			return
		}

		wait := c.reconnectBackoffDuration(attempt)
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-c.done:
			timer.Stop()
			return
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (c *WSClient) dispatch(msg wsMessage) error {
	if wsErr, ok := parseWSError(msg.Raw, msg.Channel); ok {
		c.logger.Errorf("websocket error received: %v", wsErr)
		c.dispatchError(msg.Channel, wsErr)
		return nil
	}

	d, dispatcherFound := c.dispatcherByChannelType[getChannelType(msg.Channel)]
	if !dispatcherFound {
		return fmt.Errorf("no dispatcher for channel: %s", msg.Channel)
	}

	finder := func(channel string) (*sharedSubscription, bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		s, sharedSubscriptionFound := c.sharedSubscriptionByChannel[canonicalChannelName(channel)]
		return s, sharedSubscriptionFound
	}

	return d(finder, msg)
}

func (c *WSClient) dispatchError(channel string, err error) {
	if channel == "" {
		c.mu.RLock()
		subs := make([]*sharedSubscription, 0, len(c.sharedSubscriptionByChannel))
		for _, s := range c.sharedSubscriptionByChannel {
			subs = append(subs, s)
		}
		c.mu.RUnlock()

		for _, s := range subs {
			s.dispatchError(err)
		}
		return
	}

	lookup := canonicalChannelName(channel)

	c.mu.RLock()
	s, ok := c.sharedSubscriptionByChannel[lookup]
	c.mu.RUnlock()

	if ok && s != nil {
		s.dispatchError(err)
	}
}

func (c *WSClient) writeJSON(v any) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	conn := c.conn
	if conn == nil {
		return fmt.Errorf("connection closed")
	}

	if c.debug {
		bts, _ := json.Marshal(v)
		c.logger.Debugf("[>] %s\n", string(bts))
	}

	return conn.WriteJSON(v)
}

func (c *WSClient) triggerReconnect() {
	select {
	case c.reconnectCh <- struct{}{}:
	default:
	}
}

func (c *WSClient) clearConnection(conn *websocket.Conn) {
	if conn != nil {
		_ = conn.Close()
	}

	c.mu.Lock()
	if c.conn == conn {
		c.conn = nil
	}
	c.mu.Unlock()
}

func (c *WSClient) dropActiveConnection() error {
	c.mu.Lock()
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	if conn != nil {
		return conn.Close()
	}

	return nil
}

func (c *WSClient) reconnectBackoffDuration(attempt int) time.Duration {
	if attempt <= 0 {
		return minReconnectBackoff
	}

	return min(minReconnectBackoff<<(attempt-1), maxReconnectBackoff)
}

func (c *WSClient) subscribeRetryBackoffDuration(attempt int) time.Duration {
	minBackoff := c.config.SubscribeRetryMinBackoff
	if minBackoff <= 0 {
		minBackoff = defaultSubscribeRetryMinBackoff
	}

	maxBackoff := c.config.SubscribeRetryMaxBackoff
	if maxBackoff <= 0 {
		maxBackoff = defaultSubscribeRetryMaxBackoff
	}

	if attempt <= 0 {
		return minBackoff
	}

	backoff := minBackoff << (attempt - 1)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	return backoff
}
