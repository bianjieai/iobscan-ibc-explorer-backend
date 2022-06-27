package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

var toKenRepo repository.ITokenRepo = new(repository.TokenRepo)
var chainRepo repository.IChainRepo = new(repository.IbcChainRepo)
var relayerRepo repository.IRelayerRepo = new(repository.IbcRelayerRepo)
var relayerCfgRepo repository.IRelayerConfigRepo = new(cache.RelayerConfigCacheRepo)
