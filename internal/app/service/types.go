package service

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"

var (
	txRepo                        repository.ITxRepo                     = new(repository.TxRepo)
	chainConfigRepo               repository.IChainConfigRepo            = new(repository.ChainConfigRepo)
	chainRegistryRepo             repository.IChainRegistryRepo          = new(repository.ChainRegistryRepo)
	ibcTxFailLogRepo              repository.IIBCTxFailLogRepo           = new(repository.IBCTxFailLogRepo)
	ibcChainInflowStatisticsRepo  repository.IChainInflowStatisticsRepo  = new(repository.ChainInflowStatisticsRepo)
	ibcChainOutflowStatisticsRepo repository.IChainOutflowStatisticsRepo = new(repository.ChainOutflowStatisticsRepo)
	ibcChainFeeStatisticsRepo     repository.IChainFeeStatisticsRepo     = new(repository.ChainFeeStatisticsRepo)
	authDenomRepo                 repository.IAuthDenomRepo              = new(repository.AuthDenomRepo)
	ibcChainAddressStatisticsRepo repository.IAddressStatisticsRepo      = new(repository.AddressStatisticsRepo)
	ibcChainRepo                  repository.IChainRepo                  = new(repository.IbcChainRepo)
)
