package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestIbcRelayerCronTask_handleIbcTxLatest(t *testing.T) {
	new(IbcRelayerCronTask).handleIbcTxLatest(0)
}

func TestIbcRelayerCronTask_Run(t *testing.T) {
	new(IbcRelayerCronTask).Run()
}

func TestIbcRelayerCronTask_getTimePeriodAndupdateTime(t *testing.T) {
	data, data1, err := new(IbcRelayerCronTask).getTimePeriodAndupdateTime(&entity.IBCRelayer{
		ChainA:        "bigbang",
		ChainB:        "irishub_qa",
		ChainAAddress: "cosmos16mrml9n668a6ywxsxvtkdymy9kh5m595ygr6g7",
		ChainBAddress: "iaa1u3tpcx5088rx3lzzt0gg73lt9zugrjp730apj8",
		UpdateTime:    1623218166,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("timePeriod:", data, "updateTime:", data1)
}

func TestIbcRelayerCronTask_getChannelsStatus(t *testing.T) {
	data := new(IbcRelayerCronTask).getChannelsStatus("irishub_1", "cosmoshub_4")
	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
}

func TestIbcRelayerCronTask_CheckAndChangeStatus(t *testing.T) {
	new(IbcRelayerCronTask).CheckAndChangeStatus()
}

func TestIbcRelayerCronTask_cacheIbcChannelRelayer(t *testing.T) {
	task := new(IbcRelayerCronTask)
	task.cacheIbcChannelRelayer()
	t.Log(task.channelRelayerCnt)
}

func TestIbcRelayerCronTask_CountRelayerPacketTxs(t *testing.T) {
	task := new(IbcRelayerCronTask)
	task.CountRelayerPacketTxs()
	task.saveOrUpdateRelayerTxs()
	t.Log(task.relayerTxsMap)
}

func TestIbcRelayerCronTask_CountRelayerPacketTxsAmount(t *testing.T) {
	task := new(IbcRelayerCronTask)
	task.CountRelayerPacketTxsAmount()
	task.saveOrUpdateRelayerTxs()
	//t.Log(task.relayerAmtsMap)
}

func TestIbcRelayerCronTask_caculateRelayerTotalValue(t *testing.T) {
	task := new(IbcRelayerCronTask)
	task.CountRelayerPacketTxsAmount()
	task.caculateRelayerTotalValue()
}

func TestIbcRelayerCronTask_getChainUnbondTimeFromLcd(t *testing.T) {
	task := new(IbcRelayerCronTask)
	task.cacheChainUnbondTimeFromLcd()
}
