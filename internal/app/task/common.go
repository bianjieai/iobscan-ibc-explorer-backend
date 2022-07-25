package task

import (
	"strings"
	"time"
)

type segment struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

func getHistorySegment() ([]*segment, error) {
	first, err := ibcTxRepo.FirstHistory()
	if err != nil {
		return nil, err
	}

	latest, err := ibcTxRepo.LatestHistory()
	if err != nil {
		return nil, err
	}

	start := time.Unix(first.CreateAt, 0)
	startUnix := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	end := time.Unix(latest.CreateAt, 0)
	endUnix := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 59, time.Local).Unix()

	var step int64 = 12 * 3600
	var segments []*segment
	for temp := startUnix; temp < endUnix; temp += step {
		segments = append(segments, &segment{
			StartTime: temp,
			EndTime:   temp + step - 1,
		})
	}
	return segments, nil
}

func getSegment() ([]*segment, error) {
	first, err := ibcTxRepo.First()
	if err != nil {
		return nil, err
	}

	start := time.Unix(first.CreateAt, 0)
	startUnix := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	end := time.Now()
	endUnix := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 59, time.Local).Unix()

	var step int64 = 24 * 3600
	var segments []*segment
	for temp := startUnix; temp < endUnix; temp += step {
		segments = append(segments, &segment{
			StartTime: temp,
			EndTime:   temp + step - 1,
		})
	}

	return segments, nil
}

// todayUnix 获取今日第一秒和最后一秒的时间戳
func todayUnix() (int64, int64) {
	now := time.Now()
	startUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}

// yesterdayUnix 获取昨日第一秒和最后一秒的时间戳
func yesterdayUnix() (int64, int64) {
	date := time.Now().AddDate(0, 0, -1)
	startUnix := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}

func isConnectionErr(err error) bool {
	return strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "i/o timeout")
}
