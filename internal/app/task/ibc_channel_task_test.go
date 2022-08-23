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
	res, err := getSegment()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.MustMarshalJsonToStr(res))

	res, err = getHistorySegment()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.MustMarshalJsonToStr(res))
}

func Test_ChannelStatistics(t *testing.T) {
	channelStatisticsTask.Run()
}
