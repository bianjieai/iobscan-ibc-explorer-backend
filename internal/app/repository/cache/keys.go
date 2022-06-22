package cache

import "time"

// redis key expiration
const (
	oneHour                    = 1 * time.Hour
	threeHours                 = 3 * time.Hour
	oneDay                     = 24 * time.Hour
	oneMin                     = 60 * time.Second
	NoExpiration time.Duration = -1
)

// redis key
const (
	tokenPrice  = "token_price"
	denomSupply = "denom_supply_%s"
)
