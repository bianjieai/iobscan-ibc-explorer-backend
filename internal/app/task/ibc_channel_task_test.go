package task

import "testing"

var channelTask ChannelTask

func TestChannelTaskRun(t *testing.T) {
	channelTask.Run()
}

func TestChannelTask_channelStatistics(t *testing.T) {
	channelTask.channelStatistics()
}
