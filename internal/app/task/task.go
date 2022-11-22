package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Task cron task
type Task interface {
	Name() string
	Cron() int // CronExpression
	Run() int
	//ExpireTime() time.Duration // redis expireTime
}

var (
	taskConf conf.Task
)

func LoadTaskConf(taskCfg conf.Task) {
	taskConf = taskCfg
}

func Start() {

	c := cron.New(cron.WithSeconds())

	if taskConf.CronJobDailyChainAddr == "" {
		taskConf.CronJobDailyChainAddr = DailyAccountsCronJobTime
	}
	if _, err := c.AddFunc(taskConf.CronJobDailyChainAddr, caculateActiveAddrsOfChains); err != nil {
		logrus.Fatal("cron job caculateActiveAddrsOfChains err", err)
	}
	c.Start()
}
