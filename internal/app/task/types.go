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

	syncTransferTxTaskWorkerNum = 5
	ibcTxRelateTaskWorkerNum    = 5
	relayerStatisticsWorkerNum  = 4
	defaultMaxHandlerTx         = 2000
	ibcTxTargetLatest           = "latest"
	ibcTxTargetHistory          = "history"

	segmentStepLatest  = 24 * 3600
	segmentStepHistory = 12 * 3600

	relayerAddressGatherRangeTime = 7 * 3600
)

var (
	//cache
	tokenPriceRepo             cache.TokenPriceCacheRepo
	denomDataRepo              cache.DenomDataCacheRepo
	relayerDataCache           cache.RelayerDataCacheRepo
	statisticsCheckRepo        cache.StatisticsCheckCacheRepo
	chainCache                 cache.ChainCacheRepo
	baseDenomCache             cache.AuthDenomCacheRepo
	lcdTxDataCacheRepo         cache.LcdTxDataCacheRepo
	tokenRepo                  repository.ITokenRepo                  = new(repository.TokenRepo)
	tokenTraceRepo             repository.ITokenTraceRepo             = new(repository.TokenTraceRepo)
	tokenStatisticsRepo        repository.ITokenStatisticsRepo        = new(repository.TokenStatisticsRepo)
	tokenTraceStatisticsRepo   repository.ITokenTraceStatisticsRepo   = new(repository.TokenTraceStatisticsRepo)
	baseDenomRepo              repository.IAuthDenomRepo              = new(repository.AuthDenomRepo)
	denomRepo                  repository.IDenomRepo                  = new(repository.DenomRepo)
	chainConfigRepo            repository.IChainConfigRepo            = new(repository.ChainConfigRepo)
	ibcTxRepo                  repository.IExIbcTxRepo                = new(repository.ExIbcTxRepo)
	chainRepo                  repository.IChainRepo                  = new(repository.IbcChainRepo)
	relayerRepo                repository.IRelayerRepo                = new(repository.IbcRelayerRepo)
	txRepo                     repository.ITxRepo                     = new(repository.TxRepo)
	channelRepo                repository.IChannelRepo                = new(repository.ChannelRepo)
	channelStatisticsRepo      repository.IChannelStatisticsRepo      = new(repository.ChannelStatisticsRepo)
	channelConfigRepo          repository.IChannelConfigRepo          = new(repository.ChannelConfigRepo)
	relayerFeeStatisticsRepo   repository.IRelayerFeeStatisticsRepo   = new(repository.RelayerFeeStatisticsRepo)
	relayerDenomStatisticsRepo repository.IRelayerDenomStatisticsRepo = new(repository.RelayerDenomStatisticsRepo)
	relayerAddressChannelRepo  repository.IRelayerAddressChannelRepo  = new(repository.RelayerAddressChannelRepo)
	relayerAddressRepo         repository.IRelayerAddressRepo         = new(repository.RelayerAddressRepo)
	statisticsRepo             repository.IStatisticRepo              = new(repository.IbcStatisticRepo)
	taskRecordRepo             repository.ITaskRecordRepo             = new(repository.TaskRecordRepo)
	syncTaskRepo               repository.ISyncTaskRepo               = new(repository.SyncTaskRepo)
	syncBlockRepo              repository.ISyncBlockRepo              = new(repository.SyncBlockRepo)
	txNewRepo                  repository.ITxNewRepo                  = new(repository.TxNewRepo)
	chainRegistryRepo          repository.IChainRegistryRepo          = new(repository.ChainRegistryRepo)
	relayerStatisticsTask      RelayerStatisticsTask
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
