package lighter

import (
	"context"
	"fmt"
)

const (
	PathGetAccounts = "/api/v1/account"
)

func (c *HTTPClient) GetAccounts(ctx context.Context, params *GetAccountsParams) ([]Account, error) {
	queryParams := map[string]string{
		"by":    string(params.By),
		"value": params.Value,
	}
	var result GetAccountsResult

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&result).
		Get(PathGetAccounts)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("get account: %s", resp.Error())
	}

	return result.Accounts, nil
}
