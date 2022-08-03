package repository

import "testing"

func TestTxRepo_GetTimePeriodByUpdateClient(t *testing.T) {
	val, val1, clientId, err := new(TxRepo).GetTimePeriodByUpdateClient("irishub_qa", "iaa1u3tpcx5088rx3lzzt0gg73lt9zugrjp730apj8", 1656557855)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val, val1, clientId)
}

func TestTxRepo_GetChannelOpenConfirmTime(t *testing.T) {
	val, err := new(TxRepo).GetChannelOpenConfirmTime("bigbang", "channel-182")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}

func TestTxRepo_GetLogByHash(t *testing.T) {
	val, err := new(TxRepo).GetLogByHash("irishub_qa", "1CE38597D5C4CA3E88E2DD59B7A639EF3BDCCCA9377E65C63051B994A2E8B22C")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val)
}

func TestTxRepo_GetActiveAccountsOfDay(t *testing.T) {
	val, err := new(TxRepo).GetActiveAccountsOfDay("bigbang", 1619107200, 1619193600)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val)
}
