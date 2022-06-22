package task

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"testing"
)

func TestMain(m *testing.M) {
	repository.InitMgo(conf.Mongo{
		Url:      "mongodb://ibc:ibcpassword@192.168.150.60:27018/?connect=direct&authSource=iobscan-ibc",
		Database: "iobscan-ibc",
	}, context.Background())
	m.Run()
}

func TestRunOnce(t *testing.T) {
	new(IbcChainCronTask).Run()
}
