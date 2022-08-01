package service

import "testing"

func TestStatisticInfoService_IbcTxStatistic(t *testing.T) {
	resp, err := new(StatisticInfoService).IbcTxStatistic()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(resp)
}
