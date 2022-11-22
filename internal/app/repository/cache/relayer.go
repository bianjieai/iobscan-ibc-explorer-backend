package cache

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	v8 "github.com/go-redis/redis/v8"
)

type RelayerAddrCacheRepo struct {
	relayer repository.IbcRelayerRepo
}

func (repo *RelayerAddrCacheRepo) FindAll() ([]string, error) {
	value, err := rc.Get(ibcRelayer)
	if err != nil && err == v8.Nil || len(value) == 0 {
		skip := int64(0)
		limit := int64(1000)
		relayAddrs := make([]string, 0, 1000)
		for {
			relayers, err := repo.relayer.FindAllRelayerAddrs(skip, limit)
			if err != nil {
				return nil, err
			}
			for _, val := range relayers {
				addrs := entity.ChannelPairInfoList(val.ChannelPairInfo).GetChainAddrs()
				relayAddrs = append(relayAddrs, addrs...)
			}
			if len(relayers) < int(limit) {
				break
			}
			skip += limit
		}

		if len(relayAddrs) > 0 {
			_ = rc.Set(ibcRelayer, utils.MarshalJsonIgnoreErr(relayAddrs), oneHour)
			return relayAddrs, nil
		}
	}
	var data []string
	utils.UnmarshalJsonIgnoreErr([]byte(value), &data)
	return data, nil
}
