package cache

import "fmt"

type StorageCacheRepo struct {
}

func (repo *StorageCacheRepo) AddMissDenom(recordId, chainId, denom string) error {
	_, err := rc.SAdd(missDenom, fmt.Sprintf("%s|%s|%s", recordId, chainId, denom))
	return err
}
