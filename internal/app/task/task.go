package task

import (
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type Task interface {
	Name() string
	Cron() int // CronExpression
	Run() int
	//ExpireTime() time.Duration // redis expireTime
}

var (
	tasks    []Task
	taskConf conf.Task
)

func RegisterTasks(task ...Task) {
	tasks = append(tasks, task...)
}

// GetTasks get all the task
func GetTasks() []Task {
	return tasks
}

func LoadTaskConf(taskCfg conf.Task) {
	taskConf = taskCfg
}

func Start() {
	if len(GetTasks()) == 0 {
		return
	}

	for _, v := range GetTasks() {
		task := v
		go RunOnce(task)
	}
}

func RunOnce(task Task) {
	redisLockExpireTime := time.Duration(RedisLockExpireTime) * time.Second
	if taskConf.RedisLockExpireTime > 0 {
		redisLockExpireTime = time.Duration(taskConf.RedisLockExpireTime) * time.Second
	}

	utils.RunTimer(task.Cron(), utils.Sec, func() {
		//lock redis mux
		if err := cache.GetRedisClient().Lock(task.Name(), time.Now().Unix(), redisLockExpireTime); err != nil {
			logrus.Errorf("redis lock failed, name:%s, err:%v", task.Name(), err.Error())
			return
		}
		startTime := time.Now().Unix()
		logrus.Infof("task %s start", task.Name())
		metricValue := task.Run()
		monitor.SetCronTaskStatusMetricValue(task.Name(), float64(metricValue))
		//unlock redis mux
		cache.GetRedisClient().Del(task.Name())
		logrus.Infof("task %s end, time use %d(s)", task.Name(), time.Now().Unix()-startTime)
	})
}
