package task

import "testing"

func Test_TraceDenom(t *testing.T) {
	_ = new(IbcSyncTransferTxTask).initCoordinator()
	worker := newSyncTransferTxWorker("sync_task", "worker-1")
	denom := worker.traceDenom("ibc/CC878E18B687447AC4D3670130E1A464DC8B8BDE7C76F5203298D22C2637C793", "transfer/channel-97/transfer/channel-10", "bigbang")
	//denom := worker.traceDenom("uiris", "uiris", "bigbang")

	t.Log(denom)
}

func Test_SyncTransferTx(t *testing.T) {
	new(IbcSyncTransferTxTask).Run()
}
