package cache

import (
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

//LcdTxDataCacheRepo 缓存从lcd查询的交易相关信息
type LcdTxDataCacheRepo struct {
}

func (repo *LcdTxDataCacheRepo) Set(chainId, hash, msgEvents string) error {
	_, err := rc.HSet(fmt.Sprintf(lcdTxData, chainId), hash, msgEvents)
	rc.Expire(fmt.Sprintf(lcdTxData, chainId), FiveMin)
	return err
}

func (repo *LcdTxDataCacheRepo) Get(chainId, hashVal string) (string, error) {
	return rc.HGet(fmt.Sprintf(lcdTxData, chainId), hashVal)
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
