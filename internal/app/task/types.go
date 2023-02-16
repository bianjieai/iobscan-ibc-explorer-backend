package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

const (
	EveryMinute          = 60
	ThreeMinute          = 180
	EveryHour            = 3600
	OneDay               = 86400
	RedisLockExpireTime  = 300
	OneOffTaskLockTime   = 86400 * 30
	ThreeHourCronJobTime = "0 0 */6 * * ?"
	statisticsCheckTimes = 5

	opInsert = 1
	opUpdate = 2

	segmentStep = 24 * 3600
)

var (
	chainConfigRepo        repository.IChainConfigRepo        = new(repository.ChainConfigRepo)
	chainFeeStatisticsRepo repository.IChainFeeStatisticsRepo = new(repository.ChainFeeStatisticsRepo)
	ibcTxRepo              repository.IExIbcTxRepo            = new(repository.ExIbcTxRepo)
	txRepo                 repository.ITxRepo                 = new(repository.TxRepo)
	ibcTxFailLogRepo       repository.IIBCTxFailLogRepo       = new(repository.IBCTxFailLogRepo)
)

type stringQueueCoordinator struct {
	stringQueue *utils.QueueString
}

func (coordinator *stringQueueCoordinator) getOne() (string, error) {
	if coordinator.stringQueue == nil {
		return "", fmt.Errorf("coordinator or string queue is nil")
	}

	return coordinator.stringQueue.Pop()
}
