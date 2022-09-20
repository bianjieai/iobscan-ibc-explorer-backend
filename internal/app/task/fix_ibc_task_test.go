package task

import "testing"

func Test_FixIbcTask(t *testing.T) {
	new(FixIbxTxTask).Run()
}

func Test_FixSegment(t *testing.T) {
	s := segment{
		StartTime: 1632814442,
		EndTime:   1632814442,
	}
	new(FixIbxTxTask).fixSegment(&s, false)
}
