package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"
	"sync"
)

var (
	taskConf      conf.Task
	TaskMetricMap = new(sync.Map)
)

func LoadTaskConf(taskCfg conf.Task) {
	taskConf = taskCfg
}
