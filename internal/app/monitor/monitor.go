package monitor

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor/metrics"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/task"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	cronTaskStatusMetric metrics.Guage
	TagName              = "taskname"
)

func NewMetricCronWorkStatus() metrics.Guage {
	syncWorkStatusMetric := metrics.NewGuage(
		"iobscan_ibc",
		"openapi",
		"cron_task_status",
		"ibc_openapi cron task working status (1:Normal  -1:UNormal)",
		[]string{TagName},
	)
	syncWorkStatus, _ := metrics.CovertGuage(syncWorkStatusMetric)
	return syncWorkStatus
}

func clientStatus(quit chan bool) {
	for {
		t := time.NewTimer(time.Duration(10) * time.Second)
		select {
		case <-t.C:
			for _, taskName := range []string{
				new(task.IBCChainFeeStatisticTask).Name(),
			} {
				if value, ok := task.TaskMetricMap.Load(taskName); ok {
					cronTaskStatusMetric.With(TagName, taskName).Set(value.(float64))
				} else {
					cronTaskStatusMetric.With(TagName, taskName).Set(0)
				}

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
	server.Report(func() {
		go clientStatus(quit)
	})
}
