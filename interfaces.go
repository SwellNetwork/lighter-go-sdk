package lighter

import "context"

type Client interface {
	GetFundingRates(ctx context.Context) ([]FundingRates, error)

	GetFundings(ctx context.Context, params *GetFundingsParams) ([]Funding, error)

	GetOrderBooks(ctx context.Context, params *GetOrderBooksParams) ([]OrderBook, error)
}

type FundingRates struct {
	MarketID uint8   `json:"market_id"`
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Rate     float64 `json:"rate"`
}

type GetFundingsParams struct {
	MarketID       uint8      `json:"market_id"`
	Resolution     Resolution `json:"resolution"`
	StartTimestamp int64      `json:"start_timestamp"`
	EndTimestamp   int64      `json:"end_timestamp,omitempty"`
	CountBack      int64      `json:"count_back"`
}
type Funding struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
	Rate      string `json:"rate"`
	Direction string `json:"direction"`
}

type Resolution string

const (
	Resolution1h Resolution = "1h"
	Resolution1d Resolution = "1d"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type GetOrderBooksParams struct {
	MarketID uint8 `json:"market_id"`
}

type OrderBook struct {
	Symbol                 string `json:"symbol"`
	MarketId               uint8  `json:"market_id"`
	Status                 Status `json:"status"`
	TakerFee               string `json:"taker_fee"`
	MakerFee               string `json:"maker_fee"`
	LiquidationFee         string `json:"liquidation_fee"`
	MinBaseAmount          string `json:"min_base_amount"`
	MinQuoteAmount         string `json:"min_quote_amount"`
	OrderQuoteLimit        string `json:"order_quote_limit"`
	SupportedSizeDecimals  uint8  `json:"supported_size_decimals"`
	SupportedPriceDecimals uint8  `json:"supported_price_decimals"`
	SupportedQuoteDecimals uint8  `json:"supported_quote_decimals"`
}
