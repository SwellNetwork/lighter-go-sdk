package lighter

type GetAccountBy string

const (
	GetAccountByIndex     GetAccountBy = "index"
	GetAccountByL1Address GetAccountBy = "l1_address"
)

type Resolution string

const (
	Resolution1h Resolution = "1h"
	Resolution1d Resolution = "1d"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)
