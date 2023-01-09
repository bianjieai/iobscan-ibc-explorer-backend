package distributiontask

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	CronTask interface {
		Name() string
		Cron() string
		BeforeHook() error // init or status's judge in this
		Run()
	}
)

func RunOnce(task CronTask) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			logrus.WithField("err", err).Errorf("task[%s] panic, stack:%v", task.Name(), string(buf))
		}
	}()

	start := time.Now().Unix()
	logrus.Infof("task[%s] start", task.Name())

	if err := task.BeforeHook(); err != nil {
		logrus.WithField("err", err.Error()).
			Errorf("task[%s] %s error", task.Name(), "beforeHooks")
	} else {
		task.Run()
	}
	logrus.Infof("task[%s] end, use %d(second)", task.Name(), time.Now().Unix()-start)
}

func GenTaskId(prefix string) string {
	value := time.Now().Unix()
	hostname, _ := os.Hostname()
	rand.Seed(value)
	return fmt.Sprintf("task_id:%s-%s-%d-%d", prefix, hostname, value, rand.Intn(100))
}
