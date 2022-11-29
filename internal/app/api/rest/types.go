package rest

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"

var (
	txService    service.ITxService    = new(service.TxService)
	chainService service.IChainService = new(service.ChainService)
)
