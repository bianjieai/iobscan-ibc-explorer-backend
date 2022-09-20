package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
)

type CacheService struct {
}

func (svc *CacheService) Del(key string) (int64, errors.Error) {
	num, err := cache.RedisDel(key)
	if err != nil {
		return num, errors.Wrap(err)
	}
	return num, nil
}
