package task

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"

var (
	taskConf conf.Task
)

func LoadTaskConf(taskCfg conf.Task) {
	taskConf = taskCfg
}
