package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

var (
	tokenRepo           repository.ITokenRepo          = new(repository.TokenRepo)
	tokenStatisticsRepo repository.ITokenTraceRepo     = new(repository.TokenTraceRepo)
	channelRepo         repository.IChannelRepo        = new(repository.ChannelRepo)
	denomRepo           repository.IDenomRepo          = new(repository.DenomRepo)
	chainRepo           repository.IChainRepo          = new(repository.IbcChainRepo)
	relayerRepo         repository.IRelayerRepo        = new(repository.IbcRelayerRepo)
	statisticRepo       repository.IStatisticRepo      = new(repository.IbcStatisticRepo)
	chainCfgRepo        repository.IChainConfigRepo    = new(repository.ChainConfigRepo)
	ibcTxRepo           repository.IExIbcTxRepo        = new(repository.ExIbcTxRepo)
	txRepo              repository.ITxRepo             = new(repository.TxRepo)
	exSearchRecordRepo  repository.IExSearchRecordRepo = new(repository.ExSearchRecordRepo)
	relayerCfgRepo      repository.IRelayerConfigRepo  = new(cache.RelayerConfigCacheRepo)
	baseDenomRepo       cache.BaseDenomCacheRepo
)
