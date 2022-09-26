package task

import (
	"fmt"
	"testing"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
)

var channelTask ChannelTask

func TestChannelTaskRun(t *testing.T) {
	channelTask.Run()
}

func Test_getSegment(t *testing.T) {
	res, err := getSegment(segmentStepLatest)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.MustMarshalJsonToStr(res))

	res, err = getHistorySegment(segmentStepHistory)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.MustMarshalJsonToStr(res))
}

func Test_ChannelStatistics(t *testing.T) {
	channelStatisticsTask.Run()
}
