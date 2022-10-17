package task

import "testing"

func TestCheckAndUpdateTraceSourceNode(t *testing.T) {
	CheckAndUpdateTraceSourceNode("crescent_1")
}

func TestIbcNodeLcdCronTask_Run(t *testing.T) {
	new(IbcNodeLcdCronTask).Run()
}
