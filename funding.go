package lighter

import (
	"context"
	"fmt"
)

const (
	PathGetFundingRates = "/api/v1/funding-rates"
)

type GetFundingRatesResult struct {
	Code         int32          `json:"code"`
	Message      string         `json:"message"`
	FundingRates []FundingRates `json:"funding_rates"`
}

func (c *HTTPClient) GetFundingRates(ctx context.Context) ([]FundingRates, error) {
	var result GetFundingRatesResult

	resp, err := c.client.R().SetContext(ctx).SetResult(&result).Get(PathGetFundingRates)
	if err != nil {
		return nil, fmt.Errorf("get funding rates: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("get funding rates: %s", resp.Status())
	}

	return result.FundingRates, nil
}
