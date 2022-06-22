package cache

import (
	"fmt"
)

type DenomDataCacheRepo struct {
}

func (repo *DenomDataCacheRepo) SetSupply(chainId, denom, supply string) error {
	_, err := rc.HSet(fmt.Sprintf(denomSupply, chainId), denom, supply)
	return err
}

func (repo *DenomDataCacheRepo) GetSupply(chainId, denom string) (string, error) {
	return rc.HGet(fmt.Sprintf(denomSupply, chainId), denom)
}

func (repo *DenomDataCacheRepo) SetTransferAmount(chainId string, amount map[string]string) error {
	_, err := rc.HSet(fmt.Sprintf(denomTransAmount, chainId), amount)
	return err
}

func (repo *DenomDataCacheRepo) GetTransferAmount(chainId, denom string) (string, error) {
	return rc.HGet(fmt.Sprintf(denomTransAmount, chainId), denom)
}
