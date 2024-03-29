package service

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"testing"
)

func TestMain(m *testing.M) {
	cache.InitRedisClient(conf.Redis{
		Addrs:    "127.0.0.1:6379",
		User:     "",
		Password: "",
		Mode:     "single",
		Db:       0,
	})
	repository.InitMgo(conf.Mongo{
		Url:      "mongodb://ibc:ibcpassword@192.168.0.135:27017/?authSource=iobscan-ibc",
		Database: "iobscan-ibc",
	}, context.Background())
	global.Config = &conf.Config{
		App: conf.App{
			MaxPageSize: 100,
		},
	}
	m.Run()
}

func TestChainService_List(t *testing.T) {
	resp, err := new(ChainService).List(&vo.ChainListReq{})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(resp)
}
