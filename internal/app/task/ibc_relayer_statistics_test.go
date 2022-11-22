package task

import "testing"

func TestRelayerStatisticsTask_Run(t *testing.T) {
	new(RelayerStatisticsTask).Run()
}

func Test_RelayerStatisticRunIncrement(t *testing.T) {
	seg := segment{
		StartTime: 1636761600,
		EndTime:   1636847999,
	}
	_ = relayerStatisticsTask.RunIncrement(&seg)
}

func Test_RelayerStatisticsRunWithParam(t *testing.T) {
	relayerStatisticsTask.RunWithParam("cosmoshub_4", 1640995200, 1641081599)
}

func TestFixRelayerStatisticsTask_Run(t *testing.T) {
	new(FixRelayerStatisticsTask).Run()
}
