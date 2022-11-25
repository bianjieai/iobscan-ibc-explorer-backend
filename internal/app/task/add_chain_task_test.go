package task

import (
	"testing"
)

func Test_AddChainTask(t *testing.T) {
	new(AddChainTask).Run()
}

func Test_UpdateIbcTx(t *testing.T) {
	chainMap, _ := getAllChainMap()
	chain := "bigbangname"
	new(AddChainTask).updateIbcTx(chain, chainMap[chain], chainMap)
}
