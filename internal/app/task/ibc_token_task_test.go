package task

import "testing"

var tokenTask TokenTask

func TestTokenTaskRun(t *testing.T) {
	tokenTask.Run()
}

func TestUpdateIBCChain(t *testing.T) {
	tokenTask.updateIBCChain()
}

func TestTokenStatisticsTaskRun(t *testing.T) {
	tokenStatisticsTask.Run()
}
