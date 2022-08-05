package cache

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	v8 "github.com/go-redis/redis/v8"
)

type BaseDenomCacheRepo struct {
	baseDenom repository.BaseDenomRepo
}

func (repo *BaseDenomCacheRepo) FindAll() (entity.IBCBaseDenomList, error) {
	value, err := rc.Get(baseDenom)
	if err != nil && err == v8.Nil || len(value) == 0 {
		baseDenoms, err := repo.baseDenom.FindAll()
		if err != nil {
			return nil, err
		}
		if len(baseDenoms) > 0 {
			_ = rc.Set(baseDenom, utils.MarshalJsonIgnoreErr(baseDenoms), oneDay)
			return baseDenoms, nil
		}
	}
	var data []*entity.IBCBaseDenom
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}

func (repo *BaseDenomCacheRepo) FindBySymbol(symbol string) (entity.IBCBaseDenom, error) {
	value, err := rc.Get(fmt.Sprintf(baseDenomSymbol, symbol))
	if err != nil && err == v8.Nil || len(value) == 0 {
		baseDenom, err := repo.baseDenom.FindBySymbol(symbol)
		if err != nil {
			return entity.IBCBaseDenom{}, err
		}
		_ = rc.Set(fmt.Sprintf(baseDenomSymbol, symbol), utils.MarshalJsonIgnoreErr(baseDenom), oneDay)
		return baseDenom, nil
	}
	var data entity.IBCBaseDenom
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}
