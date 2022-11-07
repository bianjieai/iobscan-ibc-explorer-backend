package cache

//no use code
type UnbondTimeCacheRepo struct {
}

func (repo *UnbondTimeCacheRepo) SetUnbondTime(chainId, unbondTime string) error {
	_, err := rc.HSet(chainUnbondTime, chainId, unbondTime)
	_ = rc.Expire(chainUnbondTime, oneDay)
	return err
}

func (repo *UnbondTimeCacheRepo) GetUnbondTime(chainId string) (string, error) {
	return rc.HGet(chainUnbondTime, chainId)
}
