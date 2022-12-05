package distributiontask

import (
	"fmt"
	v8 "github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DistributedTask struct {
	l          LockHasExpired
	expiration time.Duration
	tasks      []CronTask
}

const defaultExpiration = 30 * time.Second

func NewDistributedTask(l LockHasExpired) *DistributedTask {
	return &DistributedTask{l: l}
}

func NewDistributedTaskWithRedis(addrs, user, password, mode string, db int) (*DistributedTask, error) {
	var redisClient v8.UniversalClient
	if mode == "cluster" {
		redisClient = v8.NewClusterClient(&v8.ClusterOptions{
			Addrs:    strings.Split(addrs, ","),
			Username: user,
			Password: password,
		})
	} else if mode == "single" {
		redisClient = v8.NewClient(&v8.Options{
			Addr:     addrs,
			Username: user,
			Password: password,
			DB:       db,
		})
	} else {
		logrus.Errorf("unknown redis server mode, %s", mode)
		return nil, fmt.Errorf("redis server mode must be cluster or single")
	}
	return &DistributedTask{l: &RedisLock{redisClient: redisClient}}, nil
}

func (d *DistributedTask) Expiration() time.Duration {
	if d.expiration <= 10 {
		return defaultExpiration
	}
	return d.expiration
}

func (d *DistributedTask) SetExpiration(expiration time.Duration) {
	d.expiration = expiration
}

func (d *DistributedTask) RegisterTasks(task ...CronTask) {
	d.tasks = append(d.tasks, task...)
}

func (d *DistributedTask) Start() {
	if len(d.tasks) == 0 {
		return
	}
	if d.l == nil {
		logrus.Warn("you need implements LockHasExpired interface")
		d.l = new(WithoutLock)
	}

	c := cron.New(cron.WithSeconds())
	for _, v := range d.tasks {
		t := v
		_, err := c.AddFunc(t.Cron(), func() {
			d.RunOnceWithLock(t)
		})
		if err != nil {
			logrus.WithField("task", t.Name()).
				WithField("err", err.Error()).
				Fatal("add cron job err", err)
		}
	}
	c.Run()
}

func (d *DistributedTask) RunOnceWithLock(task CronTask) {
	value := GenTaskId(task.Name())
	if err := d.l.Lock(task.Name(), value, d.Expiration()); err != nil {
		logrus.WithField("key", task.Name()).
			WithField("err", err.Error()).
			Warnf("redis lock failed")

		return
	}

	stop := make(chan struct{})
	t := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-t.C:
				d.ReNewExpiration(task.Name(), value)
			case <-stop:
				logrus.WithField("key", task.Name()).Debug("expired")
				return
			}
		}
	}()

	RunOnce(task)
	stop <- struct{}{}
	if _, err := d.l.UnLock(task.Name(), value); err != nil {
		logrus.WithField("key", task.Name()).
			WithField("err", err.Error()).
			Error("unlock failed")
	}
}

func (d *DistributedTask) ReNewExpiration(key string, value interface{}) {
	ttl, err := d.l.TTL(key)
	if err != nil {
		logrus.WithField("key", key).Error("get ttl failed")
		return
	}
	if ttl == 0 {
		return
	}

	if ttl <= d.Expiration()/3 {
		if _, err := d.l.ScriptExpire(key, value, d.Expiration()); err != nil {
			logrus.WithField("key", key).
				WithField("err", err.Error()).
				Error("expired failed")
		}
		logrus.WithField("key", key).Debug("renew")
	}
}
