package cache

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	v8 "github.com/go-redis/redis/v8"
)

type AuthDenomCacheRepo struct {
	authDenom repository.AuthDenomRepo
}

func (repo *AuthDenomCacheRepo) FindAll() (entity.AuthDenomList, error) {
	value, err := rc.Get(baseDenom)
	if err != nil && err == v8.Nil || len(value) == 0 {
		authDenoms, err := repo.authDenom.FindAll()
		if err != nil {
			return nil, err
		}
		if len(authDenoms) > 0 {
			_ = rc.Set(baseDenom, utils.MarshalJsonIgnoreErr(authDenoms), oneDay)
			return authDenoms, nil
		}
	}
	var data []*entity.AuthDenom
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}

func (repo *AuthDenomCacheRepo) FindBySymbol(symbol string) (entity.AuthDenom, error) {
	value, err := rc.Get(fmt.Sprintf(baseDenomSymbol, symbol))
	if err != nil && err == v8.Nil || len(value) == 0 {
		authDenom, err := repo.authDenom.FindBySymbol(symbol)
		if err != nil {
			return entity.AuthDenom{}, err
		}
		_ = rc.Set(fmt.Sprintf(baseDenomSymbol, symbol), utils.MarshalJsonIgnoreErr(authDenom), oneDay)
		return authDenom, nil
	}
	var data entity.AuthDenom
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}
