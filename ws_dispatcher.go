package lighter

import (
	"fmt"

	"github.com/goccy/go-json"
)

type dispatcher func(find sharedSubscriptionFinder, msg wsMessage) error

func newMarketStatsDispatcher() dispatcher {
	type marketStatsPayload struct {
		MarketStats MarketStats `json:"market_stats"`
	}

	return func(find sharedSubscriptionFinder, msg wsMessage) error {
		if getChannelType(msg.Channel) != wsChannelTypeMarketStats {
			return nil
		}

		var payload marketStatsPayload
		if err := json.Unmarshal(msg.Raw, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal market stats message: %w", err)
		}

		s, ok := find(msg.Channel)
		if !ok || s == nil {
			return nil
		}

		s.dispatch(payload.MarketStats)
		return nil
	}
}
