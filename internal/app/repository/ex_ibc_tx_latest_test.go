package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestAggrChainAddress(t *testing.T) {
	res, err := new(ExIbcTxRepo).AggrChainAddress(1672588800, 1672675199, true, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(res)))
}
