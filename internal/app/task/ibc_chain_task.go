package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type IbcChainCronTask struct {
}

func (t *IbcChainCronTask) Name() string {
	return "ibc_chain_task"
}
func (t *IbcChainCronTask) Cron() string {
	return EveryMinute
}
func (t *IbcChainCronTask) Run() {
	chainCfgs, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %s", t.Name(), err.Error())
		return
	}
	var chains []entity.IBCChain
	//set redis key expired time
	for _, chainCfg := range chainCfgs {
		hashVal, _ := ibcInfoHashCache.Get(chainCfg.ChainId)
		//check hashValLcd if have change for reduce update or insert times
		if hashVal == chainCfg.IbcInfoHashLcd {
			continue
		}
		_ = ibcInfoHashCache.Set(chainCfg.ChainId, chainCfg.IbcInfoHashLcd)
		conntectedChains := len(chainCfg.IbcInfo)
		channels := 0
		for _, val := range chainCfg.IbcInfo {
			channels += len(val.Paths)
			//save cache
			if len(val.Paths) > 0 {
				_ = ibcInfoCache.Set(chainCfg.ChainId, val.ChainId, string(utils.MarshalJsonIgnoreErr(val.Paths)))
			}
		}
		data := createChainData(chainCfg.ChainId, channels, conntectedChains)
		chains = append(chains, data)
	}

	for _, val := range chains {
		if err := chainRepo.InserOrUpdate(val); err != nil {
			logrus.Errorf("ibc_chain inser or update fail, %s", err.Error())
		}
	}

}
func createChainData(chainId string, channels int, conntectedChains int) entity.IBCChain {
	return entity.IBCChain{
		ChainId:         chainId,
		Channels:        int64(channels),
		ConnectedChains: int64(conntectedChains),
		CreateAt:        time.Now().Unix(),
		UpdateAt:        time.Now().Unix(),
	}
}
func (t *IbcChainCronTask) ExpireTime() time.Duration {
	return 1*time.Minute - 1*time.Second
}
