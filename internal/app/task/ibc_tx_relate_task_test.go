package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func Test_IbxTxRelateTask(t *testing.T) {
	new(IbcTxRelateTask).Run()
}

func Test_IbxTxRelateHistoryTask(t *testing.T) {
	new(IbcTxRelateHistoryTask).Run()
}

func Test_HandlerIbcTxs(t *testing.T) {
	chainMap, _ := getAllChainMap()
	w := newSyncTransferTxWorker("transfer", "worker", chainMap)
	chainId := "irishub_qa"
	denomMap, _ := w.getChainDenomMap(chainId)
	hashes := []string{"6BDD5E93A3E9DEC5402D8674508A15C52FC80105089DADA896B1AC67F65D275C"}
	txList, _ := txRepo.GetTxByHashes(chainId, hashes)
	ibcTxList, _ := w.handleSourceTx(chainId, txList, denomMap)

	rw := newIbcTxRelateWorker("relate", "worker", ibcTxTargetLatest, chainMap)
	rw.handlerIbcTxs(chainId, ibcTxList, denomMap)
	t.Log(utils.MustMarshalJsonToStr(ibcTxList))
}

func Test_HandlerIbcTxs2(t *testing.T) {
	tx, err := ibcTxRepo.FindByRecordId("8b139dadcbdceff9c3b924087ca16d86", false)
	if err != nil {
		t.Fatal(err)
	}

	chainMap, _ := getAllChainMap()
	chainId := "irishub_qa"
	denomMap := make(map[string]*entity.IBCDenom)

	rw := newIbcTxRelateWorker("relate", "worker", ibcTxTargetLatest, chainMap)
	rw.handlerIbcTxs(chainId, []*entity.ExIbcTx{tx}, denomMap)
}
