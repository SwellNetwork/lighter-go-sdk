package lighter

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

type sharedSubscriptionFinder func(string) (*sharedSubscription, bool)
type callback func(any)

func subscribeTyped[T any](
	c *WSClient,
	channel string,
	callback func(T, error),
) (*Subscription, error) {
	if callback == nil {
		return nil, fmt.Errorf("callback cannot be nil")
	}

	var zero T

	return c.subscribe(channel, func(msg any) {
		typed, ok := msg.(T)
		if !ok {
			callback(zero, fmt.Errorf("invalid message type: %T", msg))
			return
		}

		callback(typed, nil)
	})
}

type wsCommand struct {
	Type    string `json:"type"`
	Channel string `json:"channel,omitempty"`
	Auth    string `json:"auth,omitempty"`
}

var (
	wsCommandPing      = wsCommand{Type: "ping"}
	wsCommandSubscribe = func(channel string, auth string) wsCommand {
		return wsCommand{Type: "subscribe", Channel: channel, Auth: auth}
	}
	wsCommandUnsubscribe = func(channel string, auth string) wsCommand {
		return wsCommand{Type: "unsubscribe", Channel: channel, Auth: auth}
	}
)

type wsChannelType string

const (
	wsChannelTypeMarketStats wsChannelType = "market_stats"
)

type wsMessage struct {
	Channel string          `json:"channel"`
	Type    string          `json:"type"`
	Raw     json.RawMessage `json:"-"`
}

func (m *wsMessage) UnmarshalJSON(data []byte) error {
	type wsMessageAlias struct {
		Channel string `json:"channel"`
		Type    string `json:"type"`
	}

	var alias wsMessageAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	m.Channel = alias.Channel
	m.Type = alias.Type
	m.Raw = append(m.Raw[:0], data...)

	return nil
}

type Subscription struct {
	ID      string
	Payload any
	Close   func()
}

func getChannelType(channel string) wsChannelType {
	splits := strings.Split(canonicalChannelName(channel), ":")

	return wsChannelType(splits[0])
}

func canonicalChannelName(channel string) string {
	return strings.ReplaceAll(channel, "/", ":")
}
