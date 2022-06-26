package rest

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/service"

var (
	tokenService   service.ITokenService   = new(service.TokenService)
	channelService service.IChannelService = new(service.ChannelService)
)
