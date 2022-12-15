package cache

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

//AddressCacheRepo 缓存从lcd查询的交易相关信息
type AddressCacheRepo struct {
}

func (repo *AddressCacheRepo) SetAccountList(chain, addr string, data *vo.AccountListResp) error {
	return rc.MarshalSet(fmt.Sprintf(addrAccounts, chain, addr), data, oneMin)
}

func (repo *AddressCacheRepo) GetAccountList(chain, addr string) (*vo.AccountListResp, error) {
	var data vo.AccountListResp
	if err := rc.UnmarshalGet(fmt.Sprintf(addrAccounts, chain, addr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *AddressCacheRepo) SetTokenList(chain, addr string, data *vo.AddrTokenListResp) error {
	return rc.MarshalSet(fmt.Sprintf(addrTokens, chain, addr), data, oneMin)
}

func (repo *AddressCacheRepo) GetTokenList(chain, addr string) (*vo.AddrTokenListResp, error) {
	var data vo.AddrTokenListResp
	if err := rc.UnmarshalGet(fmt.Sprintf(addrTokens, chain, addr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
