package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

// ChainFlowCacheRepo chain ibc in/out cache
type ChainFlowCacheRepo struct {
}

func (repo *ChainFlowCacheRepo) SetInflowTrend(days int, chain string, value []vo.VolumeItem) error {
	bytes, _ := json.Marshal(value)
	_, err := rc.HSet(fmt.Sprintf(chainInflowVolumeTrend, days), chain, bytes)
	return err
}

func (repo *ChainFlowCacheRepo) GetInflowTrend(days int, chain string) ([]vo.VolumeItem, error) {
	var res []vo.VolumeItem
	err := rc.UnmarshalHGet(fmt.Sprintf(chainInflowVolumeTrend, days), chain, &res)
	return res, err
}

func (repo *ChainFlowCacheRepo) ExpireInflowTrend(days int, expire time.Duration) bool {
	return rc.Expire(fmt.Sprintf(chainInflowVolumeTrend, days), expire)
}

func (repo *ChainFlowCacheRepo) SetInflowVolume(days int, chain string, value string) error {
	_, err := rc.HSet(fmt.Sprintf(chainInflowVolume, days), chain, value)
	return err
}

func (repo *ChainFlowCacheRepo) GetInflowVolume(days int, chain string) (string, error) {
	return rc.HGet(fmt.Sprintf(chainInflowVolume, days), chain)
}

func (repo *ChainFlowCacheRepo) GetAllInflowVolume(days int) (map[string]string, error) {
	var res map[string]string
	err := rc.UnmarshalHGetAll(fmt.Sprintf(chainInflowVolume, days), &res)
	return res, err
}

func (repo *ChainFlowCacheRepo) ExpireInflowVolume(days int, expire time.Duration) bool {
	return rc.Expire(fmt.Sprintf(chainInflowVolume, days), expire)
}

func (repo *ChainFlowCacheRepo) SetOutflowTrend(days int, chain string, value []vo.VolumeItem) error {
	bytes, _ := json.Marshal(value)
	_, err := rc.HSet(fmt.Sprintf(chainOutflowVolumeTrend, days), chain, bytes)
	return err
}

func (repo *ChainFlowCacheRepo) GetOutflowTrend(days int, chain string) ([]vo.VolumeItem, error) {
	var res []vo.VolumeItem
	err := rc.UnmarshalHGet(fmt.Sprintf(chainOutflowVolumeTrend, days), chain, &res)
	return res, err
}

func (repo *ChainFlowCacheRepo) ExpireOutflowTrend(days int, expire time.Duration) bool {
	return rc.Expire(fmt.Sprintf(chainOutflowVolumeTrend, days), expire)
}

func (repo *ChainFlowCacheRepo) SetOutflowVolume(days int, chain string, value string) error {
	_, err := rc.HSet(fmt.Sprintf(chainOutflowVolume, days), chain, value)
	return err
}

func (repo *ChainFlowCacheRepo) GetOutflowVolume(days int, chain string) (string, error) {
	return rc.HGet(fmt.Sprintf(chainOutflowVolume, days), chain)
}

func (repo *ChainFlowCacheRepo) GetAllOutflowVolume(days int) (map[string]string, error) {
	var res map[string]string
	err := rc.UnmarshalHGetAll(fmt.Sprintf(chainOutflowVolume, days), &res)
	return res, err
}

func (repo *ChainFlowCacheRepo) ExpireOutflowVolume(days int, expire time.Duration) bool {
	return rc.Expire(fmt.Sprintf(chainOutflowVolume, days), expire)
}
