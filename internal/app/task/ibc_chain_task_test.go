package task

import (
	"context"
	"testing"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
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
		//Url: "mongodb://ibc:ibcpassword@192.168.0.122:27017,192.168.0.126:27017,192.168.0.127:27017/?authSource=iobscan-ibc",
		Url: "mongodb://ibc:ibcpassword@192.168.150.60:27018/?connect=direct&authSource=iobscan-ibc",
		//Url:      "mongodb://ibc:ibcpassword@35.236.185.62:38129/?connect=direct&authSource=iobscan-ibc",

		Database: "iobscan-ibc",
	}, context.Background())

	time.Local = time.UTC
	global.Config = &conf.Config{Task: conf.Task{
		SingleChainSyncTransferTxMax:      1000,
		SingleChainIbcTxRelateMax:         1000,
		FixDenomTraceDataStartTime:        1634081359,
		FixDenomTraceDataEndTime:          1658814309,
		FixDenomTraceHistoryDataStartTime: 1620369550,
		FixDenomTraceHistoryDataEndTime:   1658830692,
	}}
	m.Run()
}

func TestRunOnce(t *testing.T) {
	new(IbcChainCronTask).Run()
}
