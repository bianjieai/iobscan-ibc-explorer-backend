package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

const (
	DailyAccountsCronJobTime = "0 0 8 * * ?"
)

var (
	relayerAddrCache cache.RelayerAddrCacheRepo
	chainConfigRepo  repository.IChainConfigRepo = new(repository.ChainConfigRepo)
	txRepo           repository.ITxRepo          = new(repository.TxRepo)
	statisticsRepo   repository.IStatisticRepo   = new(repository.IbcStatisticRepo)
)
