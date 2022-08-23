package cache

import (
	"fmt"
	"time"
)

type IbcInfoCacheRepo struct {
}

func (repo *IbcInfoCacheRepo) Set(chainId, dcChainId string, channels string) error {
	_, err := rc.HSet(fmt.Sprintf(ibcInfo, chainId), dcChainId, channels)
	return err
}

func (repo *IbcInfoCacheRepo) Get(chainId, dcChainId string) (interface{}, error) {
	return rc.HGet(fmt.Sprintf(ibcInfo, chainId), dcChainId)
}

func (repo *IbcInfoCacheRepo) SetExpiredTime(chainId string, expiration time.Duration) {
	rc.Expire(fmt.Sprintf(ibcInfo, chainId), expiration)
}
