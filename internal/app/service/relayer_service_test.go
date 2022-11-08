package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestRelayerService_List(t *testing.T) {
	resp, err := new(RelayerService).List(&vo.RelayerListReq{
		Page: vo.Page{
			PageNum:  1,
			PageSize: 10,
		},
		UseCount: false,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(resp)
}

func Test_getMsgAmtDenom(t *testing.T) {
	denom := getMsgAmtDenom(&model.TxMsg{
		Type: "acknowledge_packet",
		Msg: bson.M{
			"packet_id": "transferchannel-164transferchannel-2392",
			"packet": bson.M{
				"source_port":         "transfer",
				"source_channel":      "channel-164",
				"destination_port":    "transfer",
				"destination_channel": "channel-239",
				"data": bson.M{
					"denom":    "transfer/channel-152/uwwtoken0913001",
					"amount":   2,
					"sender":   "iaa1n6yx57k3rp2pfk8rekakfahy740ny60s94c3c9",
					"receiver": "cosmos1hfqy9knpa35te0hz7f2xpy99kz8ljh4sux92fz",
				},
			},
			"signer": "iaa1n9wuxk2d69xt0q996rdetewqg6uwd8rmhdz0a3",
		},
	})
	t.Log(denom)
}

func TestRelayerService_DetailRelayerTxs(t *testing.T) {
	global.Config.App.MaxPageSize = 100
	res, err := new(RelayerService).DetailRelayerTxs("6364f740177ccd71260b3fa0", &vo.DetailRelayerTxsReq{
		Page: vo.Page{
			PageNum:  1,
			PageSize: 10,
		},
		Chain: "irishub_qa",
	})
	if err != nil {
		t.Fatal(err.Msg())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(res)))
}

func TestRelayerService_Detail(t *testing.T) {
	res, err := new(RelayerService).Detail("6364f740177ccd71260b3fa0")
	if err != nil {
		t.Fatal(err.Msg())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(res)))
}
