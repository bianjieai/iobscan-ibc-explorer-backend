package utils

import "time"

func RunTimer(num int, uint Unit, fn func()) {
	go func() {
		// run once right now
		fn()
		for {
			now := time.Now()
			next := now.Add(ParseDuration(num, uint))
			next = TruncateTime(next, uint)
			t := time.NewTimer(next.Sub(now))
			select {
			case <-t.C:
				fn()
			}
		}
	}()
}

// TodayUnix 获取今日第一秒和最后一秒的时间戳
func TodayUnix() (int64, int64) {
	now := time.Now()
	startUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}

// YesterdayUnix 获取昨日第一秒和最后一秒的时间戳
func YesterdayUnix() (int64, int64) {
	date := time.Now().AddDate(0, 0, -1)
	startUnix := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}
