package cache

import (
	"encoding/json"
	"strconv"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

// RelayerDataCacheRepo 缓存relayer各项数据
type RelayerDataCacheRepo struct {
}

func (repo *RelayerDataCacheRepo) SetRelayedTrend(relayerId, days string, value string) error {
	hashKey := relayerId + "_" + days
	_, err := rc.HSet(relayerRelayedTrend, hashKey, value)
	if err != nil {
		return err
	}
	rc.Expire(relayerRelayedTrend, 2*oneHour)
	return nil
}

func (repo *RelayerDataCacheRepo) GetRelayedTrend(relayerId, days string) (string, error) {
	hashKey := relayerId + "_" + days
	value, err := rc.HGet(relayerRelayedTrend, hashKey)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (repo *RelayerDataCacheRepo) DelRelayedTrend() error {
	_, err := rc.Del(relayerRelayedTrend)
	return err
}

func (repo *RelayerDataCacheRepo) SetTotalTxs(relayerId, chain string, value int64) error {
	hashKey := relayerId + "_" + chain
	_, err := rc.HSet(ibcRelayerTotalTxs, hashKey, value)
	if err != nil {
		return err
	}
	rc.Expire(ibcRelayerTotalTxs, 2*oneHour)
	return nil
}

func (repo *RelayerDataCacheRepo) GetTotalTxs(relayerId, chain string) (int64, error) {
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

func (repo *RelayerDataCacheRepo) DelTotalTxs() error {
	_, err := rc.Del(ibcRelayerTotalTxs)
	return err
}

func (repo *RelayerDataCacheRepo) SetTransferTypeTxs(relayerId string, data *vo.TransferTypeTxsResp) error {
	bz, _ := json.Marshal(data)
	_, err := rc.HSet(ibcRelayerTransferTypeTxs, relayerId, bz)
	if err != nil {
		return err
	}
	rc.Expire(ibcRelayerTransferTypeTxs, 2*oneHour)
	return nil
}

func (repo *RelayerDataCacheRepo) GetTransferTypeTxs(relayerId string) (*vo.TransferTypeTxsResp, error) {
	value, err := rc.HGet(ibcRelayerTransferTypeTxs, relayerId)
	if err != nil {
		return nil, err
	}

	var res vo.TransferTypeTxsResp
	if err := json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (repo *RelayerDataCacheRepo) DelTransferTypeTxs() error {
	_, err := rc.Del(ibcRelayerTransferTypeTxs)
	return err
}

func (repo *RelayerDataCacheRepo) SetTotalRelayedValue(relayerId string, data *vo.TotalRelayedValueResp) error {
	bz, _ := json.Marshal(data)
	_, err := rc.HSet(ibcRelayerTotalRelayedValue, relayerId, bz)
	if err != nil {
		return err
	}
	rc.Expire(ibcRelayerTotalRelayedValue, 7*oneDay)
	return nil
}

func (repo *RelayerDataCacheRepo) GetTotalRelayedValue(relayerId string) (*vo.TotalRelayedValueResp, error) {
	value, err := rc.HGet(ibcRelayerTotalRelayedValue, relayerId)
	if err != nil {
		return nil, err
	}

	var res vo.TotalRelayedValueResp
	if err := json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (repo *RelayerDataCacheRepo) DelTotalRelayedValue() error {
	_, err := rc.Del(ibcRelayerTotalRelayedValue)
	return err
}

func (repo *RelayerDataCacheRepo) SetTotalFeeCost(relayerId string, data *vo.TotalFeeCostResp) error {
	bz, _ := json.Marshal(data)
	_, err := rc.HSet(ibcRelayerTotalFeeCost, relayerId, bz)
	if err != nil {
		return err
	}
	rc.Expire(ibcRelayerTotalFeeCost, 7*oneDay)
	return nil
}

func (repo *RelayerDataCacheRepo) GetTotalFeeCost(relayerId string) (*vo.TotalFeeCostResp, error) {
	value, err := rc.HGet(ibcRelayerTotalFeeCost, relayerId)
	if err != nil {
		return nil, err
	}

	var res vo.TotalFeeCostResp
	if err := json.Unmarshal([]byte(value), &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (repo *RelayerDataCacheRepo) DelTotalFeeCost() error {
	_, err := rc.Del(ibcRelayerTotalFeeCost)
	return err
}
