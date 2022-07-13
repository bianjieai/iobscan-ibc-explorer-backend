package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

var (
	task IbcRelayerCronTask
)

func TestIbcRelayerCronTask_Run(t *testing.T) {
	task.Run()
}

func TestIbcRelayerCronTask_getTimePeriodAndupdateTime(t *testing.T) {
	data, data1 := task.getTimePeriodAndupdateTime(&entity.IBCRelayer{
		ChainA:        "bigbang",
		ChainB:        "irishub_qa",
		ChainAAddress: "cosmos16mrml9n668a6ywxsxvtkdymy9kh5m595ygr6g7",
		ChainBAddress: "iaa1u3tpcx5088rx3lzzt0gg73lt9zugrjp730apj8",
		UpdateTime:    1623218166,
	})
	t.Log("timePeriod:", data, "updateTime:", data1)
}

func TestIbcRelayerCronTask_getChannelsStatus(t *testing.T) {
	data := task.getChannelsStatus("irishub_1", "cosmoshub_4")
	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
}

func TestIbcRelayerCronTask_CheckAndChangeStatus(t *testing.T) {
	task.CheckAndChangeRelayer()
}

func TestIbcRelayerCronTask_cacheIbcChannelRelayer(t *testing.T) {

	task.cacheIbcChannelRelayer()
	t.Log(task.channelRelayerCnt)
}

func TestIbcRelayerCronTask_CountRelayerPacketTxs(t *testing.T) {
	task.AggrRelayerPacketTxs()
	t.Log(task.relayerTxsDataMap)
}

func TestIbcRelayerCronTask_getChainUnbondTimeFromLcd(t *testing.T) {
	task.cacheChainUnbondTimeFromLcd()
}

func TestIbcRelayerCronTask_handleOneRelayerStatusAndTime(t *testing.T) {
	task.handleOneRelayerStatusAndTime(&entity.IBCRelayer{
		RelayerId:  "cf0fb3209ec3323c539e0e24c44e576d",
		ChainA:     "irishub_qa",
		ChainB:     "bigbang",
		Status:     2,
		TimePeriod: -1,
		UpdateTime: 0,
	}, 1656558771, 146)
}

func TestIbcRelayerCronTask_updateRelayerStatus(t *testing.T) {
	task.updateRelayerStatus(&entity.IBCRelayer{
		RelayerId:  "cf0fb3209ec3323c539e0e24c44e576d",
		Status:     2,
		TimePeriod: 146,
		UpdateTime: 1656558771,
	})
}

func TestIbcRelayerCronTask_DistinctRelayer(t *testing.T) {
	var datas = []entity.IBCRelayer{
		{ChainA: "irishub_1", ChainB: "cosmoshub_4", ChannelA: "channel-12", ChannelB: "channel-182", ChainAAddress: "iaa148zzqgulnly3wgx35s5f0z4l4vwf30tj03wgaq", ChainBAddress: "cosmos148zzqgulnly3wgx35s5f0z4l4vwf30tj6nwel3"},
		{ChainA: "cosmoshub", ChainB: "irishub_1", ChannelA: "channel-182", ChannelB: "channel-12", ChainAAddress: "cosmos148zzqgulnly3wgx35s5f0z4l4vwf30tj6nwel3", ChainBAddress: "iaa148zzqgulnly3wgx35s5f0z4l4vwf30tj03wgaq"},
	}
	relayerStatisticsTask.initdistRelayerMap()
	value := distinctRelayer(datas, relayerStatisticsTask.distRelayerMap)
	t.Log(value)
}

func TestIbcRelayerCronTask_checkDbExist(t *testing.T) {
	var datas = []entity.IBCRelayer{
		{ChainA: "irishub_1", ChainB: "cosmoshub_4", ChannelA: "channel-12", ChannelB: "channel-182", ChainAAddress: "iaa148zzqgulnly3wgx35s5f0z4l4vwf30tj03wgaq", ChainBAddress: "cosmos148zzqgulnly3wgx35s5f0z4l4vwf30tj6nwel3"},
		{ChainA: "cosmoshub", ChainB: "irishub_1", ChannelA: "channel-182", ChannelB: "channel-12", ChainAAddress: "cosmos148zzqgulnly3wgx35s5f0z4l4vwf30tj6nwel3", ChainBAddress: "iaa148zzqgulnly3wgx35s5f0z4l4vwf30tj03wgaq"},
	}
	relayerStatisticsTask.initdistRelayerMap()
	value := filterDbExist(datas, relayerStatisticsTask.distRelayerMap)
	t.Log(value)
}

func TestIbcRelayerCronTask_caculateRelayerTotalValue(t *testing.T) {
	task.getTokenPriceMap()
	task.caculateRelayerTotalValue()
	for key, val := range task.relayerValueMap {
		t.Log(key, val.String())
	}
}

func TestRelayerStatisticsTask_Run(t *testing.T) {
	new(RelayerStatisticsTask).Run()
}

func TestIbcRelayerCronTask_checkAndUpdateEmptyAddr(t *testing.T) {
	task.checkAndUpdateRelayerSrcChainAddr()
}
