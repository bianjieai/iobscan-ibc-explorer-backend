package task

import (
	"testing"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

func Test_TraceDenom(t *testing.T) {
	chainMap, _ := getAllChainMap()
	denom := traceDenom("transfer/channel-97/transfer/channel-10", "bigbang", chainMap)
	//denom := worker.traceDenom("uiris", "uiris", "bigbang")

	t.Log(denom)
}

func Test_SyncTransferTx(t *testing.T) {
	new(IbcSyncTransferTxTask).Run()
}

func Test_HandleSourceTx(t *testing.T) {
	chainMap, _ := getAllChainMap()
	w := newSyncTransferTxWorker("transfer", "worker", chainMap)
	chain := "irishub_qa"
	denomMap, _ := w.getChainDenomMap(chain)
	hashes := []string{"3115FB1C39C2156321C175974C9C7EFE9DC5009C2C7A2EF98EA2A70785E45B89", "6BDD5E93A3E9DEC5402D8674508A15C52FC80105089DADA896B1AC67F65D275C"}
	txList, _ := txRepo.GetTxByHashes(chain, hashes)
	ibcTxList, ibcDenomList := w.handleSourceTx(chain, txList, denomMap)
	t.Log(utils.MustMarshalJsonToStr(ibcTxList))
	t.Log(utils.MustMarshalJsonToStr(ibcDenomList))
}
