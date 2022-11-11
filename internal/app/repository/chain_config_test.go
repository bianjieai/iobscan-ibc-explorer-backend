package repository

import (
	"context"
	"encoding/json"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"testing"
)

func TestMain(m *testing.M) {
	InitMgo(conf.Mongo{
		Url:      "mongodb://ibc:ibcpassword@192.168.0.135:27017/?authSource=iobscan-ibc",
		Database: "iobscan-ibc",
	}, context.Background())
	m.Run()
}

func TestChainConfigRepo_FindAll(t *testing.T) {
	data, err := new(ChainConfigRepo).FindAll()
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func TestChainConfigRepo_FindOne(t *testing.T) {
	t.Log(ibcDatabase)
	data, err := new(ChainConfigRepo).FindOne("irishub_qa")
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func Test_Update(t *testing.T) {
	conf := entity.ChainConfig{
		ChainId: "qa_iris_snapshot",
		IbcInfo: []*entity.IbcInfo{
			{
				ChainId: "cosmos_4",
				Paths: []*entity.ChannelPath{
					{
						State:     "OPEN",
						PortId:    "transfer",
						ChannelId: "channel-1",
						ChainId:   "cosmos",
						ScChainId: "qa_iris_snapshot",
						Counterparty: entity.CounterParty{
							State:     "OPEN",
							PortId:    "transfer",
							ChannelId: "channel-9",
						},
					},
				},
			},
		},
		IbcInfoHashLcd: "4bda2cdb211fbad4de8dc26ba03abaccc",
	}
	err := new(ChainConfigRepo).UpdateIbcInfo(&conf)
	t.Log(err)
}
