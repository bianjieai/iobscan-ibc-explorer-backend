package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/task"
)

var (
	txService      service.ITxService      = new(service.TxService)
	chainService   service.IChainService   = new(service.ChainService)
	feeService     service.IFeeService     = new(service.FeeService)
	addressService service.IAddressService = new(service.AddressService)
)

var (
	// task
	ibcTxFailLogTask         task.IBCTxFailLogTask
	iBCChainFeeStatisticTask task.IBCChainFeeStatisticTask
	ibcAddressStatisticTask  task.IBCAddressStatisticTask
)
