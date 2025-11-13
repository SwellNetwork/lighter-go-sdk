package lighter

type GetFundingRatesResult struct {
	Code         int32          `json:"code"`
	Message      string         `json:"message"`
	FundingRates []FundingRates `json:"funding_rates"`
}

type FundingRates struct {
	MarketID uint8   `json:"market_id"`
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Rate     float64 `json:"rate"`
}
