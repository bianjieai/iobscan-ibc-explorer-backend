package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

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

func TestTxRepo_ChainFeeStatistics(t *testing.T) {
	val, err := new(TxRepo).ChainFeeStatistics("cosmoshub_4", 1662249600, 1662335999)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(val)))
}

func TestTxRepo_ChainUserFeeStatistics(t *testing.T) {
	val, err := new(TxRepo).ChainUserFeeStatistics("cosmoshub_4", 1662249600, 1662335999)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(val)))
}
