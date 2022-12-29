package task

import (
	"testing"
)

var channelTask ChannelTask

func TestChannelTaskRun(t *testing.T) {
	channelTask.Run()
}

func Test_ChannelStatistics(t *testing.T) {
	channelStatisticsTask.Run()
}
