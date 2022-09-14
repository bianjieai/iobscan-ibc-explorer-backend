package task

import "testing"

func TestName(t *testing.T) {
	new(FixAcknowledgeTxTask).Run()
}
