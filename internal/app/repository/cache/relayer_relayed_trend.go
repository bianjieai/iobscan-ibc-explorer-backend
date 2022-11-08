package cache

//RelayerRelayedTrendCacheRepo 缓存从relayer最近N天的交易总数
type RelayerRelayedTrendCacheRepo struct {
}

func (repo *RelayerRelayedTrendCacheRepo) Set(relayerId, days string, value string) error {
	hashKey := relayerId + "_" + days
	_, err := rc.HSet(relayerRelayedTrend, hashKey, value)
	rc.Expire(relayerRelayedTrend, 2*oneHour)
	return err
}

func (repo *RelayerRelayedTrendCacheRepo) Get(relayerId, days string) (string, error) {
	hashKey := relayerId + "_" + days
	value, err := rc.HGet(relayerRelayedTrend, hashKey)
	if err != nil {
		return "", err
	}

	return value, nil
}
