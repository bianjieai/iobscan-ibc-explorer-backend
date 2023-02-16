package repository

import (
	"context"
	"encoding/json"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"testing"
)

func TestMain(m *testing.M) {
	InitMgo(conf.Mongo{
		Url:      "mongodb://iobscan:iobscanPassword@192.168.150.40:27017/?connect=direct&authSource=iobscan-ibc_0805",
		Database: "iobscan-ibc_0805",
		//Url:      "mongodb://ibcreader:idy45Eth@34.80.213.223:27017/?connect=direct&authSource=iobscan_ibc",
		//Database: "iobscan_ibc",
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
