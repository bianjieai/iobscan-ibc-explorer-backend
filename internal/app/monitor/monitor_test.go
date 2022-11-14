package monitor

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
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

func Test_checkAndUpdateLcd(t *testing.T) {
	lcd := "https://api.sifchain.finance:443"
	chainId := "sifchain_1"
	ok := checkAndUpdateLcd(lcd, &entity.ChainConfig{
		ChainId: chainId,
	})
	t.Log(ok)
}