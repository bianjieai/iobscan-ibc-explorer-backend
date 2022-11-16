package task

import "testing"

func Test_ModifyChainIdHandlerOne(t *testing.T) {
	idNameMap, _ := getChainIdNameMap()
	newModifyChainIdHandlerOne(idNameMap).exec("")
}

func Test_ModifyChainIdHandlerOneUpdateColl(t *testing.T) {
	idNameMap, _ := getChainIdNameMap()
	newModifyChainIdHandlerOne(idNameMap).updateColl("ibc_token_trace_statistics")
}
