package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
	"time"
)

type IbcChainCronTask struct {
}

func (t *IbcChainCronTask) Name() string {
	return "ibc_chain_task"
}
func (t *IbcChainCronTask) Cron() int {
	if taskConf.CronTimeChainTask > 0 {
		return taskConf.CronTimeChainTask
	}
	return EveryMinute
}
func (t *IbcChainCronTask) Run() int {
	chainCfgs, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %s", t.Name(), err.Error())
		return -1
	}
	var chains []entity.IBCChain
	for _, chainCfg := range chainCfgs {
		conntectedChains := len(chainCfg.IbcInfo)
		channels := 0
		for _, val := range chainCfg.IbcInfo {
			channels += len(val.Paths)
		}
		data := createChainData(chainCfg.CurrentChainId, channels, conntectedChains)
		chains = append(chains, data)
	}

	for _, val := range chains {
		if err := chainRepo.InserOrUpdate(val); err != nil {
			logrus.Errorf("ibc_chain inser or update fail, %s", err.Error())
		}
	}
	return 1

}
func createChainData(chainId string, channels int, conntectedChains int) entity.IBCChain {
	return entity.IBCChain{
		Chain:           chainId,
		Channels:        int64(channels),
		ConnectedChains: int64(conntectedChains),
		CreateAt:        time.Now().Unix(),
		UpdateAt:        time.Now().Unix(),
	}
}
func (t *IbcChainCronTask) ExpireTime() time.Duration {
	return 1*time.Minute - 1*time.Second
}
