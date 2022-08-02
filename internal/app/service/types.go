package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

var (
	tokenRepo           repository.ITokenRepo         = new(repository.TokenRepo)
	tokenStatisticsRepo repository.ITokenTraceRepo    = new(repository.TokenTraceRepo)
	channelRepo         repository.IChannelRepo       = new(repository.ChannelRepo)
	baseDenomRepo       repository.IBaseDenomRepo     = new(repository.BaseDenomRepo)
	denomRepo           repository.IDenomRepo         = new(repository.DenomRepo)
	chainRepo           repository.IChainRepo         = new(repository.IbcChainRepo)
	relayerRepo         repository.IRelayerRepo       = new(repository.IbcRelayerRepo)
	statisticRepo       repository.IStatisticRepo     = new(repository.IbcStatisticRepo)
	ibcTxRepo           repository.IExIbcTxRepo       = new(repository.ExIbcTxRepo)
	chainConfigRepo     repository.IChainConfigRepo   = new(repository.ChainConfigRepo)
	relayerCfgRepo      repository.IRelayerConfigRepo = new(cache.RelayerConfigCacheRepo)
	logCacheRepo        cache.TxLogCacheRepo
)
