package repository

import "testing"

func TestTxRepo_GetActiveAccountsOfDay(t *testing.T) {
	val, err := new(TxRepo).GetActiveAccountsOfDay("bigbang", 1619107200, 1619193600)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val)
}
