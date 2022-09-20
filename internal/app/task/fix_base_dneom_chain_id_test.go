package task

import "testing"

var _fixBaseDenomChainIdTask = new(FixBaseDenomChainIdTask)

func Test_FixBaseDenomChainId(t *testing.T) {
	_fixBaseDenomChainIdTask.Run()
}

func Test_FixBaseDenomChainIdHandleSegment(t *testing.T) {
	if err := _fixBaseDenomChainIdTask.init(); err != nil {
		t.Fatal(err)
	}
	_fixBaseDenomChainIdTask.handleSegment(1617883200, 1617926399, "uphoton", true)
}
