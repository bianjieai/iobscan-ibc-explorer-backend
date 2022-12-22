package task

import (
	"testing"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
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

func Test_singleSideAddressMatchPair(t *testing.T) {
	pairInfoList := []entity.ChannelPairInfo{{
		PairId:        "7f2b583500b5d09b06aa5b71d8cc4d44",
		ChainA:        "cosmoshub",
		ChainB:        "osmosis",
		ChannelA:      "channel-2",
		ChannelB:      "channel-135",
		ChainAAddress: "cosmos1yvejj22t78s2vfk7slty2d7fs5lkc8rnnt3j9u",
		ChainBAddress: "osmo1ptm9dyyya6erqj9h9ydurtzvla2h6lv9pdwwnx",
	}, {
		PairId:        "ed72954c46710b1d3e3b1073a02716ee",
		ChainA:        "osmosis",
		ChainB:        "",
		ChannelA:      "channel-135",
		ChannelB:      "",
		ChainAAddress: "osmo19vpyrq9wt7hecx6czp7x8lkspnyy64v7zpkl6m",
		ChainBAddress: "",
	}, {
		PairId:        "ed72954c46710b1d3e3b1k73a02716ee",
		ChainA:        "cosmoshub",
		ChainB:        "",
		ChannelA:      "channel-3",
		ChannelB:      "",
		ChainAAddress: "cosmos1yvejj22t78s2vfk7slty2d7fs5lkc8rnnt3j9u",
		ChainBAddress: "",
	}}

	pair, b, err := singleSideAddressMatchPair(pairInfoList)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)
	t.Log(utils.MustMarshalJsonToStr(pair))

	info, _, _ := matchRelayerChannelPairInfo(pair)
	t.Log(utils.MustMarshalJsonToStr(info))
}
