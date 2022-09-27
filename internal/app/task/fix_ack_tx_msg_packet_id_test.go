package task

import (
	"testing"
)

func Test_fixAckTxTask(t *testing.T) {
	NewfixAckTxTask("bigbang", "", 4016630, 4977706).Run()
}

func TestFixAckTxPacketIdTask_Run(t *testing.T) {
	new(FixAckTxPacketIdTask).RunWithParam("qa_iris_snapshot", "")
}
