package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
)

type RelayerConfigCacheRepo struct {
	relayerCfg repository.RelayerConfigRepo
}

func (repo *RelayerConfigCacheRepo) FindAll() ([]*entity.IBCRelayerConfig, error) {
	fn := func() (interface{}, error) {
		return repo.relayerCfg.FindAll()
	}

	var res []*entity.IBCRelayerConfig
	err := rc.StringTemplateUnmarshal(ibcRelayerCfg, oneHour, fn, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *RelayerConfigCacheRepo) FindRelayerPairIds() ([]*dto.RelayerPairIdDTO, error) {
	fn := func() (interface{}, error) {
		return repo.relayerCfg.FindRelayerPairIds()
	}

	var res []*dto.RelayerPairIdDTO
	err := rc.StringTemplateUnmarshal(ibcRelayerCfgPairIds, oneHour, fn, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *RelayerConfigCacheRepo) Insert(cfg *entity.IBCRelayerConfig) error {
	return repo.relayerCfg.Insert(cfg)
}

func (repo *RelayerConfigCacheRepo) DelCacheFindAll() (int64, error) {
	return rc.Del(ibcRelayerCfg)
}

func (repo *RelayerConfigCacheRepo) DelCacheFindRelayerPairIds() (int64, error) {
	return rc.Del(ibcRelayerCfgPairIds)
}
