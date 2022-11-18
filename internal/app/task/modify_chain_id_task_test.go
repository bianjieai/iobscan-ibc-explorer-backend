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

func Test_ModifyChainIdHandlerThree(t *testing.T) {
	idNameMap, _ := getChainIdNameMap()
	newModifyChainIdHandlerThree(idNameMap).exec()
}

func Test_ModifyChainIdHandlerThreeUpdateColl(t *testing.T) {
	idNameMap, _ := getChainIdNameMap()
	seg := segment{
		StartTime: 1664110315,
		EndTime:   1664110315,
	}
	newModifyChainIdHandlerThree(idNameMap).updateColl(&seg, false)
}
