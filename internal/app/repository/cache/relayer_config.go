package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	v8 "github.com/go-redis/redis/v8"
)

type RelayerConfigCacheRepo struct {
	relayerCfg repository.RelayerConfigRepo
}

func (repo *RelayerConfigCacheRepo) Set(value string) error {
	err := rc.Set(ibcRelayerCfg, value, oneHour)
	return err
}

func (repo *RelayerConfigCacheRepo) FindAll() ([]*entity.IBCRelayerConfig, error) {
	value, err := rc.Get(ibcRelayerCfg)
	if err != nil && err == v8.Nil {
		relayerCfgs, err := repo.relayerCfg.FindAll()
		if err != nil {
			return nil, err
		}

		_ = rc.Set(ibcRelayerCfg, string(utils.MarshalJsonIgnoreErr(relayerCfgs)), oneHour)
		return relayerCfgs, nil
	}
	datas := make([]*entity.IBCRelayerConfig, 0, 10)
	utils.UnmarshalJsonIgnoreErr([]byte(value), &datas)
	return datas, nil
}
