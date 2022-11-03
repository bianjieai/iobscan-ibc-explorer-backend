package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
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

	syncTransferTxTaskWorkerNum    = 5
	ibcTxRelateTaskWorkerNum       = 5
	relayerStatisticsWorkerNum     = 4
	fixDenomTraceDataTaskWorkerNum = 8
	fixIbxTxWorkerNum              = 5
	defaultMaxHandlerTx            = 2000
	ibcTxTargetLatest              = "latest"
	ibcTxTargetHistory             = "history"

	segmentStepLatest  = 24 * 3600
	segmentStepHistory = 12 * 3600
)
const (
	channelMatchSuccess = 1
	channelNotFound     = 0
	channelMatchFail    = -1
)

var (
	//cache
	tokenPriceRepo      cache.TokenPriceCacheRepo
	denomDataRepo       cache.DenomDataCacheRepo
	unbondTimeCache     cache.UnbondTimeCacheRepo
	statisticsCheckRepo cache.StatisticsCheckCacheRepo
	chainCache          cache.ChainCacheRepo
	baseDenomCache      cache.BaseDenomCacheRepo
	storageCache        cache.StorageCacheRepo
	lcdTxDataCacheRepo  cache.LcdTxDataCacheRepo

	// mongo
	tokenRepo                  repository.ITokenRepo                  = new(repository.TokenRepo)
	tokenTraceRepo             repository.ITokenTraceRepo             = new(repository.TokenTraceRepo)
	tokenStatisticsRepo        repository.ITokenStatisticsRepo        = new(repository.TokenStatisticsRepo)
	tokenTraceStatisticsRepo   repository.ITokenTraceStatisticsRepo   = new(repository.TokenTraceStatisticsRepo)
	baseDenomRepo              repository.IBaseDenomRepo              = new(repository.BaseDenomRepo)
	denomRepo                  repository.IDenomRepo                  = new(repository.DenomRepo)
	denomCalculateRepo         repository.IDenomCalculateRepo         = new(repository.DenomCalculateRepo)
	chainConfigRepo            repository.IChainConfigRepo            = new(repository.ChainConfigRepo)
	ibcTxRepo                  repository.IExIbcTxRepo                = new(repository.ExIbcTxRepo)
	chainRepo                  repository.IChainRepo                  = new(repository.IbcChainRepo)
	relayerRepo                repository.IRelayerRepo                = new(repository.IbcRelayerRepo)
	txRepo                     repository.ITxRepo                     = new(repository.TxRepo)
	channelRepo                repository.IChannelRepo                = new(repository.ChannelRepo)
	channelStatisticsRepo      repository.IChannelStatisticsRepo      = new(repository.ChannelStatisticsRepo)
	channelConfigRepo          repository.IChannelConfigRepo          = new(repository.ChannelConfigRepo)
	relayerStatisticsRepo      repository.IRelayerStatisticsRepo      = new(repository.RelayerStatisticsRepo)
	relayerFeeStatisticsRepo   repository.IRelayerFeeStatisticsRepo   = new(repository.RelayerFeeStatisticsRepo)
	relayerDenomStatisticsRepo repository.IRelayerDenomStatisticsRepo = new(repository.RelayerDenomStatisticsRepo)
	relayerAddressChannelRepo  repository.IRelayerAddressChannelRepo  = new(repository.RelayerAddressChannelRepo)
	statisticsRepo             repository.IStatisticRepo              = new(repository.IbcStatisticRepo)
	taskRecordRepo             repository.ITaskRecordRepo             = new(repository.TaskRecordRepo)
	syncTaskRepo               repository.ISyncTaskRepo               = new(repository.SyncTaskRepo)
	syncBlockRepo              repository.ISyncBlockRepo              = new(repository.SyncBlockRepo)
	txNewRepo                  repository.ITxNewRepo                  = new(repository.TxNewRepo)
	chainRegistryRepo          repository.IChainRegistryRepo          = new(repository.ChainRegistryRepo)
	relayerStatisticsTask      RelayerStatisticsTask
)

type chainQueueCoordinator struct {
	chainQueue *utils.QueueString
}

func (coordinator *chainQueueCoordinator) getChain() (string, error) {
	if coordinator.chainQueue == nil {
		return "", fmt.Errorf("coordinator or chain queue is nil")
	}

	return coordinator.chainQueue.Pop()
}
