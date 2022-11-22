package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
)

var (
	chainRepo       repository.IChainRepo       = new(repository.IbcChainRepo)
	statisticRepo   repository.IStatisticRepo   = new(repository.IbcStatisticRepo)
	ibcTxRepo       repository.IExIbcTxRepo     = new(repository.ExIbcTxRepo)
	txRepo          repository.ITxRepo          = new(repository.TxRepo)
	chainConfigRepo repository.IChainConfigRepo = new(repository.ChainConfigRepo)
)
