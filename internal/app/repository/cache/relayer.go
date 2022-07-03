package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	v8 "github.com/go-redis/redis/v8"
)

type RelayerCacheRepo struct {
	relayer repository.IbcRelayerRepo
}

func (repo *RelayerCacheRepo) FindAll() ([]*entity.IBCRelayer, error) {
	value, err := rc.Get(ibcRelayer)
	if err != nil && err == v8.Nil || len(value) == 0 {
		relayers, err := repo.relayer.FindAll(0, 0)
		if err != nil {
			return nil, err
		}
		if len(relayers) > 0 {
			_ = rc.Set(ibcRelayer, utils.MarshalJsonIgnoreErr(relayers), FiveMin)
			return relayers, nil
		}
	}
	var data []*entity.IBCRelayer
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}
