package monitor

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor/metrics"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	cronTaskStatusMetric  metrics.Guage
	lcdConnectStatsMetric metrics.Guage
	redisStatusMetric     metrics.Guage
	TagName               = "taskname"
	ChainTag              = "chain_id"
)

func NewMetricCronWorkStatus() metrics.Guage {
	syncWorkStatusMetric := metrics.NewGuage(
		"ibc_explorer_backend",
		"",
		"cron_task_status",
		"ibc_explorer_backend cron task working status (1:Normal  -1:UNormal)",
		[]string{TagName},
	)
	syncWorkStatus, _ := metrics.CovertGuage(syncWorkStatusMetric)
	return syncWorkStatus
}

func NewMetricRedisStatus() metrics.Guage {
	redisNodeStatusMetric := metrics.NewGuage(
		"ibc_explorer_backend",
		"redis",
		"connection_status",
		"ibc_explorer_backend  node connection status of redis service (1:Normal  -1:UNormal)",
		nil,
	)
	redisStatus, _ := metrics.CovertGuage(redisNodeStatusMetric)
	return redisStatus
}

func NewMetricLcdStatus() metrics.Guage {
	lcdConnectionStatusMetric := metrics.NewGuage(
		"ibc_explorer_backend",
		"lcd",
		"connection_status",
		"ibc_explorer_backend  lcd connection status of blockchain (1:Normal  -1:UNormal)",
		[]string{ChainTag},
	)
	connectionStatus, _ := metrics.CovertGuage(lcdConnectionStatusMetric)
	return connectionStatus
}

func SetCronTaskStatusMetricValue(taskName string, value float64) {
	if cronTaskStatusMetric != nil {
		cronTaskStatusMetric.With(TagName, taskName).Set(value)
	}
}

func lcdConnectionStatus(quit chan bool) {
	for {
		t := time.NewTimer(time.Duration(15) * time.Second)
		select {
		case <-t.C:
			chainCfgs, err := new(repository.ChainConfigRepo).FindAllChainInfs()
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			for _, val := range chainCfgs {
				_, err = utils.HttpGet(fmt.Sprintf("%s/node_info", val.Lcd))
				if err != nil {
					_, err = utils.HttpGet(fmt.Sprintf("%s/blocks/latest", val.Lcd))
					if err != nil {
						lcdConnectStatsMetric.With(ChainTag, val.ChainId).Set(float64(-1))
						logrus.Error(err.Error())
						continue
					}
				}
				lcdConnectStatsMetric.With(ChainTag, val.ChainId).Set(float64(1))
			}

		case <-quit:
			logrus.Debug("quit signal recv  lcdConnectionStatus")
			return

		}
	}
}

func redisClientStatus(quit chan bool) {
	for {
		t := time.NewTimer(time.Duration(10) * time.Second)
		select {
		case <-t.C:
			if cache.RedisStatus() {
				redisStatusMetric.Set(float64(1))
			} else {
				redisStatusMetric.Set(float64(-1))
			}
		case <-quit:
			logrus.Debug("quit signal recv redisClientStatus")
			return
		}
	}
}

func Start(port string) {
	quit := make(chan bool)
	defer func() {
		close(quit)
		if err := recover(); err != nil {
			logrus.Error("monitor server occur error ", err)
			os.Exit(1)
		}
	}()
	logrus.Info("monitor server start")
	// start monitor
	server := metrics.NewMonitor(port)
	cronTaskStatusMetric = NewMetricCronWorkStatus()
	redisStatusMetric = NewMetricRedisStatus()
	lcdConnectStatsMetric = NewMetricLcdStatus()
	server.Report(func() {
		go redisClientStatus(quit)
		go lcdConnectionStatus(quit)
	})
}
