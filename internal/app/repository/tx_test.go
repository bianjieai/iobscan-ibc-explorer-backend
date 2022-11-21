package repository

import "testing"

func TestTxRepo_GetTimePeriodByUpdateClient(t *testing.T) {
	val1, err := new(TxRepo).GetUpdateTimeByUpdateClient("irishub_qa", "iaa1u3tpcx5088rx3lzzt0gg73lt9zugrjp730apj8", "adb", 1656557855)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(val1)
}

func TestTxRepo_GetChannelOpenConfirmTime(t *testing.T) {
	val, err := new(TxRepo).GetChannelOpenConfirmTime("bigbang", "channel-182")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val)
}
