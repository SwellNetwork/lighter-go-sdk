package lighter

import "fmt"

type MarketStatsParams struct {
	MarketID string `json:"marketId"`
	IsAll    bool   `json:"isAll"`
}

type MarketStats struct {
	MarketId              int     `json:"market_id"`
	IndexPrice            string  `json:"index_price"`
	MarkPrice             string  `json:"mark_price"`
	OpenInterest          string  `json:"open_interest"`
	LastTradePrice        string  `json:"last_trade_price"`
	CurrentFundingRate    string  `json:"current_funding_rate"`
	FundingRate           string  `json:"funding_rate"`
	FundingTimestamp      int64   `json:"funding_timestamp"`
	DailyBaseTokenVolume  float64 `json:"daily_base_token_volume"`
	DailyQuoteTokenVolume float64 `json:"daily_quote_token_volume"`
	DailyPriceLow         float64 `json:"daily_price_low"`
	DailyPriceHigh        float64 `json:"daily_price_high"`
	DailyPriceChange      float64 `json:"daily_price_change"`
}

func (c *WSClient) MarketStats(
	params MarketStatsParams,
	callback func(MarketStats, error),
) (*Subscription, error) {
	channel, err := getMarketStatsChannel(params)
	if err != nil {
		return nil, err
	}

	return subscribeTyped(c, channel, callback)
}

func getMarketStatsChannel(params MarketStatsParams) (string, error) {
	if params.IsAll {
		return "market_stats/all", nil
	}

	if params.MarketID == "" {
		return "", fmt.Errorf("marketId is required when isAll is false")
	}

	return fmt.Sprintf("market_stats/%s", params.MarketID), nil
}
