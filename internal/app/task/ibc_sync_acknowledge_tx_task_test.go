package task

import "testing"

func Test_IbcSyncAcknowledgeTxTask(t *testing.T) {
	new(IbcSyncAcknowledgeTxTask).Run()
}
