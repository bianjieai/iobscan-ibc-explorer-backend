package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/task"
)

var (
	tokenService    service.ITokenService    = new(service.TokenService)
	channelService  service.IChannelService  = new(service.ChannelService)
	chainService    service.IChainService    = new(service.ChainService)
	addressService  service.IAddressService  = new(service.AddressService)
	relayerService  service.IRelayerService  = new(service.RelayerService)
	homeService     service.IHomeService     = new(service.HomeService)
	transferService service.ITransferService = new(service.TransferService)
	overviewService service.IOverviewService = new(service.OverviewService)
	cacheService    service.CacheService

	// task
	addChainTask              task.AddChainTask
	tokenStatisticsTask       task.TokenStatisticsTask
	channelStatisticsTask     task.ChannelStatisticsTask
	relayerStatisticsTask     task.RelayerStatisticsTask
	addTransferDataTask       task.AddTransferDataTask
	ibcNodeLcdCronTask        task.IbcNodeLcdCronTask
	ibcStatisticCronTask      task.IbcStatisticCronTask
	fixRelayerStatisticsTask  task.FixRelayerStatisticsTask
	relayerAddressInitTask    task.IbcRelayerAddressInitTask
	chainInflowStatisticsTask task.ChainInflowStatisticsTask
)
