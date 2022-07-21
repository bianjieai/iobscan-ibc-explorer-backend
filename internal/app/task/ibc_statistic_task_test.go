package task

import "testing"

var (
	_statistic IbcStatisticCronTask
)

func TestIbcStatisticCronTask_updateChannel24h(t *testing.T) {
	if err := _statistic.updateChannel24h(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestIbcStatisticCronTask_updateChannelInfo(t *testing.T) {
	if err := _statistic.updateChannelInfo(); err != nil {
		t.Fatal(err.Error())
	}
}
