package monitor

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor/metrics"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	cronTaskStatusMetric metrics.Guage
	redisStatusMetric    metrics.Guage
	TagName              = "taskname"
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

func SetCronTaskStatusMetricValue(taskName string, value float64) {
	if cronTaskStatusMetric != nil {
		cronTaskStatusMetric.With(TagName, taskName).Set(value)
	}
}

func redisClientStatus() {
	for {
		t := time.NewTimer(time.Duration(10) * time.Second)
		select {
		case <-t.C:
			if cache.RedisStatus() {
				redisStatusMetric.Set(float64(1))
			} else {
				redisStatusMetric.Set(float64(-1))
			}
		}
	}
}

func Start(port string) {
	defer func() {
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
	server.Report(func() {
		redisClientStatus()
	})
}
