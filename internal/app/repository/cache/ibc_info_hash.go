package cache

import "time"

type IbcInfoHashCacheRepo struct {
}

func (repo *IbcInfoHashCacheRepo) Set(chainId, hashVal string) error {
	_, err := rc.HSet(ibcInfoHash, chainId, hashVal)
	return err
}

func (repo *IbcInfoHashCacheRepo) Get(chainId string) (string, error) {
	return rc.HGet(ibcInfoHash, chainId)
}

func (repo *IbcInfoHashCacheRepo) SetExpiredTime(expiration time.Duration) {
	rc.Expire(ibcInfoHash, expiration)
}

func (repo *IbcInfoHashCacheRepo) GetAll() (map[string]string, error) {
	var res map[string]string
	err := rc.UnmarshalHGetAll(ibcInfoHash, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
