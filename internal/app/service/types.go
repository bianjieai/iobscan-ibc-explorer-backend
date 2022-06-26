package service

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"

var (
	tokenRepo           repository.ITokenRepo           = new(repository.TokenRepo)
	tokenStatisticsRepo repository.ITokenStatisticsRepo = new(repository.TokenStatisticsRepo)
	channelRepo         repository.IChannelRepo         = new(repository.ChannelRepo)
	baseDenomRepo       repository.IBaseDenomRepo       = new(repository.BaseDenomRepo)
	denomRepo           repository.IDenomRepo           = new(repository.DenomRepo)
)
