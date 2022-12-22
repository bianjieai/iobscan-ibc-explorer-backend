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
