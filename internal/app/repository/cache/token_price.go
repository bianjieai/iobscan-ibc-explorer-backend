package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/sirupsen/logrus"
	"strconv"
)

type TokenPriceCacheRepo struct {
}

func (repo *TokenPriceCacheRepo) Set(coinId string, price string) error {
	_, err := rc.HSet(tokenPrice, coinId, price)
	return err
}

func (repo *TokenPriceCacheRepo) BatchSet(price map[string]string) error {
	_, err := rc.HSet(tokenPrice, price)
	return err
}

func (repo *TokenPriceCacheRepo) Get(coinId string) (float64, error) {
	str, err := rc.HGet(tokenPrice, coinId)
	if err != nil {
		return 0, err
	}

	float, _ := strconv.ParseFloat(str, 64)
	return float, nil
}

func (repo *TokenPriceCacheRepo) GetAll() (map[string]float64, error) {
	var res map[string]float64
	err := rc.UnmarshalHGetAll(tokenPrice, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func TokenPriceMap() map[string]dto.CoinItem {
	coinIdPriceMap, _ := new(TokenPriceCacheRepo).GetAll()
	baseDenoms, err := new(BaseDenomCacheRepo).FindAll()
	if err != nil {
		logrus.Error("find base_denom fail, ", err.Error())
		return nil
	}
	if len(coinIdPriceMap) == 0 {
		return nil
	}
	denomPriceMap := make(map[string]dto.CoinItem, len(baseDenoms))
	for _, val := range baseDenoms {
		if price, ok := coinIdPriceMap[val.CoinId]; ok {
			denomPriceMap[val.Denom+val.Chain] = dto.CoinItem{Price: price, Scale: val.Scale}
		}
	}
	return denomPriceMap
}
