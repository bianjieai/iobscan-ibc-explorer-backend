package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestOverviewService_TokenDistribution(t *testing.T) {
	resp, err := new(OverviewService).TokenDistribution(&vo.TokenDistributionReq{
		BaseDenomChain: "cosmoshub",
		BaseDenom:      "uatom",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(resp)))
}

func TestOverviewService_ChainVolume(t *testing.T) {
	resp, err := new(OverviewService).ChainVolume(&vo.ChainVolumeReq{})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(resp)))
}

func TestOverviewService_ChainVolumeTrend(t *testing.T) {
	resp, err := new(OverviewService).ChainVolumeTrend(&vo.ChainVolumeTrendReq{})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(resp)))
}
