package task

import (
	"context"
	"testing"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	time.Local = time.UTC
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   constant.DefaultTimeFormat,
		DisableHTMLEscape: true,
	})
	cache.InitRedisClient(conf.Redis{
		Addrs:    "127.0.0.1:6379",
		User:     "",
		Password: "",
		Mode:     "single",
		Db:       0,
	})
	repository.InitMgo(conf.Mongo{
		Url:      "mongodb://iobscan:iobscanPassword@192.168.150.40:27017/?connect=direct&authSource=iobscan-ibc_0805",
		Database: "iobscan-ibc_0805",
	}, context.Background())

	time.Local = time.UTC
	global.Config = &conf.Config{
		Task: conf.Task{},
	}
	m.Run()
}
