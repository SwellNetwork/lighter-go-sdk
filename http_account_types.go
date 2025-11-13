package lighter

type GetAccountsParams struct {
	By    GetAccountBy `json:"by"`
	Value string       `json:"value"`
}
type GetAccountsResult struct {
	Code     int32     `json:"code"`
	Message  string    `json:"message"`
	Total    int64     `json:"total"`
	Accounts []Account `json:"accounts"`
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

type GetAccountBy string

const (
	GetAccountByIndex     GetAccountBy = "index"
	GetAccountByL1Address GetAccountBy = "l1_address"
)
