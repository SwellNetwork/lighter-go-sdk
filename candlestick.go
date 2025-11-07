package lighter

import (
	"context"
	"fmt"
	"strconv"
)

const (
	PathGetFundings = "/api/v1/fundings"
)

type GetFundingsResult struct {
	Code       int32     `json:"code"`
	Resolution string    `json:"resolution"`
	Fundings   []Funding `json:"fundings"`
}

func (c *HTTPClient) GetFundings(ctx context.Context, params *GetFundingsParams) ([]Funding, error) {
	queryParams := map[string]string{
		"market_id":       strconv.FormatInt(int64(params.MarketID), 10),
		"resolution":      string(params.Resolution),
		"start_timestamp": strconv.FormatInt(params.StartTimestamp, 10),
		"end_timestamp":   strconv.FormatInt(params.EndTimestamp, 10),
		"count_back":      strconv.FormatInt(int64(params.CountBack), 10),
	}

	var result GetFundingsResult

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&result).
		Get(PathGetFundings)
	if err != nil {
		return nil, fmt.Errorf("get fundings: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("get fundings: %v", resp.Error())
	}

	return result.Fundings, nil
}
