package task

import (
	"testing"
)

func Test_FixDenomTraceData(t *testing.T) {
	new(FixDenomTraceDataTask).Run()
}

func Test_FixDenomTraceHistoryData(t *testing.T) {
	new(FixDenomTraceHistoryDataTask).Run()
}

func Test_TraitSegment(t *testing.T) {
	segments := new(fixDenomTraceDataTrait).getSegment(1659910755, 1659949156, 3600)
	t.Log(segments)
}
