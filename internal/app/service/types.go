package service

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"

var (
	txRepo            repository.ITxRepo            = new(repository.TxRepo)
	chainConfigRepo   repository.IChainConfigRepo   = new(repository.ChainConfigRepo)
	chainRegistryRepo repository.IChainRegistryRepo = new(repository.ChainRegistryRepo)
	ibcTxFailLogRepo  repository.IIBCTxFailLogRepo  = new(repository.IBCTxFailLogRepo)
)
