package task

import (
	"fmt"
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
	chain := "irishub_qa"
	denomMap, _ := w.getChainDenomMap(chain)
	hashes := []string{"6BDD5E93A3E9DEC5402D8674508A15C52FC80105089DADA896B1AC67F65D275C"}
	txList, _ := txRepo.GetTxByHashes(chain, hashes)
	ibcTxList, _ := w.handleSourceTx(chain, txList, denomMap)

	rw := newIbcTxRelateWorker("relate", "worker", ibcTxTargetLatest, chainMap)
	rw.handlerIbcTxs(chain, ibcTxList, denomMap)
	t.Log(utils.MustMarshalJsonToStr(ibcTxList))
}

func Test_parseRecvPacketTxEvents(t *testing.T) {
	txs, err := txRepo.GetTxByHashes("uptick", []string{"C476E603D7A3329FCB8486897B465815F10E3A6B70F7BC89657986409CBC3FB6"})
	if err != nil {
		t.Log(err.Error())
	}

	dcConnection, packetAck, exists := parseRecvPacketTxEvents(1, txs[0])
	fmt.Println(dcConnection)
	fmt.Println(packetAck)
	fmt.Println(exists)
}
