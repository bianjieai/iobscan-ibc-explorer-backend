package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

const (
	EveryMinute         = 60
	ThreeMinute         = 180
	RedisLockExpireTime = 300
	OneOffTaskLockTime  = 86400 * 30

	statisticsCheckTimes = 5
)

const (
	opInsert = 1
	opUpdate = 2
)

var (
	//cache
	tokenPriceRepo      cache.TokenPriceCacheRepo
	denomDataRepo       cache.DenomDataCacheRepo
	ibcInfoHashCache    cache.IbcInfoHashCacheRepo
	ibcInfoCache        cache.IbcInfoCacheRepo
	unbondTimeCache     cache.UnbondTimeCacheRepo
	statisticsCheckRepo cache.StatisticsCheckCacheRepo
	relayerCache        cache.RelayerCacheRepo
	chainCache          cache.ChainCacheRepo
	baseDenomCache      cache.BaseDenomCacheRepo

	// mongo
	tokenRepo                repository.ITokenRepo                = new(repository.TokenRepo)
	tokenTraceRepo           repository.ITokenTraceRepo           = new(repository.TokenTraceRepo)
	tokenStatisticsRepo      repository.ITokenStatisticsRepo      = new(repository.TokenStatisticsRepo)
	tokenTraceStatisticsRepo repository.ITokenTraceStatisticsRepo = new(repository.TokenTraceStatisticsRepo)
	baseDenomRepo            repository.IBaseDenomRepo            = new(repository.BaseDenomRepo)
	denomRepo                repository.IDenomRepo                = new(repository.DenomRepo)
	denomCaculateRepo        repository.IDenomCaculateRepo        = new(repository.DenomCaculateRepo)
	chainConfigRepo          repository.IChainConfigRepo          = new(repository.ChainConfigRepo)
	ibcTxRepo                repository.IExIbcTxRepo              = new(repository.ExIbcTxRepo)
	chainRepo                repository.IChainRepo                = new(repository.IbcChainRepo)
	relayerRepo              repository.IRelayerRepo              = new(repository.IbcRelayerRepo)
	txRepo                   repository.ITxRepo                   = new(repository.TxRepo)
	channelRepo              repository.IChannelRepo              = new(repository.ChannelRepo)
	channelStatisticsRepo    repository.IChannelStatisticsRepo    = new(repository.ChannelStatisticsRepo)
	channelConfigRepo        repository.IChannelConfigRepo        = new(repository.ChannelConfigRepo)
	relayerStatisticsRepo    repository.IRelayerStatisticsRepo    = new(repository.RelayerStatisticsRepo)
	relayerStatisticsTask    RelayerStatisticsTask
)
