package repository

import (
	"encoding/json"
	"testing"
)

func TestRelayerStatisticsRepo_CountRelayerTotalValue(t *testing.T) {
	data, err := new(RelayerStatisticsRepo).CountRelayerBaseDenomAmt()
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}
