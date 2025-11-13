package lighter

import "context"

type Client interface {
	GetFundingRates(ctx context.Context) ([]FundingRates, error)
	GetFundings(ctx context.Context, params *GetFundingsParams) ([]Funding, error)
	GetOrderBooks(ctx context.Context, params *GetOrderBooksParams) ([]OrderBook, error)
	GetAccounts(ctx context.Context, params *GetAccountsParams) ([]Account, error)
}
