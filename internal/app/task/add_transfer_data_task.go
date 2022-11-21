package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
	"strings"
)

type AddTransferDataTask struct {
}

var _ OneOffTask = new(AddTransferDataTask)

func (t *AddTransferDataTask) Name() string {
	return "add_transfer_data_task"
}

func (t *AddTransferDataTask) Switch() bool {
	return false
}

func (t *AddTransferDataTask) Run() int {
	return 1
	//return t.handle(global.Config.ChainConfig.AddTransferChains)
}

func (t *AddTransferDataTask) RunWithParam(chainsStr string) int {
	return t.handle(chainsStr)
}

func (t *AddTransferDataTask) handle(chainsStr string) int {
	newChains := strings.Split(chainsStr, ",")
	if len(newChains) == 0 {
		logrus.Errorf("task %s don't have new chains", t.Name())
		return 1
	}
	chainMap, err := getAllChainMap()
	if err != nil {
		return -1
	}
	w := &syncTransferTxWorker{
		chainMap: chainMap,
	}

	chainCureight := make(map[string]int64, len(newChains))
	for _, val := range newChains {
		logrus.Info("start handle chain:", val)
		for {
			curH, size, err := t.DoChain(w, val, chainCureight[val], defaultMaxHandlerTx)
			if err != nil {
				logrus.Error(err.Error())
				return -1
			}
			if size < defaultMaxHandlerTx {
				logrus.Info("finish handle chain:", val)
				break
			}
			chainCureight[val] = curH
		}
	}
	return 1
}

func (t *AddTransferDataTask) DoChain(w *syncTransferTxWorker, chain string, height, limit int64) (int64, int, error) {
	maxHeight := int64(-1)
	denomMap, err := w.getChainDenomMap(chain)
	if err != nil {
		return maxHeight, 0, err
	}
	transferHashDatas, err := txNewRepo.GetTransferTx(chain, height, limit)
	if err != nil {
		return maxHeight, 0, err
	}
	var hashes []string
	for _, val := range transferHashDatas {
		hashes = append(hashes, val.TxHash)
	}
	txList, err := txRepo.GetTxByHashes(chain, hashes)
	if err != nil {
		return maxHeight, 0, err
	}
	total := len(txList)
	if err := t.handleChain(chain, w, txList, denomMap); err != nil {
		return maxHeight, 0, err
	}
	if len(txList) > 0 {
		maxHeight = txList[len(txList)-1].Height
	}
	return maxHeight, total, nil
}

func (t *AddTransferDataTask) handleChain(chain string, w *syncTransferTxWorker, txList []*entity.Tx, denomMap map[string]*entity.IBCDenom) error {
	if len(txList) == 0 {
		return nil
	}

	ibcTxList, ibcDenomList := w.handleSourceTx(chain, txList, denomMap)
	if len(ibcDenomList) > 0 {
		if err := denomRepo.InsertBatch(ibcDenomList); err != nil {
			logrus.Errorf("task %s worker %s denomRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chain, err)
			return err
		}
	}
	if len(ibcTxList) > 0 {
		if err := ibcTxRepo.InsertBatch(ibcTxList); err != nil {
			logrus.Errorf("task %s worker %s ibcTxRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chain, err)
			return err
		}
	}
	return nil
}
