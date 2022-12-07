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
