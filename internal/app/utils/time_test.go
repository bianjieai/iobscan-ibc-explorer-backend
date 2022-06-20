package utils

import (
	"testing"
	"time"
)

var testTime = time.Date(2021, time.August, 23, 0, 0, 0, 0, time.Local)

type testTimeAnalyzer struct {
	analyzer TimeAnalyzer
	want     int64
}

func TestTimeAnalyzer(t *testing.T) {
	testYearTimeAnalyzer(t)
	testUnNaturalYearTimeAnalyzer(t)

	testNaturalMonthTimeAnalyzer(t)
	testUnNaturalMonthTimeAnalyzer(t)

	testWeekTimeAnalyzer(t)
	testDayTimeAnalyzer(t)
}

func testYearTimeAnalyzer(t *testing.T) {
	tests := []testTimeAnalyzer{
		{
			analyzer: naturalYearTimeAnalyzer{
				now:   testTime,
				digit: 0,
			},
			want: 1609430400, // 2021-01-01
		},
		{
			analyzer: naturalYearTimeAnalyzer{
				now:   testTime,
				digit: 1,
			},
			want: 1609430400, // 2021-01-01
		},
		{
			analyzer: naturalYearTimeAnalyzer{
				now:   testTime,
				digit: 2,
			},
			want: 1577808000, // 2020-01-01
		},
		{
			analyzer: naturalYearTimeAnalyzer{
				now:   testTime,
				digit: 10,
			},
			want: 1325347200, // 2012-01-01
		},
	}

	for _, v := range tests {
		res := v.analyzer.CalStatisticTimestamp()
		if v.want != res {
			t.Fatalf("expected: %d, actual: %d", v.want, res)
		}
	}
}

func testUnNaturalYearTimeAnalyzer(t *testing.T) {
	tests := []testTimeAnalyzer{
		{
			analyzer: unNaturalYearTimeAnalyzer{
				now:   testTime,
				digit: 0,
			},
			want: 1609430400, // 2021-01-01
		},
		{
			analyzer: unNaturalYearTimeAnalyzer{
				now:   testTime,
				digit: 1,
			},
			want: 1598112000, // 2020-08-23
		},
		{
			analyzer: unNaturalYearTimeAnalyzer{
				now:   testTime,
				digit: 2,
			},
			want: 1566489600, // 2019-08-23
		},
		{
			analyzer: unNaturalYearTimeAnalyzer{
				now:   testTime,
				digit: 10,
			},
			want: 1314028800, // 2011-08-23
		},
	}

	for _, v := range tests {
		res := v.analyzer.CalStatisticTimestamp()
		if v.want != res {
			t.Fatalf("expected: %d, actual: %d", v.want, res)
		}
	}
}

func testNaturalMonthTimeAnalyzer(t *testing.T) {
	tests := []testTimeAnalyzer{
		{
			analyzer: naturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 0,
			},
			want: 1627747200, // 2021-08-01
		},
		{
			analyzer: naturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 1,
			},
			want: 1627747200, // 2021-08-01
		},
		{
			analyzer: naturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 2,
			},
			want: 1625068800, // 2021-07-01
		},
		{
			analyzer: naturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 10,
			},
			want: 1604160000, // 2020-11-01
		},
	}

	for _, v := range tests {
		res := v.analyzer.CalStatisticTimestamp()
		if v.want != res {
			t.Fatalf("expected: %d, actual: %d", v.want, res)
		}
	}
}

func testUnNaturalMonthTimeAnalyzer(t *testing.T) {
	tests := []testTimeAnalyzer{
		{
			analyzer: unNaturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 0,
			},
			want: 1627747200, // 2021-08-01
		},
		{
			analyzer: unNaturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 1,
			},
			want: 1626969600, // 2021-07-023
		},
		{
			analyzer: unNaturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 2,
			},
			want: 1624377600, // 2021-06-23
		},
		{
			analyzer: unNaturalMonthTimeAnalyzer{
				now:   testTime,
				digit: 10,
			},
			want: 1603382400, // 2020-10-23
		},
	}

	for _, v := range tests {
		res := v.analyzer.CalStatisticTimestamp()
		if v.want != res {
			t.Fatalf("expected: %d, actual: %d", v.want, res)
		}
	}
}

func testWeekTimeAnalyzer(t *testing.T) {
	tests := []testTimeAnalyzer{
		{
			analyzer: weekTimeAnalyzer{
				now:   testTime,
				digit: 0,
			},
			want: 1629648000, // 2021-08-23, monday
		},
		{
			analyzer: weekTimeAnalyzer{
				now:   testTime,
				digit: 1,
			},
			want: 1629043200, // 2021-08-16
		},
		{
			analyzer: weekTimeAnalyzer{
				now:   testTime,
				digit: 2,
			},
			want: 1628438400, // 2021-08-09
		},
		{
			analyzer: weekTimeAnalyzer{
				now:   testTime,
				digit: 10,
			},
			want: 1623600000, // 2021-06-14
		},
	}

	for _, v := range tests {
		res := v.analyzer.CalStatisticTimestamp()
		if v.want != res {
			t.Fatalf("expected: %d, actual: %d", v.want, res)
		}
	}
}

func testDayTimeAnalyzer(t *testing.T) {
	tests := []testTimeAnalyzer{
		{
			analyzer: naturaldayTimeAnalyzer{
				now:   testTime,
				digit: 0,
			},
			want: 1629648000, // 2021-08-23, monday
		},
		{
			analyzer: naturaldayTimeAnalyzer{
				now:   testTime,
				digit: 1,
			},
			want: 1629561600, // 2021-08-22
		},
		{
			analyzer: naturaldayTimeAnalyzer{
				now:   testTime,
				digit: 2,
			},
			want: 1629475200, // 2021-08-21
		},
		{
			analyzer: naturaldayTimeAnalyzer{
				now:   testTime,
				digit: 10,
			},
			want: 1628784000, // 2021-08-13
		},
	}

	for _, v := range tests {
		res := v.analyzer.CalStatisticTimestamp()
		if v.want != res {
			t.Fatalf("expected: %d, actual: %d", v.want, res)
		}
	}
}
