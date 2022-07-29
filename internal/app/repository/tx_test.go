package repository

import "testing"

func TestTxRepo_GetTimePeriodByUpdateClient(t *testing.T) {
	val, val1, clientId, err := new(TxRepo).GetTimePeriodByUpdateClient("irishub_qa", "iaa1u3tpcx5088rx3lzzt0gg73lt9zugrjp730apj8", 1656557855)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val, val1, clientId)
}
