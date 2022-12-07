package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	v8 "github.com/go-redis/redis/v8"
)

type ChainCacheRepo struct {
	chain repository.IbcChainRepo
}

func (repo *ChainCacheRepo) FindAll() ([]*entity.IBCChain, error) {
	value, err := rc.Get(ibcChain)
	if err != nil && err == v8.Nil || len(value) == 0 {
		chains, err := repo.chain.FindAll(0, 0)
		if err != nil {
			return nil, err
		}
		if len(chains) > 0 {
			_ = rc.Set(ibcChain, utils.MarshalJsonIgnoreErr(chains), FiveMin)
			return chains, nil
		}
	}
	var data []*entity.IBCChain
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}

func (repo *ChainCacheRepo) SetChainsConnection(value string) error {
	err := rc.Set(ibcChainsConnection, value, oneDay)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ChainCacheRepo) GetChainsConnection() (string, error) {
	value, err := rc.Get(ibcChainsConnection)
	if err != nil {
		return "", err
	}
	return value, nil
}
