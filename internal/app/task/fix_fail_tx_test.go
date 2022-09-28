package task

import "testing"

func TestFixFailTxTask_Run(t *testing.T) {
	new(FixFailTxTask).Run()
}

func TestFixFailTxTask_FixRecvPacketTxs(t *testing.T) {
	new(FixFailTxTask).fixFailTxs(ibcTxTargetHistory, []*segment{
		{StartTime: 1652457600, EndTime: 1652500800},
		{StartTime: 1653926400, EndTime: 1653969600},
		{StartTime: 1654444800, EndTime: 1654488000},
	})
}
