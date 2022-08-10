package cache

import "fmt"

//缓存chain_id链上的client_id与channel对的对应关系
type ClientIdChannelCacheRepo struct {
}

func (repo *ClientIdChannelCacheRepo) Set(chainId, clientId, hashVal string) error {
	_, err := rc.HSet(fmt.Sprintf(clientIdInfo, chainId), clientId, hashVal)
	return err
}

func (repo *ClientIdChannelCacheRepo) Get(chainId, clientId string) (string, error) {
	return rc.HGet(fmt.Sprintf(clientIdInfo, chainId), clientId)
}
