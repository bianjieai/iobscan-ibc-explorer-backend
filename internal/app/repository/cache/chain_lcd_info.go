package cache

//缓存配置的lcd相关信息
type LcdInfoCacheRepo struct {
}

func (repo *LcdInfoCacheRepo) Set(chainId, hashVal string) error {
	_, err := rc.HSet(lcdInfo, chainId, hashVal)
	return err
}

func (repo *LcdInfoCacheRepo) Get(chainId string) (string, error) {
	return rc.HGet(lcdInfo, chainId)
}

func (repo *LcdInfoCacheRepo) GetAll() (map[string]string, error) {
	var res map[string]string
	err := rc.UnmarshalHGetAll(lcdInfo, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
