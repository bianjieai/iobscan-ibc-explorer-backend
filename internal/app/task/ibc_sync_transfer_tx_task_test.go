package task

import "testing"

func Test_TraceDenom(t *testing.T) {
	chainMap, _ := getAllChainMap()
	denom := traceDenom("transfer/channel-97/transfer/channel-10", "bigbang", chainMap)
	//denom := worker.traceDenom("uiris", "uiris", "bigbang")

	t.Log(denom)
}

func Test_SyncTransferTx(t *testing.T) {
	new(IbcSyncTransferTxTask).Run()
}
