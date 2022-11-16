package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/task"
)

var (
	tokenService    service.ITokenService    = new(service.TokenService)
	channelService  service.IChannelService  = new(service.ChannelService)
	chainService    service.IChainService    = new(service.ChainService)
	relayerService  service.IRelayerService  = new(service.RelayerService)
	homeService     service.IHomeService     = new(service.HomeService)
	transferService service.ITransferService = new(service.TransferService)
	cacheService    service.CacheService

	// task
	addChainTask                 task.AddChainTask
	fixDcChainIdTask             task.FixDcChainIdTask
	fixBaseDenomChainIdTask      task.FixBaseDenomChainIdTask
	fixDenomTraceDataTask        task.FixDenomTraceDataTask
	fixDenomTraceHistoryDataTask task.FixDenomTraceHistoryDataTask
	tokenStatisticsTask          task.TokenStatisticsTask
	channelStatisticsTask        task.ChannelStatisticsTask
	relayerStatisticsTask        task.RelayerStatisticsTask
	relayerDataTask              task.RelayerDataTask
	fixFailRecvPacketTask        task.FixFailRecvPacketTask
	addTransferDataTask          task.AddTransferDataTask
	fixFailTxTask                task.FixFailTxTask
	fixAcknowledgeTxTask         task.FixAcknowledgeTxTask
	fixAckTxPacketIdTask         task.FixAckTxPacketIdTask
	fixIbxTxTask                 task.FixIbxTxTask
	ibcNodeLcdCronTask           task.IbcNodeLcdCronTask
	ibcStatisticCronTask         task.IbcStatisticCronTask
	modifyChainIdTask            task.ModifyChainIdTask
)
