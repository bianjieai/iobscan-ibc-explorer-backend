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

func Test_RelayerStatisticRunIncrement(t *testing.T) {
	seg := segment{
		StartTime: 1636761600,
		EndTime:   1636847999,
	}
	_ = relayerStatisticsTask.RunIncrement(&seg)
}

func Test_RelayerStatisticsRunWithParam(t *testing.T) {
	relayerStatisticsTask.RunWithParam("kichain_2", 1636761600, 1636847999)
}

func TestIbcRelayerCronTask_getUpdateTime(t *testing.T) {
	one, err := relayerRepo.FindOneByRelayerId("6364f740177ccd71260b3fa0")
	if err != nil {
		t.Fatal(err.Error())
	}
	data := task.getUpdateTime(one)
	t.Log("updateTime:", data)
}

//func TestIbcRelayerCronTask_getChannelsStatus(t *testing.T) {
//	data := task.getChannelsStatus("irishub_1", "cosmoshub_4")
//	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
//}
//
//func TestIbcRelayerCronTask_CheckAndChangeStatus(t *testing.T) {
//	_ = task.init()
//	task.CheckAndChangeRelayer()
//}
//
//func TestIbcRelayerCronTask_cacheIbcChannelRelayer(t *testing.T) {
//
//	task.cacheIbcChannelRelayer()
//	t.Log(task.channelRelayerCnt)
//}
//
//func TestIbcRelayerCronTask_CountRelayerPacketTxs(t *testing.T) {
//	task.AggrRelayerTxsAndAmt()
//	t.Log(task.relayerTxsDataMap)
//}
//
//func TestIbcRelayerCronTask_getChainUnbondTimeFromLcd(t *testing.T) {
//	task.cacheChainUnbondTimeFromLcd()
//}
//
//func TestIbcRelayerCronTask_handleOneRelayerStatusAndTime(t *testing.T) {
//	task.updateOneRelayerUpdateTime(&entity.IBCRelayer{
//		RelayerId:  "cf0fb3209ec3323c539e0e24c44e576d",
//		ChainA:     "irishub_qa",
//		ChainB:     "bigbang",
//		Status:     2,
//		TimePeriod: -1,
//		UpdateTime: 0,
//	}, 1656558771, 146, 0)
//}
//

func Test_caculateRelayerTotalValue(t *testing.T) {
	denomPrice := cache.TokenPriceMap()
	one, err := relayerRepo.FindOneByRelayerId("6364f740177ccd71260b3fa0")
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

//func TestIbcRelayerCronTask_checkAndUpdateEmptyAddr(t *testing.T) {
//	checkAndUpdateRelayerSrcChainAddr()
//}
