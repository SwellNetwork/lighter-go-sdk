package lighter

type GetOrderBooksParams struct {
	MarketID uint8 `json:"market_id"`
}

type GetOrderBooksResult struct {
	Code       int32       `json:"code"`
	Message    string      `json:"message"`
	OrderBooks []OrderBook `json:"order_books"`
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

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)
