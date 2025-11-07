package lighter

import "context"

type Client interface {
	GetFundingRates(ctx context.Context) ([]FundingRates, error)

	GetFundings(ctx context.Context, params *GetFundingsParams) ([]Funding, error)

	GetOrderBooks(ctx context.Context, params *GetOrderBooksParams) ([]OrderBook, error)

	GetAccounts(ctx context.Context, params *GetAccountsParams) ([]Account, error)
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

type GetAccountsParams struct {
	By    GetAccountBy `json:"by"`
	Value string       `json:"value"`
}

type Account struct {
	Code                    int32      `json:"code"`
	Message                 string     `json:"message"`
	AccountType             int8       `json:"account_type"`
	Index                   int64      `json:"index"`
	L1Address               string     `json:"l1_address"`
	CancelAllTime           int64      `json:"cancel_all_time"`
	TotalOrderCount         int64      `json:"total_order_count"`
	TotalIsolatedOrderCount int64      `json:"total_isolated_order_count"`
	PendingOrderCount       int64      `json:"pending_order_count"`
	AvailableBalance        string     `json:"available_balance"`
	Status                  uint8      `json:"status"`
	Collateral              string     `json:"collateral"`
	AccountIndex            int64      `json:"account_index"`
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	Positions               []Position `json:"positions"`
	TotalAssetValue         string     `json:"total_asset_value"`
	CrossAssetValue         string     `json:"cross_asset_value"`
}

type Position struct {
	MarketID               uint8  `json:"market_id"`
	Symbol                 string `json:"symbol"`
	InitialMarginFraction  string `json:"initial_margin_fraction"`
	OpenOrderCount         int64  `json:"open_order_count"`
	PendingOrderCount      int64  `json:"pending_order_count"`
	PositionTiedOrderCount int64  `json:"position_tied_order_count"`
	Sign                   int32  `json:"sign"`
	Position               string `json:"position"`
	AVGEntryPrice          string `json:"avg_entry_price"`
	PositionValue          string `json:"position_value"`
	UnrealizedPNL          string `json:"unrealized_pnl"`
	RealizedPNL            string `json:"realized_pnl"`
	LiquidationPrice       string `json:"liquidation_price"`
	TotalFundingPaidOut    string `json:"total_funding_paid_out,omitempty"`
	MarginMode             int32  `json:"margin_mode"`
	AllocatedMargin        string `json:"allocated_margin"`
}
