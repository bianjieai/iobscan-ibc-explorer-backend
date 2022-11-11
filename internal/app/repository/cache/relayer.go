package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
)

type RelayerCacheRepo struct {
	relayer repository.IbcRelayerRepo
}

func (repo *RelayerCacheRepo) FindAll() ([]*entity.IBCRelayerNew, error) {
	fn := func() (interface{}, error) {
		return repo.relayer.FindAllRelayerForCache()
	}

	var res []*entity.IBCRelayerNew
	err := rc.StringTemplateUnmarshal(ibcRelayer, oneHour, fn, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *RelayerCacheRepo) DelCacheFindAll() (int64, error) {
	return rc.Del(ibcRelayer)
}
