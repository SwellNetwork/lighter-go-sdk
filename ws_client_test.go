//go:build integration
// +build integration

package lighter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestWSClientIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WSClientIntegrationTestSuite))
}

type WSClientIntegrationTestSuite struct {
	suite.Suite

	client *WSClient
}

func (ts *WSClientIntegrationTestSuite) SetupTest() {
	ts.client = NewWSClient(
		DefaultMainnetWSClientConfig(),
		WithWSClientDebug(true),
	)
}

func (ts *WSClientIntegrationTestSuite) TestMarketStats() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ts.T().Log("Connecting to websocket")
	if err := ts.client.Connect(ctx); err != nil {
		ts.T().Fatalf("Failed to connect: %v", err)
	}
	defer ts.client.Close()

	sub, err := ts.client.MarketStats(
		MarketStatsParams{MarketID: "0"},
		func(marketStats MarketStats, err error) {
			if err != nil {
				ts.T().Fatalf("Failed to receive trades: %v", err)
			}
			ts.T().Logf("Received: %v", marketStats)
		},
	)

	if err != nil {
		ts.T().Fatalf("Failed to subscribe to market stats: %v", err)
	}

	ts.T().Log("Subscribed to market stats")

	defer sub.Close()

	<-time.After(time.Second * 5)
	ts.T().Log("Unsubscribing from market stats")
	sub.Close()
}
