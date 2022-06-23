package task

import (
	"time"
)

type ChannelTask struct {
}

func (t *ChannelTask) Name() string {
	return "ibc_channel_task"
}

func (t *ChannelTask) Cron() string {
	return ThreeMinute
}

func (t *ChannelTask) ExpireTime() time.Duration {
	return 3*time.Minute - 1*time.Second
}

func (t *ChannelTask) Run() {

}
