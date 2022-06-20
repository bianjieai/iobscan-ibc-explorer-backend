package task

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

type Task interface {
	Name() string
	Cron() string // CronExpression
	Run()
	ExpireTime() time.Duration // redis expireTime
}

var tasks []Task

func RegisterTasks(task ...Task) {
	tasks = append(tasks, task...)
}

// GetTasks get all the task
func GetTasks() []Task {
	return tasks
}

func Start() {
	if len(GetTasks()) == 0 {
		return
	}

	c := cron.New(cron.WithSeconds())
	for _, v := range GetTasks() {
		task := v
		RunOnce(task)

		_, err := c.AddFunc(task.Cron(), func() {
			RunOnceWithLock(task)
		})
		if err != nil {
			logrus.Fatal("cron job err", err)
		}
	}
	c.Start()
}

func RunOnce(task Task) {
	logrus.Infof("task %s start", task.Name())
	task.Run()
	logrus.Infof("task %s end", task.Name())
}

func RunOnceWithLock(task Task) {
	if err := redisClient.Lock(task.Name(), time.Now().Unix(), task.ExpireTime()); err != nil {
		logrus.Errorf("redis lock failed, name:%s, err:%v", task.Name(), err.Error())
		return
	}

	RunOnce(task)
}
