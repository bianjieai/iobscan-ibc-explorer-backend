package lcd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

var lcdTxDataCacheRepo cache.LcdTxDataCacheRepo

const (
	replaceHolderAddress = "{address}"
)

func GetAccount(chain, address, lcd, apiPath string, crossCache bool) (*vo.AccountResp, error) {
	lcdGet := func() (*vo.AccountResp, error) {
		apiPath = strings.ReplaceAll(apiPath, replaceHolderAddress, address)
		url := fmt.Sprintf("%s%s", lcd, apiPath)

		bz, err := utils.HttpGet(url)
		if err != nil {
			return nil, err
		}

		var resp vo.AccountResp
		err = json.Unmarshal(bz, &resp)
		if err != nil {
			return nil, err
		}

		_ = lcdTxDataCacheRepo.SetAccount(chain, address, &resp)
		return &resp, err
	}

	if crossCache { // 绕过缓存，取链上的最新数据
		return lcdGet()
	}

	if state, err := lcdTxDataCacheRepo.GetAccount(chain, address); err == nil {
		return state, nil
	}

	return lcdGet()
}

func GetBalances(chain, address, lcd, apiPath string) (*vo.BalancesResp, error) {
	if state, err := lcdTxDataCacheRepo.GetBalances(chain, address); err == nil {
		return state, nil
	}
	apiPath = strings.ReplaceAll(apiPath, replaceHolderAddress, address)
	url := fmt.Sprintf("%s%s", lcd, apiPath)

	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp vo.BalancesResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return nil, err
	}

	_ = lcdTxDataCacheRepo.SetBalances(chain, address, &resp)

	return &resp, nil
}

func GetUnbonding(chain, address, lcd, apiPath string) (*vo.UnbondingResp, error) {
	if state, err := lcdTxDataCacheRepo.GetUnbonding(chain, address); err == nil {
		return state, nil
	}
	apiPath = strings.ReplaceAll(apiPath, replaceHolderAddress, address)
	url := fmt.Sprintf("%s%s", lcd, apiPath)

	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp vo.UnbondingResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return nil, err
	}
	_ = lcdTxDataCacheRepo.SetUnbonding(chain, address, &resp)
	return &resp, nil
}

func GetDelegation(chain, address, lcd, apiPath string) (*vo.DelegationResp, error) {
	if state, err := lcdTxDataCacheRepo.GetDelegation(chain, address); err == nil {
		return state, nil
	}
	apiPath = strings.ReplaceAll(apiPath, replaceHolderAddress, address)
	url := fmt.Sprintf("%s%s", lcd, apiPath)

	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp vo.DelegationResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return nil, err
	}
	_ = lcdTxDataCacheRepo.SetDelegation(chain, address, &resp)
	return &resp, nil
}

func GetRewards(chain, address, lcd, apiPath string) (*vo.RewardsResp, error) {
	if state, err := lcdTxDataCacheRepo.GetRewards(chain, address); err == nil {
		return state, nil
	}
	apiPath = strings.ReplaceAll(apiPath, replaceHolderAddress, address)
	url := fmt.Sprintf("%s%s", lcd, apiPath)

	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp vo.RewardsResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return nil, err
	}
	_ = lcdTxDataCacheRepo.SetRewards(chain, address, &resp)
	return &resp, nil
}
