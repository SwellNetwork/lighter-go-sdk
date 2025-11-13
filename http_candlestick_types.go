package lighter

type GetFundingsParams struct {
	MarketID       uint8      `json:"market_id"`
	Resolution     Resolution `json:"resolution"`
	StartTimestamp int64      `json:"start_timestamp"`
	EndTimestamp   int64      `json:"end_timestamp,omitempty"`
	CountBack      int64      `json:"count_back"`
}
type GetFundingsResult struct {
	Code       int32     `json:"code"`
	Resolution string    `json:"resolution"`
	Fundings   []Funding `json:"fundings"`
}
type Funding struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
	Rate      string `json:"rate"`
	Direction string `json:"direction"`
}

type Resolution string

const (
	Resolution1h Resolution = "1h"
	Resolution1d Resolution = "1d"
)
