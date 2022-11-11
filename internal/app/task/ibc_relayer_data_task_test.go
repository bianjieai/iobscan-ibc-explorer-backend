package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"testing"
)

func TestRelayerDataTask_Run(t *testing.T) {
	relayerDataTask.Run()
}

func TestRelayerDataTask_doOneSegment(t *testing.T) {
	//segmentData := segment{
	//	StartTime: ,
	//	EndTime: ,
	//}
	//relayerDataTask.doOneSegment(segmentData,true)
}

func TestRelayerDataTask_aggrUnknowRelayerChannelPair(t *testing.T) {
	relayerDataTask.aggrUnknowRelayerChannelPair()
}

func TestRelayerDataTask_matchRegisterRelayerChannelPairInfo(t *testing.T) {
	addrPairs := []entity.ChannelPairInfo{
		{ChainA: "qa_iris_snapshot", ChainAAddress: "iaa1vx32zg7aj62w906cwrjqhpv4xlsx4k4t4l6d2m", ChainB: "bigbang", ChainBAddress: "cosmos1vx32zg7aj62w906cwrjqhpv4xlsx4k4tqa6ug2"},
		{ChainA: "qa_iris_snapshot", ChainAAddress: "iaa1fspgk7hu2ezlpasrf5tw4dwkrxmys8umtpum3a", ChainB: "bigbang", ChainBAddress: "cosmos1fspgk7hu2ezlpasrf5tw4dwkrxmys8um7ru2nv"},
	}
	pairs, _, err := matchRegisterRelayerChannelPairInfo(addrPairs)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(pairs)
}

func Test_getRelayerStatisticData(t *testing.T) {
	one, err := relayerRepo.FindOneByRelayerId("636c60806754c75a88adc0d2")
	if err != nil {
		t.Fatal(err.Error())
	}
	item := getRelayerStatisticData(cache.TokenPriceMap(), one)
	t.Log(item.RelayedTotalTxsValue)
	t.Log(item.TotalFeeValue)
	t.Log(item.RelayedSuccessTxs)
	t.Log(item.RelayedTotalTxs)
}

func Test_doRegisterRelayer(t *testing.T) {
	doRegisterRelayer(cache.TokenPriceMap())
}
