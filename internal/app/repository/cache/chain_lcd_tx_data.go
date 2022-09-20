package cache

import "fmt"

//缓存配置的lcd相关信息
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
