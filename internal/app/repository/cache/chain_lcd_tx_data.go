package cache

import (
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

//LcdTxDataCacheRepo 缓存从lcd查询的交易相关信息
type LcdTxDataCacheRepo struct {
}

func (repo *LcdTxDataCacheRepo) Set(chain, hash, msgEvents string) error {
	_, err := rc.HSet(fmt.Sprintf(lcdTxData, chain), hash, msgEvents)
	rc.Expire(fmt.Sprintf(lcdTxData, chain), oneHour)
	return err
}

func (repo *LcdTxDataCacheRepo) Get(chain, hashVal string) (string, error) {
	return rc.HGet(fmt.Sprintf(lcdTxData, chain), hashVal)
}

func (repo *LcdTxDataCacheRepo) SetClientState(clientKey string, data *vo.ClientStateResp) error {
	return rc.MarshalSet(fmt.Sprintf(clientState, clientKey), data, oneDay)
}

func (repo *LcdTxDataCacheRepo) GetClientState(clientKey string) (*vo.ClientStateResp, error) {
	var data vo.ClientStateResp
	if err := rc.UnmarshalGet(fmt.Sprintf(clientState, clientKey), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo *LcdTxDataCacheRepo) SetAccount(chain, addr string, data *vo.AccountResp) error {
	return rc.MarshalSet(fmt.Sprintf(lcdAccount, chain, addr), data, threeHours)
}

func (repo *LcdTxDataCacheRepo) GetAccount(chain, addr string) (*vo.AccountResp, error) {
	var data vo.AccountResp
	if err := rc.UnmarshalGet(fmt.Sprintf(lcdAccount, chain, addr), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo *LcdTxDataCacheRepo) SetBalances(chain, addr string, data *vo.BalancesResp) error {
	return rc.MarshalSet(fmt.Sprintf(lcdBalances, chain, addr), data, oneMin)
}

func (repo *LcdTxDataCacheRepo) GetBalances(chain, addr string) (*vo.BalancesResp, error) {
	var data vo.BalancesResp
	if err := rc.UnmarshalGet(fmt.Sprintf(lcdBalances, chain, addr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *LcdTxDataCacheRepo) SetUnbonding(chain, addr string, data *vo.UnbondingResp) error {
	return rc.MarshalSet(fmt.Sprintf(lcdUnbonding, chain, addr), data, oneMin)
}

func (repo *LcdTxDataCacheRepo) GetUnbonding(chain, addr string) (*vo.UnbondingResp, error) {
	var data vo.UnbondingResp
	if err := rc.UnmarshalGet(fmt.Sprintf(lcdUnbonding, chain, addr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *LcdTxDataCacheRepo) SetDelegation(chain, addr string, data *vo.DelegationResp) error {
	return rc.MarshalSet(fmt.Sprintf(lcdDelegation, chain, addr), data, oneMin)
}

func (repo *LcdTxDataCacheRepo) GetDelegation(chain, addr string) (*vo.DelegationResp, error) {
	var data vo.DelegationResp
	if err := rc.UnmarshalGet(fmt.Sprintf(lcdDelegation, chain, addr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *LcdTxDataCacheRepo) SetRewards(chain, addr string, data *vo.RewardsResp) error {
	return rc.MarshalSet(fmt.Sprintf(lcdRewards, chain, addr), data, oneMin)
}

func (repo *LcdTxDataCacheRepo) GetRewards(chain, addr string) (*vo.RewardsResp, error) {
	var data vo.RewardsResp
	if err := rc.UnmarshalGet(fmt.Sprintf(lcdRewards, chain, addr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
