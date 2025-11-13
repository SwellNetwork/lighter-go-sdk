//go:build integration
// +build integration

package lighter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestHTTPClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPClientIntegrationTestSuite))
}

type HTTPClientIntegrationTestSuite struct {
	suite.Suite

	client *HTTPClient
}

func (ts *HTTPClientIntegrationTestSuite) SetupSuite() {
	config := DefaultTestnetHTTPClientConfig()

	ts.client = NewHTTPClient(config)
}

func (ts *HTTPClientIntegrationTestSuite) TestGetFundingRates() {
	ctx := context.Background()

	result, err := ts.client.GetFundingRates(ctx)

	ts.T().Log("get funding rates", result)

	ts.NoError(err)
	ts.NotNil(result)
}

func (ts *HTTPClientIntegrationTestSuite) TestGetFundings() {
	ctx := context.Background()

	params := &GetFundingsParams{
		MarketID:       0,
		Resolution:     Resolution1h,
		StartTimestamp: time.Now().Add(-3 * time.Hour).Unix(),
		EndTimestamp:   time.Now().Unix(),
		CountBack:      3,
	}
	result, err := ts.client.GetFundings(ctx, params)

	ts.T().Log("get fundings", result)

	ts.NoError(err)
	ts.NotNil(result)
}

func (ts *HTTPClientIntegrationTestSuite) TestGetOrderBooks() {
	ctx := context.Background()

	result, err := ts.client.GetOrderBooks(ctx, nil)

	ts.T().Log("get orderbooks", result)

	ts.NoError(err)
	ts.NotNil(result)
}

func (ts *HTTPClientIntegrationTestSuite) TestGetAccounts() {
	ctx := context.Background()

	params := &GetAccountsParams{
		By:    GetAccountByL1Address,
		Value: "",
	}

	result, err := ts.client.GetAccounts(ctx, params)

	ts.T().Log("get accounts", result)

	ts.NoError(err)
	ts.NotNil(result)
}
