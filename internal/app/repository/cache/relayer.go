package cache

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	v8 "github.com/go-redis/redis/v8"
)

type RelayerCacheRepo struct {
	relayer repository.IbcRelayerRepo
}

func (repo *RelayerCacheRepo) BatchSet(price map[string]string) error {
	_, err := rc.HSet(ibcRelayer, price)
	return err
}

func (repo *RelayerCacheRepo) FindAll() (map[string]string, error) {
	value, err := rc.HGetAll(ibcRelayer)
	if err != nil && err == v8.Nil || len(value) == 0 {
		relayers, err := repo.relayer.FindAll(0, 0)
		if err != nil {
			return nil, err
		}
		if len(relayers) > 0 {
			relayerMapCache := make(map[string]string, len(relayers))
			for _, val := range relayers {
				key := fmt.Sprintf("%s:%s:%s", val.ChainA, val.ChainAAddress, val.ChannelA)
				key1 := fmt.Sprintf("%s:%s:%s", val.ChainB, val.ChainBAddress, val.ChannelB)
				relayerMapCache[key] = ""
				relayerMapCache[key1] = ""
			}
			_, _ = rc.HSet(ibcRelayer, relayerMapCache)
			return relayerMapCache, nil
		}
	}
	return value, nil
}
