package repository

import "testing"

func TestTxRepo_GetTimePeriodByUpdateClient(t *testing.T) {
	val, val1, err := new(TxRepo).GetTimePeriodByUpdateClient("bigbang")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val, val1)
}
