package task

type IbcSyncTransferTxTask struct {
}

var _ Task = new(IbcSyncTransferTxTask)

func (t *IbcSyncTransferTxTask) Name() string {
	return "ibc_sync_transfer_tx_task"
}
func (t *IbcSyncTransferTxTask) Cron() int {
	if taskConf.CronTimeSyncTransferTxTask > 0 {
		return taskConf.CronTimeSyncTransferTxTask
	}
	return ThreeMinute
}

func (t *IbcSyncTransferTxTask) Run() int {
	return 1
}
