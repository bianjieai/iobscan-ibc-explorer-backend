package cache

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"time"
)

type TxLogCacheRepo struct {
	txRepo repository.TxRepo
}

func (repo *TxLogCacheRepo) Set(chainId, txHash string, log string) error {
	_, err := rc.HSet(fmt.Sprintf(ibcTxLog, chainId), txHash, log)
	return err
}

func (repo *TxLogCacheRepo) Get(chainId, txHash string) (string, error) {
	return rc.HGet(fmt.Sprintf(ibcTxLog, chainId), txHash)
}

func (repo *TxLogCacheRepo) SetExpiredTime(chainId string, expiration time.Duration) {
	rc.Expire(fmt.Sprintf(ibcTxLog, chainId), expiration)
}

func (repo *TxLogCacheRepo) GetLogByHash(chainId, txHash string) (string, error) {
	value, _ := repo.Get(chainId, txHash)
	if len(value) == 0 {
		tx, err := repo.txRepo.GetTxByHash(chainId, txHash)
		if err != nil {
			return "", err
		}
		_ = repo.Set(chainId, txHash, tx.Log)
		return tx.Log, nil
	}
	return value, nil
}
