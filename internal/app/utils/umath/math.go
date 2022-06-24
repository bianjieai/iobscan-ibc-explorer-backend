package umath

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func PercentFromInt(dividend, divisor, precision int64) string {
	if divisor == 0 {
		return "-"
	}
	if dividend == 0 {
		return "0%"
	}

	value := decimal.NewFromInt(dividend).
		DivRound(decimal.NewFromInt(divisor), int32(precision+2)).
		Mul(decimal.NewFromInt(100)).
		String()
	return fmt.Sprintf("%s%s", value, "%")
}

func CalculateRate(dividend, divisor, precision int64) float64 {
	if divisor == 0 {
		return -1
	}
	if dividend == 0 {
		return 0
	}

	value, _ := decimal.NewFromInt(dividend).
		DivRound(decimal.NewFromInt(divisor), int32(precision)).
		Float64()
	return value
}
