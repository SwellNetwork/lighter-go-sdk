package lighter

import (
	"fmt"

	"github.com/google/uuid"
)

func (c *WSClient) subscribe(
	channel string,
	callback func(any),
) (*Subscription, error) {
	if callback == nil {
		return nil, fmt.Errorf("callback cannot be nil")
	}

	lookupKey := canonicalChannelName(channel)

	c.mu.Lock()
	s, exists := c.sharedSubscriptionByChannel[lookupKey]
	if !exists {
		s = newSharedSubscription(
			channel,
			c.config.AuthToken,
			func(channel string, auth string) {
				if err := c.writeJSON(wsCommandSubscribe(channel, auth)); err != nil {
					c.logger.Errorf("failed to subscribe: %v", err)
				}
			},
			func(channel string, auth string) {
				c.mu.Lock()
				defer c.mu.Unlock()
				delete(c.sharedSubscriptionByChannel, lookupKey)

				if err := c.writeJSON(wsCommandUnsubscribe(channel, auth)); err != nil {
					c.logger.Errorf("failed to unsubscribe: %v", err)
				}
			},
		)

		c.sharedSubscriptionByChannel[lookupKey] = s
	}

	c.mu.Unlock()

	subscriberID := uuid.New().String()
	s.addSubscriber(subscriberID, callback)

	return &Subscription{
		id: subscriberID,
		close: func() {
			s.removeSubscriber(subscriberID)
		},
	}, nil
}

func (c *WSClient) resubscribeAll() error {
	for _, s := range c.sharedSubscriptionByChannel {
		if err := c.writeJSON(wsCommandSubscribe(s.channel, s.auth)); err != nil {
			return err
		}
	}

	return nil
}
