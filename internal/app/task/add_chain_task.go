package task

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type AddChainTask struct {
}

var _ OneOffTask = new(AddChainTask)

func (t *AddChainTask) Name() string {
	return "add_chain_task"
}

func (t *AddChainTask) Run() int {
	chainsStr := global.Config.ChainConfig.NewChains
	newChainIds := strings.Split(chainsStr, ",")
	if len(newChainIds) == 0 {
		logrus.Errorf("task %s don't have new chains", t.Name())
		return 1
	}

	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", t.Name(), err)
		return -1
	}

	denomList, err := denomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s denomRepo.FindAll error, %v", t.Name(), err)
		return -1
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	// update ibc tx
	go func() {
		for _, chainId := range newChainIds {
			chainConfig, ok := chainMap[chainId]
			if !ok {
				logrus.Warningf("task %s %s dont't have chain config", t.Name(), chainId)
				continue
			}

			t.updateIbcTx(chainId, chainConfig)
		}
		waitGroup.Done()
	}()

	// update denom
	go func() {
		t.updateDenom(denomList, chainMap)
		waitGroup.Done()
	}()

	waitGroup.Wait()
	return 1
}

func (t *AddChainTask) updateIbcTx(chainId string, chainConfig *entity.ChainConfig) {
	logrus.Infof("task %s start updating %s ibc tx", t.Name(), chainId)
	if len(chainConfig.IbcInfo) == 0 {
		logrus.Warningf("task %s %s dont't have ibc info", t.Name(), chainId)
		return
	}

	for _, ibcInfo := range chainConfig.IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.State != constant.ChannelStateOpen {
				logrus.Warningf("task %s %s channel %s is not open", t.Name(), chainId, path.ChannelId)
				continue
			}

			// todo update failed dc chain id??
			counterpartyChainId := path.ChainId
			counterpartyChannelId := path.Counterparty.ChannelId
			if err := ibcTxRepo.AddNewChainUpdate(counterpartyChainId, counterpartyChannelId, chainId); err != nil {
				logrus.Errorf("task %s %s AddNewChainUpdate error, counterpartyChainId: %s, counterpartyChannelId: %s", t.Name(), chainId, counterpartyChainId, counterpartyChannelId)
				_ = storageCache.AddChainError(chainId, counterpartyChainId, counterpartyChannelId)
			}

			if err := ibcTxRepo.AddNewChainUpdateHistory(counterpartyChainId, counterpartyChannelId, chainId); err != nil {
				logrus.Errorf("task %s %s AddNewChainUpdateHistory error, counterpartyChainId: %s, counterpartyChannelId: %s", t.Name(), chainId, counterpartyChainId, counterpartyChannelId)
				_ = storageCache.AddChainError(chainId, counterpartyChainId, counterpartyChannelId)
			}
		}
	}

	logrus.Infof("task %s update %s ibc tx end", t.Name(), chainId)
}

func (t *AddChainTask) updateDenom(denomList entity.IBCDenomList, chainMap map[string]*entity.ChainConfig) {
	logrus.Infof("task %s update denom start", t.Name())

	for _, v := range denomList {
		if v.DenomPath == "" {
			continue
		}

		denomFullPath := fmt.Sprintf("%s/%s", v.DenomPath, v.RootDenom)
		denomNew := traceDenom(denomFullPath, v.ChainId, chainMap)
		if v.BaseDenom != denomNew.BaseDenom || v.BaseDenomChainId != denomNew.BaseDenomChainId || v.PrevDenom != denomNew.PrevDenom ||
			v.PrevChainId != denomNew.PrevChainId || v.IsBaseDenom != denomNew.IsBaseDenom {
			if err := denomRepo.UpdateDenom(denomNew); err != nil {
				logrus.Errorf("task %s update denom %s-%s error", t.Name(), denomNew.ChainId, denomNew.Denom)
			}
		}

		if v.BaseDenom != denomNew.BaseDenom || v.BaseDenomChainId != denomNew.BaseDenomChainId {
			if err := ibcTxRepo.UpdateBaseDenomInfo(v.BaseDenom, v.BaseDenomChainId, denomNew.BaseDenom, denomNew.BaseDenomChainId); err != nil {
				logrus.Errorf("task %s UpdateBaseDenomInfo error, %s-%s => %s-%s", t.Name(), v.BaseDenomChainId, v.BaseDenom, denomNew.BaseDenomChainId, denomNew.BaseDenom)
				_ = storageCache.UpdateBaseDenomError(v.BaseDenom, v.BaseDenomChainId, denomNew.BaseDenom, denomNew.BaseDenomChainId)
			}
			if err := ibcTxRepo.UpdateBaseDenomInfoHistory(v.BaseDenom, v.BaseDenomChainId, denomNew.BaseDenom, denomNew.BaseDenomChainId); err != nil {
				logrus.Errorf("task %s UpdateBaseDenomInfoHistory error, %s-%s => %s-%s", t.Name(), v.BaseDenomChainId, v.BaseDenom, denomNew.BaseDenomChainId, denomNew.BaseDenom)
				_ = storageCache.UpdateBaseDenomError(v.BaseDenom, v.BaseDenomChainId, denomNew.BaseDenom, denomNew.BaseDenomChainId)
			}
		}
	}

	logrus.Infof("task %s update denom end", t.Name())
}
