package cache

import "fmt"

type StorageCacheRepo struct {
}

func (repo *StorageCacheRepo) AddMissDenom(recordId, chainId, denom string) error {
	_, err := rc.SAdd(missDenom, fmt.Sprintf("%s|%s|%s", recordId, chainId, denom))
	return err
}

func (repo *StorageCacheRepo) AddChainError(chainId, counterpartyChainId, counterpartyChannelId string) error {
	_, err := rc.SAdd(addChainError, fmt.Sprintf("%s|%s|%s", chainId, counterpartyChainId, counterpartyChannelId))
	return err
}

func (repo *StorageCacheRepo) UpdateBaseDenomError(baseDenom, baseDenomChainId, baseDenomNew, baseDenomChainIdNew string) error {
	_, err := rc.SAdd(updateBaseDenomError, fmt.Sprintf("%s|%s|%s|%s", baseDenomChainId, baseDenom, baseDenomChainIdNew, baseDenomNew))
	return err
}
