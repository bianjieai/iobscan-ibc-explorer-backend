package cache

import "strconv"

//RelayerTotalTxsCacheRepo 缓存从relayer交易列表查询的交易总数
type RelayerTotalTxsCacheRepo struct {
}

func (repo *RelayerTotalTxsCacheRepo) Set(relayerId, chain string, value int64) error {
	hashKey := relayerId + "_" + chain
	_, err := rc.HSet(ibcRelayerTotalTxs, hashKey, value)
	rc.Expire(ibcRelayerTotalTxs, 2*oneHour)
	return err
}

func (repo *RelayerTotalTxsCacheRepo) Get(relayerId, chain string) (int64, error) {
	hashKey := relayerId + "_" + chain
	value, err := rc.HGet(ibcRelayerTotalTxs, hashKey)
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}
