package rest

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"

var (
	tokenService    service.ITokenService    = new(service.TokenService)
	channelService  service.IChannelService  = new(service.ChannelService)
	chainService    service.IChainService    = new(service.ChainService)
	relayerService  service.IRelayerService  = new(service.RelayerService)
	homeService     service.IHomeService     = new(service.HomeService)
	transferService service.ITransferService = new(service.TransferService)
	cacheService    service.CacheService
)
