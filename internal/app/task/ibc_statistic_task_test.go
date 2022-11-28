package task

import "testing"

var (
	_statistic IbcStatisticCronTask
)

func TestIbcStatisticCronTask_updateChannel24h(t *testing.T) {
	if err := _statistic.updateChannelAndChains24h(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestIbcStatisticCronTask_updateChannelInfo(t *testing.T) {
	if err := _statistic.updateChannelInfo(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestIbcStatisticCronTask_updateChains(t *testing.T) {
	if err := _statistic.updateChains(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestIbcStatisticCronTask_updateTxsIncre(t *testing.T) {
	if err := _statistic.updateTxsIncre(); err != nil {
		t.Fatal(err.Error())
	}
}
