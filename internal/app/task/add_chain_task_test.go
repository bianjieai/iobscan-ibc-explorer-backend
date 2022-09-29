package task

import (
	"testing"
)

func Test_AddChainTask(t *testing.T) {
	new(AddChainTask).Run()
}

func Test_UpdateIbcTx(t *testing.T) {
	chainMap, _ := getAllChainMap()
	chainId := "bigbang"
	new(AddChainTask).updateIbcTx(chainId, chainMap[chainId], chainMap)
}
