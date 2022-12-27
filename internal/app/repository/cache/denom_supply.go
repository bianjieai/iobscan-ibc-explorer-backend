package cache

import (
	"fmt"
)

type DenomDataCacheRepo struct {
}

func (repo *DenomDataCacheRepo) SetSupply(chain, denom, supply string) error {
	_, err := rc.HSet(fmt.Sprintf(denomSupply, chain), denom, supply)
	return err
}

func (repo *DenomDataCacheRepo) BatchSetSupply(chain string, value map[string]string) error {
	_, err := rc.HSet(fmt.Sprintf(denomSupply, chain), value)
	return err
}

func (repo *DenomDataCacheRepo) GetSupply(chain, denom string) (string, error) {
	return rc.HGet(fmt.Sprintf(denomSupply, chain), denom)
}

func (repo *DenomDataCacheRepo) DelSupply(chain string) (int64, error) {
	return rc.Del(fmt.Sprintf(denomSupply, chain))
}

func (repo *DenomDataCacheRepo) SetTransferAmount(chain string, amount map[string]string) error {
	_, err := rc.HSet(fmt.Sprintf(denomTransAmount, chain), amount)
	return err
}

func (repo *DenomDataCacheRepo) GetTransferAmount(chain, denom string) (string, error) {
	return rc.HGet(fmt.Sprintf(denomTransAmount, chain), denom)
}

func (repo *DenomDataCacheRepo) DelTransferAmount(chain string) (int64, error) {
	return rc.Del(fmt.Sprintf(denomTransAmount, chain))
}
