package task

import (
	"fmt"
	"sync"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IbcTxRelateHistoryTask struct {
}

var _ Task = new(IbcTxRelateHistoryTask)
var relateHistoryCoordinator *chainQueueCoordinator

func (t *IbcTxRelateHistoryTask) Name() string {
	return "ibc_tx_relate_history_task"
}

func (t *IbcTxRelateHistoryTask) Cron() int {
	if taskConf.CronTimeIbcTxRelateTask > 0 {
		return taskConf.CronTimeIbcTxRelateTask
	}
	return ThreeMinute
}

func (t *IbcTxRelateHistoryTask) workerNum() int {
	if global.Config.Task.IbcTxRelateWorkerNum > 0 {
		return global.Config.Task.IbcTxRelateWorkerNum
	}
	return ibcTxRelateTaskWorkerNum
}

func (t *IbcTxRelateHistoryTask) Run() int {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", t.Name(), err)
		return -1
	}

	// init coordinator
	chainQueue := new(utils.QueueString)
	for _, v := range chainMap {
		chainQueue.Push(v.ChainId)
	}
	relateHistoryCoordinator = &chainQueueCoordinator{
		chainQueue: chainQueue,
	}

	workerNum := t.workerNum()
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	for i := 1; i <= workerNum; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			newIbcTxRelateWorker(t.Name(), wn, ibcTxTargetHistory, chainMap).exec()
			waitGroup.Done()
		}(workName)
	}
	waitGroup.Wait()

	return 1
}
