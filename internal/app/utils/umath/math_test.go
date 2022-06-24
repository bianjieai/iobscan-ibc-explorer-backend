package umath

import (
	"fmt"
	"testing"
)

func TestPercent(t *testing.T) {
	fmt.Println(PercentFromInt(1, 2, 2))

	fmt.Println(PercentFromInt(211, 3, 2))

	fmt.Println(PercentFromInt(211, 232, 3))

	fmt.Println(PercentFromInt(0, 232, 2))

	fmt.Println(PercentFromInt(211, 0, 2))

	fmt.Println(CalculateRate(321, 122, 2))

	fmt.Println(CalculateRate(321, 0, 2))

}
