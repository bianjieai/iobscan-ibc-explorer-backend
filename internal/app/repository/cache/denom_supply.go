package cache

import (
	"fmt"
)

type DenomSupplyCacheRepo struct {
}

func (repo *DenomSupplyCacheRepo) Set(chainId, denom, supply string) error {
	_, err := rc.HSet(fmt.Sprintf(denomSupply, chainId), denom, supply)
	return err
}

func (repo *DenomSupplyCacheRepo) Get(chainId, denom string) (string, error) {
	return rc.HGet(fmt.Sprintf(denomSupply, chainId), denom)
}
