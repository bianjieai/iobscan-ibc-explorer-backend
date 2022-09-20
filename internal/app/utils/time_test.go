package utils

import "testing"

func TestValidateDenom(t *testing.T) {
	err := ValidateDenom("ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2")
	t.Log(err)
}
