package task

import "testing"

var chainOutflowStatisticsTask ChainOutflowStatisticsTask
var chainInflowStatisticsTask ChainInflowStatisticsTask

func Test_chainOutflowStatisticsTaskRunFullStatistics(t *testing.T) {
	chainOutflowStatisticsTask.RunFullStatistics()
}

func Test_chainOutflowStatisticsTaskRun(t *testing.T) {
	chainOutflowStatisticsTask.Run()
}

func Test_chainInflowStatisticsTaskRunFullStatistics(t *testing.T) {
	chainInflowStatisticsTask.RunFullStatistics()
}

func Test_chainInflowStatisticsTaskRun(t *testing.T) {
	chainInflowStatisticsTask.Run()
}
