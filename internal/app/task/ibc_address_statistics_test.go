package task

import "testing"

func TestIBCAddressStatisticTask_Run(t *testing.T) {
	new(IBCAddressStatisticTask).Run()
}

func TestIBCAddressStatisticTask_RunAllChain(t *testing.T) {
	new(IBCAddressStatisticTask).RunAllChain()
}

func TestIBCAddressStatisticTask_RunIncrement(t *testing.T) {
	new(IBCAddressStatisticTask).RunIncrement(&segment{StartTime: 1673222400, EndTime: 1673308799})
}

func TestIBCAddressStatisticTask_RunWithParam(t *testing.T) {
	new(IBCAddressStatisticTask).RunWithParam("", 1675555200, 1675641599)
}
