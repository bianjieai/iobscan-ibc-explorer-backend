package task

import (
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

func TestIbcRelayerCronTask_getTimePeriodAndupdateTime(t *testing.T) {
	//data := task.getUpdateTime(&entity.IBCRelayerNew{
	//	UpdateTime:    1623218166,
	//})
	//t.Log( "updateTime:", data1)
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
//func TestIbcRelayerCronTask_updateRelayerStatus(t *testing.T) {
//	chainMap, _ := getAllChainMap()
//	task.chainConfigMap = chainMap
//	task.updateRelayerStatus(&entity.IBCRelayer{
//		RelayerId:  "bf8d73cd76b3b6b4a53e1b8c956b7978",
//		ChainA:     "irishub_qa",
//		ChainB:     "bigbang",
//		ChannelA:   "channel-115",
//		ChannelB:   "channel-199",
//		Status:     1,
//		TimePeriod: -1,
//		UpdateTime: 1660705368,
//	})
//}
//

func Test_caculateRelayerTotalValue(t *testing.T) {
	denomPrice := getTokenPriceMap()
	one, err := relayerRepo.FindOneByRelayerId("6364e39e50255e66b989c04d")
	if err != nil {
		t.Fatal(err.Error())
	}
	txsAmt := AggrRelayerTxsAndAmt(one)
	caculateRelayerTotalValue(denomPrice, txsAmt)
	feeAmt := AggrRelayerFeeAmt(one)
	caculateRelayerTotalValue(denomPrice, feeAmt)
}

//func TestIbcRelayerCronTask_checkAndUpdateEmptyAddr(t *testing.T) {
//	checkAndUpdateRelayerSrcChainAddr()
//}
