package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"testing"
)

func TestRun(t *testing.T) {
	global.Config.Task.SwitchAddTransferDataTask = true
	global.Config.ChainConfig.AddTransferChains = "columbus_5"
	new(AddTransferDataTask).Run()
}

func TestStart(t *testing.T) {
	chainMap, err := getAllChainMap()
	if err != nil {

	}
	w := &syncTransferTxWorker{
		chainMap: chainMap,
	}
	new(AddTransferDataTask).DoChain(w, "columbus_5", 0, 100)
}
