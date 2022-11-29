package task

import (
	"fmt"
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
)

const (
	opInsert = 1
	opUpdate = 2

	ibcTxCount = 500000

	fixCreateAtErrTime = 1656950400

	replaceHolderOffset  = "OFFSET"
	replaceHolderLimit   = "LIMIT"
	replaceHolderChannel = "CHANNEL"
	replaceHolderPort    = "PORT"
)
const (
	channelMatchSuccess = 1
	channelNotFound     = 0
	channelMatchFail    = -1
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
