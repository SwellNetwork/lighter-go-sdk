package lighter

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

type sharedSubscriptionFinder func(string) (*sharedSubscription, bool)
type callback func(any)

type wsErrorPayload struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type WSError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Channel string `json:"-"`
}

func (e WSError) Error() string {
	if e.Channel != "" {
		return fmt.Sprintf("websocket error (%s): %d %s", e.Channel, e.Code, e.Message)
	}

	return fmt.Sprintf("websocket error: %d %s", e.Code, e.Message)
}

func parseWSError(raw json.RawMessage, channel string) (*WSError, bool) {
	var payload wsErrorPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, false
	}

	if payload.Error == nil {
		return nil, false
	}

	return &WSError{
		Code:    payload.Error.Code,
		Message: payload.Error.Message,
		Channel: canonicalChannelName(channel),
	}, true
}

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
		if err, ok := msg.(error); ok {
			callback(zero, err)
			return
		}

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
	id      string
	payload any
	close   func()
}

func (s *Subscription) GetID() string {
	return s.id
}

func (s *Subscription) GetPayload() any {
	return s.payload
}

func (s *Subscription) Close() {
	s.close()
}

func getChannelType(channel string) wsChannelType {
	splits := strings.Split(canonicalChannelName(channel), ":")

	return wsChannelType(splits[0])
}

func canonicalChannelName(channel string) string {
	return strings.ReplaceAll(channel, "/", ":")
}
