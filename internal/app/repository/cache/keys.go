package cache

import "time"

// redis key
const (
	oneHour                    = 1 * time.Hour
	threeHours                 = 3 * time.Hour
	oneDay                     = 24 * time.Hour
	oneMin                     = 60 * time.Second
	NoExpiration time.Duration = -1
)
