package task

import "testing"

var channelTask ChannelTask

func TestChannelTaskRun(t *testing.T) {
	channelTask.Run()
}

func TestChannelTask_channelStatistics(t *testing.T) {
	channelTask.channelStatistics()
}

func TestName(t *testing.T) {
	res := channelTask.channelEqual("irishub_qa|channel-33|bigbang|channel-126", "irishub_qa|channel-32|bigbang|channel-125")
	t.Log(res)
}
