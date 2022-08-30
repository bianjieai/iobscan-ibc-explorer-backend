package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type AddDataTask struct {
}

var _ OneOffTask = new(AddDataTask)

func (t *AddDataTask) Name() string {
	return "add_data_task"
}

func (t *AddDataTask) Switch() bool {
	return global.Config.Task.SwitchAddDataTask
}

func (t *AddDataTask) Run() int {
	chainMap, err := getAllChainMap()
	if err != nil {
		return -1
	}
	w := &syncTransferTxWorker{
		chainMap: chainMap,
	}

	doChain := func(chainId string, height, limit int64) (int64, int, error) {
		maxHeight := int64(-1)
		denomMap, err := w.getChainDenomMap(chainId)
		if err != nil {
			return maxHeight, 0, err
		}
		txList, err := w.getTxList(chainId, height, limit)
		if err != nil {
			return maxHeight, 0, err
		}
		total := len(txList)
		if err := t.handleChain(chainId, w, txList, denomMap); err != nil {
			return maxHeight, 0, err
		}
		if len(txList) > 0 {
			maxHeight = txList[len(txList)-1].Height
		}
		return maxHeight, total, nil
	}

	chainCureight := make(map[string]int64, len(chainMap))
	for _, val := range chainMap {
		logrus.Info("start handle chain:", val.ChainId)
		for {
			curH, size, err := doChain(val.ChainId, chainCureight[val.ChainId], defaultMaxHandlerTx)
			if err != nil {
				logrus.Error(err.Error())
				return -1
			}
			if size < defaultMaxHandlerTx {
				logrus.Info("finish handle chain:", val.ChainId)
				break
			}
			chainCureight[val.ChainId] = curH
		}
	}
	return 1
}

func (t *AddDataTask) handleChain(chainId string, w *syncTransferTxWorker, txList []*entity.Tx, denomMap map[string]*entity.IBCDenom) error {
	if len(txList) == 0 {
		return nil
	}

	ibcTxList, ibcDenomList := w.handleSourceTx(chainId, txList, denomMap)
	if len(ibcDenomList) > 0 {
		if err := denomRepo.InsertBatch(ibcDenomList); err != nil {
			logrus.Errorf("task %s worker %s denomRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chainId, err)
			return err
		}
	}
	if len(ibcTxList) > 0 {
		if err := ibcTxRepo.InsertBatch(ibcTxList); err != nil {
			logrus.Errorf("task %s worker %s ibcTxRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chainId, err)
			return err
		}
	}
	return nil
}
