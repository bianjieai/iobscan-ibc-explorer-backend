package task

import "testing"

func TestIBCChainFeeStatisticTask_Run(t *testing.T) {
	new(IBCChainFeeStatisticTask).Run()
}

func TestIBCChainFeeStatisticTask_RunAllChain(t *testing.T) {
	new(IBCChainFeeStatisticTask).RunAllChain()
}

func TestIBCChainFeeStatisticTask_RunIncrement(t *testing.T) {
	new(IBCChainFeeStatisticTask).RunIncrement(&segment{StartTime: 1673222400, EndTime: 1673308799})
}

func TestIBCChainFeeStatisticTask_RunWithParam(t *testing.T) {
	new(IBCChainFeeStatisticTask).RunWithParam("mantle_1", 1650755006, 1650816760)
}
