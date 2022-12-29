package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

// OverviewCacheRepo overview api cache
type OverviewCacheRepo struct {
}

func (repo *OverviewCacheRepo) SetTokenDistribution(denom string, chain string, value vo.TokenDistributionResp) error {
	bytes, _ := json.Marshal(value)
	_, err := rc.HSet(fmt.Sprintf(overviewTokenDistribution, chain, denom), chain, bytes)
	rc.Expire(fmt.Sprintf(overviewTokenDistribution, chain, denom), oneMin)
	return err
}

func (repo *OverviewCacheRepo) GetTokenDistribution(denom string, chain string) (*vo.TokenDistributionResp, error) {
	var res *vo.TokenDistributionResp
	err := rc.UnmarshalHGet(fmt.Sprintf(overviewTokenDistribution, chain, denom), chain, &res)
	return res, err
}
