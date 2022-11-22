package rest

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"

var (
	chainService      service.IChainService         = new(service.ChainService)
	staticInfoService service.IStatisticInfoService = new(service.StatisticInfoService)
	ibcTxService      service.IbcTxServerI          = new(service.IbcTxService)
)
