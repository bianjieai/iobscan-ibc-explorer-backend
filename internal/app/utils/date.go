package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DateFmtYYYYMM         = "2006-01"
	DateFmtYYYYMMDD       = "2006-01-02"
	DateFmtYYYYMMDDHHmmss = "2006-01-02 15:04:05"
	DateISO8601           = "2006-01-02T15:04:05Z"
	DateISO8601WithZone   = "2006-01-02T15:04:05Z07:00"
)

const (
	_ Unit = iota
	Day
	Hour
	Min
	Sec
)

type Unit int

var CstZone = time.FixedZone("CST", 8*3600) // 东八

func TruncateTime(t time.Time, unit Unit) time.Time {
	switch unit {
	case Day:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Hour:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case Min:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	case Sec:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	}
	panic("not exist unit")
}

func ParseDuration(num int, unit Unit) time.Duration {
	switch unit {
	case Day:
		return time.Duration(num*24) * time.Hour
	case Hour:
		return time.Duration(num) * time.Hour
	case Min:
		return time.Duration(num) * time.Minute
	case Sec:
		return time.Duration(num) * time.Second
	}
	panic("not exist unit")
}

// Get the first and last day of the month, input e.g "2020-08"  "2013-01" "1998-3"
func GetMonthStartAndEnd(yearAndMonth string) (string, string, error) {
	matchString, _ := regexp.MatchString("^\\d{4}-(0[1-9]|1[012]|[0-9])$", yearAndMonth)
	if !matchString {
		return "", "", fmt.Errorf("failed to parse time, please input the correct pattern")
	}

	split := strings.Split(yearAndMonth, "-")
	year := split[0]
	month := split[1]

	yInt, _ := strconv.Atoi(year)
	yMonth, _ := strconv.Atoi(month)

	startDate := time.Date(yInt, time.Month(yMonth), 1, 0, 0, 0, 0, time.Local).Format(DateFmtYYYYMMDD)
	endDate := time.Date(yInt, time.Month(yMonth+1), 0, 0, 0, 0, 0, time.Local).Format(DateFmtYYYYMMDD)
	return startDate, endDate, nil
}

func FmtTime(t time.Time, fmt string) string {
	return t.Format(fmt)
}

// return e.g. 2020-09-15
func GetCurrentTime() string {
	return time.Now().Format(DateFmtYYYYMMDD)
}

func StringToDate(str string) time.Time {
	res, _ := time.Parse(DateFmtYYYYMMDD, str)
	return res
}

func StringToMonthDate(str string) time.Time {
	res, _ := time.Parse(DateFmtYYYYMM, str)
	return res
}

func MonthDateToString(t time.Time) string {
	return t.Format(DateFmtYYYYMM)
}

func StringToDateWithCST(str string) time.Time {
	res, _ := time.ParseInLocation(DateFmtYYYYMMDD, str, CstZone)
	return res
}

func ISO8601ToGMT(str string) string {
	t, err := time.Parse(DateISO8601, str)
	if err != nil {
		return ""
	}
	t = t.Add(8 * time.Hour)
	return t.Format(DateFmtYYYYMMDD)
}

func ISO8601ToGMTWithSecond(str string) string {
	t, err := time.Parse(DateISO8601WithZone, str)
	if err != nil {
		return ""
	}
	t = t.Add(8 * time.Hour)
	return t.Format(DateFmtYYYYMMDDHHmmss)
}

func ISO8601StrToTime(str string) (time.Time, error) {
	t, err := time.Parse(DateISO8601WithZone, str)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
