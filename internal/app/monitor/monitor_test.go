package monitor

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"testing"
)

func TestMain(m *testing.M) {
	repository.InitMgo(conf.Mongo{
		Url:      "mongodb://ibc:ibcpassword@192.168.0.135:27017/?authSource=iobscan-ibc",
		Database: "iobscan-ibc",
	}, context.Background())
	m.Run()
}

func Test_checkLcd(t *testing.T) {
	value := checkLcd("https://emoney.validator.network/api", "/ibc/core/channel/v1beta1/channels?pagination.offset=OFFSET&pagination.limit=LIMIT&pagination.count_total=true")
	if value {
		t.Log("pass")
	}
}
