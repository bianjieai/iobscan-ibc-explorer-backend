package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"testing"
)

var (
	task IbcRelayerCronTask
)

func TestIbcRelayerCronTask_Run(t *testing.T) {
	task.Run()
}

func Test_updateRegisterRelayerChannelPairInfo(t *testing.T) {
	//updateRegisterRelayerChannelPairInfo()
}

func TestIbcRelayerCronTask_getUpdateTime(t *testing.T) {
	one, err := relayerRepo.FindOneByRelayerId("6364f740177ccd71260b3fa0")
	if err != nil {
		t.Fatal(err.Error())
	}
	data := task.getUpdateTime(one)
	t.Log("updateTime:", data)
}

func Test_caculateRelayerTotalValue(t *testing.T) {
	denomPrice := cache.TokenPriceMap()
	one, err := relayerRepo.FindOneByRelayerId("63743d45e2427f9a04d8f42b")
	if err != nil {
		t.Fatal(err.Error())
	}
	txsAmt := AggrRelayerTxsAndAmt(one)
	txsAmtValue := caculateRelayerTotalValue(denomPrice, txsAmt)
	feeAmt := AggrRelayerFeeAmt(one)
	feeAmtValue := caculateRelayerTotalValue(denomPrice, feeAmt)
	t.Log(txsAmtValue)
	t.Log(feeAmtValue)
}
