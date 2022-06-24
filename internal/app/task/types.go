package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

const (
	EveryMinute     = "0 */1 * * * ?"
	EverySecond     = "*/1 * * * * ?"
	EveryFiveSecond = "*/5 * * * * ?"
	EveryTenSecond  = "*/10 * * * * ?"
	OneHour         = "0 0 */1 * * ?"
	TwelveHour      = "0 0 */12 * * ?"
	ThreeMinute     = "0 */3 * * * ?"
	FiveMinute      = "0 */5 * * * ?"
	TwentyMinute    = "0 */20 * * * ?"
)

var (
	//cache
	tokenPriceRepo   cache.TokenPriceCacheRepo
	denomDataRepo    cache.DenomDataCacheRepo
	ibcInfoHashCache cache.IbcInfoHashCacheRepo
	ibcInfoCache     cache.IbcInfoCacheRepo

	// mongo
	tokenRepo             repository.ITokenRepo             = new(repository.TokenRepo)
	baseDenomRepo         repository.IBaseDenomRepo         = new(repository.BaseDenomRepo)
	denomRepo             repository.IDenomRepo             = new(repository.DenomRepo)
	denomCaculateRepo     repository.IDenomCaculateRepo     = new(repository.DenomCaculateRepo)
	tokenStatisticsRepo   repository.ITokenStatisticsRepo   = new(repository.TokenStatisticsRepo)
	chainConfigRepo       repository.IChainConfigRepo       = new(repository.ChainConfigRepo)
	ibcTxRepo             repository.IExIbcTxRepo           = new(repository.ExIbcTxRepo)
	chainRepo             repository.IChainRepo             = new(repository.IbcChainRepo)
	relayerRepo           repository.IRelayerRepo           = new(repository.IbcRelayerRepo)
	txRepo                repository.ITxRepo                = new(repository.TxRepo)
	channelRepo           repository.IChannelRepo           = new(repository.ChannelRepo)
	channelStatisticsRepo repository.IChannelStatisticsRepo = new(repository.ChannelStatisticsRepo)
	relayerStatisticsRepo repository.IRelayerStatisticsRepo = new(repository.RelayerStatisticsRepo)
)
