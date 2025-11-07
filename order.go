package lighter

import (
	"context"
	"fmt"
	"strconv"
)

const (
	PathOrderBooks = "/api/v1/orderBooks"
)

type GetOrderBooksResult struct {
	Code       int32       `json:"code"`
	Message    string      `json:"message"`
	OrderBooks []OrderBook `json:"order_books"`
}

func (c *HTTPClient) GetOrderBooks(ctx context.Context, params *GetOrderBooksParams) ([]OrderBook, error) {
	var queryParams map[string]string

	if params != nil {
		queryParams["market_id"] = strconv.FormatInt(int64(params.MarketID), 10)
	}

	var result GetOrderBooksResult

	resp, err := c.client.R().SetContext(ctx).SetQueryParams(queryParams).SetResult(&result).Get(PathOrderBooks)
	if err != nil {
		return nil, fmt.Errorf("get order books: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("get order books: %s", resp.Error())
	}

	return result.OrderBooks, nil
}
