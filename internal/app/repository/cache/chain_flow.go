package cache

import (
	"fmt"
)

// ChainFlowCacheRepo chain ibc in/out cache
type ChainFlowCacheRepo struct {
}

func (repo *ChainFlowCacheRepo) GetInflowVolume(days int, chain string) (string, error) {
	return rc.HGet(fmt.Sprintf(chainInflowVolume, days), chain)
}

func (repo *ChainFlowCacheRepo) GetAllInflowVolume(days int) (map[string]float64, error) {
	var res map[string]float64
	err := rc.UnmarshalHGetAll(fmt.Sprintf(chainInflowVolume, days), &res)
	return res, err
}

func (repo *ChainFlowCacheRepo) GetOutflowVolume(days int, chain string) (string, error) {
	return rc.HGet(fmt.Sprintf(chainOutflowVolume, days), chain)
}

func (repo *ChainFlowCacheRepo) GetAllOutflowVolume(days int) (map[string]float64, error) {
	var res map[string]float64
	err := rc.UnmarshalHGetAll(fmt.Sprintf(chainOutflowVolume, days), &res)
	return res, err
}
