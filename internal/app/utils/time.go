package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimeAnalyzer interface {
	CalStatisticTimestamp() int64
	GetSegmentTimeList() []string
}

func NewTimeAnalyzer(digit int, timeUnit string) (TimeAnalyzer, error) {
	switch strings.ToUpper(timeUnit) {
	case "Y":
		return naturalYearTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	case "UY":
		return unNaturalYearTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	case "M":
		return naturalMonthTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	case "UM":
		return unNaturalMonthTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	case "W":
		return weekTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	case "D":
		return naturaldayTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	case "UD":
		return unNaturalDayTimeAnalyzer{
			now:   time.Now(),
			digit: digit,
		}, nil
	default:
		return nil, errors.New("unknown time unit")
	}
}

type (
	// naturalYearTimeAnalyzer with unit y/Y
	naturalYearTimeAnalyzer struct {
		now   time.Time
		digit int
	}

	// naturalYearTimeAnalyzer with unit uy/UY
	unNaturalYearTimeAnalyzer struct {
		now   time.Time
		digit int
	}

	// naturalMonthTimeAnalyzer with unit m/M
	naturalMonthTimeAnalyzer struct {
		now   time.Time
		digit int
	}

	// unNaturalMonthTimeAnalyzer with unit um/UM
	unNaturalMonthTimeAnalyzer struct {
		now   time.Time
		digit int
	}

	// weekTimeAnalyzer with unit w/W
	weekTimeAnalyzer struct {
		now   time.Time
		digit int
	}

	// unNaturalDayTimeAnalyzer with unit ud/UD
	unNaturalDayTimeAnalyzer struct {
		now   time.Time
		digit int
	}

	// naturaldayTimeAnalyzer with unit d/D
	naturaldayTimeAnalyzer struct {
		now   time.Time
		digit int
	}
)

// CalStatisticTimestamp
// see testNaturalYearTimeAnalyzer
func (y naturalYearTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch y.digit {
	case 0:
		year, _, _ := y.now.Date()
		return time.Date(year, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	default:
		date, _, _ := y.now.AddDate(-(y.digit - 1), 0, 0).Date()
		return time.Date(date, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	}
}

func (y naturalYearTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerY{startTimestamp: y.CalStatisticTimestamp()}.GetFormatList()
}

// CalStatisticTimestamp
// see testUnNaturalYearTimeAnalyzer
func (u unNaturalYearTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch u.digit {
	case 0:
		year, _, _ := u.now.Date()
		return time.Date(year, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	default:
		date, month, day := u.now.AddDate(-u.digit, 0, 0).Date()
		return time.Date(date, month, day, 0, 0, 0, 0, time.Local).Unix()
	}
}

func (u unNaturalYearTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerY{startTimestamp: u.CalStatisticTimestamp()}.GetFormatList()
}

// CalStatisticTimestamp
// see testNaturalMonthTimeAnalyzer
func (m naturalMonthTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch m.digit {
	case 0:
		year, month, _ := m.now.Date()
		return time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	default:
		date, month, _ := m.now.AddDate(0, -(m.digit - 1), 0).Date()
		return time.Date(date, month, 1, 0, 0, 0, 0, time.Local).Unix()
	}
}

func (m naturalMonthTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerM{startTimestamp: m.CalStatisticTimestamp()}.GetFormatList()
}

// CalStatisticTimestamp
// see testUnNaturalMonthTimeAnalyzer
func (u unNaturalMonthTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch u.digit {
	case 0:
		year, month, _ := u.now.Date()
		return time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix()
	default:
		date, month, day := u.now.AddDate(0, -u.digit, 0).Date()
		return time.Date(date, month, day, 0, 0, 0, 0, time.Local).Unix()
	}
}

func (u unNaturalMonthTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerM{startTimestamp: u.CalStatisticTimestamp()}.GetFormatList()
}

// CalStatisticTimestamp
// see testWeekTimeAnalyzer
func (w weekTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch w.digit {
	default:
		timestamp := w.now.Unix() - 7*24*3600*int64(w.digit)
		date, month, day := time.Unix(timestamp, 0).Date()
		return time.Date(date, month, day, 0, 0, 0, 0, time.Local).Unix()
	case 0:
		now := w.now
		offset := int(time.Monday - now.Weekday())
		if offset > 0 {
			offset = -6
		}
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset).Unix()
	}
}

// GetSegmentTimeList waiting to be realized, use the `segmentTimeFormatAnalyzeD` temporarily
func (w weekTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerD{startTimestamp: w.CalStatisticTimestamp()}.GetFormatList()
}

func (u unNaturalDayTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch u.digit {
	case 0:
		year, month, day := u.now.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
	default:
		date, month, day := u.now.AddDate(0, 0, -u.digit).Date()
		return time.Date(date, month, day, 0, 0, 0, 0, time.Local).Unix()
	}
}

func (u unNaturalDayTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerD{startTimestamp: u.CalStatisticTimestamp()}.GetFormatList()
}

// CalStatisticTimestamp
// see testDayTimeAnalyzer
func (d naturaldayTimeAnalyzer) CalStatisticTimestamp() int64 {
	switch d.digit {
	case 0:
		year, month, day := d.now.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Unix()
	default:
		date, month, day := d.now.AddDate(0, 0, -(d.digit - 1)).Date()
		return time.Date(date, month, day, 0, 0, 0, 0, time.Local).Unix()
	}
}

func (d naturaldayTimeAnalyzer) GetSegmentTimeList() []string {
	return segmentTimeFormatAnalyzerD{startTimestamp: d.CalStatisticTimestamp()}.GetFormatList()
}

func TimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetPastMonthString return ["yyyy-mm", "yyyy-mm"] format
func GetPastMonthString(n int) []string {
	var res []string
	for i := n; i > 0; i-- {
		date, month, _ := time.Now().AddDate(0, -(i - 1), 0).Date()
		if month < 10 {
			res = append(res, fmt.Sprintf("%d-0%d", date, month))
		} else {
			res = append(res, fmt.Sprintf("%d-%d", date, month))
		}
	}
	return res
}

func GetMonthFromMonthString(s string) int {
	ms := strings.Split(s, "-")[1]
	m, _ := strconv.Atoi(ms)
	return m
}

// GetThreeHaveDataMonth
// build 3 months data, if only have 1 or 2 months have data then align to 3
func GetThreeHaveDataMonth(months []string) []string {
	res := make([]string, 0, 3)
	switch len(months) {
	case 0:
		return nil
	case 1:
		monthDate := StringToMonthDate(months[0])
		lastMothDate := monthDate.AddDate(0, -1, 0)
		lastLastMothDate := lastMothDate.AddDate(0, -1, 0)
		lastMonth := MonthDateToString(lastMothDate)
		lastLastMonth := MonthDateToString(lastLastMothDate)

		res = append(res, lastLastMonth, lastMonth, months[0])
		return res
	case 2:
		monthDate := StringToMonthDate(months[0])
		lastMothDate := monthDate.AddDate(0, -1, 0)
		lastMonth := MonthDateToString(lastMothDate)

		res = append(res, lastMonth)
		res = append(res, months...)
		return res
	default:
		return months[len(months)-3:]
	}
}

// ==========================================================================

type SegmentTimeFormatAnalyzer interface {
	GetFormatList() []string
}

func NewCycleDateFormatAnalyzer(timeUnit string, timestamp int64) (SegmentTimeFormatAnalyzer, error) {
	switch strings.ToUpper(timeUnit) {
	case "Y":
		return segmentTimeFormatAnalyzerY{
			startTimestamp: timestamp,
		}, nil
	case "M":
		return segmentTimeFormatAnalyzerM{
			startTimestamp: timestamp,
		}, nil
	case "D":
		return segmentTimeFormatAnalyzerD{
			startTimestamp: timestamp,
		}, nil
	default:
		return nil, errors.New("unknown timeUnit")
	}
}

type (
	segmentTimeFormatAnalyzerY struct {
		startTimestamp int64
	}

	segmentTimeFormatAnalyzerM struct {
		startTimestamp int64
	}

	segmentTimeFormatAnalyzerD struct {
		startTimestamp int64
	}
)

func (t segmentTimeFormatAnalyzerY) GetFormatList() []string {
	now := time.Now().Unix()
	if t.startTimestamp >= now {
		year, _, _ := time.Unix(t.startTimestamp, 0).Date()
		return []string{strconv.Itoa(year)}
	}

	timeUnix := time.Unix(t.startTimestamp, 0)
	var res []string
	for {
		year, _, _ := timeUnix.Date()
		res = append(res, strconv.Itoa(year))
		timeUnix = timeUnix.AddDate(1, 0, 0)
		if timeUnix.Unix() >= now {
			break
		}
	}

	year, _, _ := time.Unix(now, 0).Date()
	if res[len(res)-1] != strconv.Itoa(year) {
		res = append(res, strconv.Itoa(year))
	}
	return res
}

func (t segmentTimeFormatAnalyzerM) GetFormatList() []string {
	now := time.Now().Unix()
	if t.startTimestamp >= now {
		year, month, _ := time.Unix(t.startTimestamp, 0).Date()
		return []string{t.format(year, month)}
	}

	year, month, _ := time.Unix(t.startTimestamp, 0).Date()
	timeUnix := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	var res []string
	for {
		year, month, _ := timeUnix.Date()
		res = append(res, t.format(year, month))
		timeUnix = timeUnix.AddDate(0, 1, 0)
		if timeUnix.Unix() >= now {
			break
		}
	}
	year, month, _ = time.Unix(now, 0).Date()
	if res[len(res)-1] != t.format(year, month) {
		res = append(res, t.format(year, month))
	}
	return res
}

func (t segmentTimeFormatAnalyzerM) format(year int, month time.Month) string {
	if month < 10 {
		return fmt.Sprintf("%d-0%d", year, month)
	} else {
		return fmt.Sprintf("%d-%d", year, month)
	}
}

func (t segmentTimeFormatAnalyzerD) GetFormatList() []string {
	now := time.Now().Unix()
	if t.startTimestamp >= now {
		date, month, day := time.Unix(t.startTimestamp, 0).Date()
		return []string{t.format(date, month, day)}
	}

	timeUnix := time.Unix(t.startTimestamp, 0)
	var res []string
	for {
		date, month, day := timeUnix.Date()
		res = append(res, t.format(date, month, day))
		timeUnix = timeUnix.AddDate(0, 0, 1)
		if timeUnix.Unix() >= now {
			break
		}
	}
	date, month, day := time.Unix(now, 0).Date()
	if res[len(res)-1] != t.format(date, month, day) {
		res = append(res, t.format(date, month, day))
	}
	return res
}

func (t segmentTimeFormatAnalyzerD) format(year int, month time.Month, day int) string {
	var monthStr string
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	} else {
		monthStr = strconv.Itoa(int(month))
	}

	var dayStr string
	if day < 10 {
		dayStr = fmt.Sprintf("0%d", day)
	} else {
		dayStr = strconv.Itoa(day)
	}

	return fmt.Sprintf("%d-%s-%s", year, monthStr, dayStr)
}
