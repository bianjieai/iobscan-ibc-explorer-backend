package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestIbcRelayerRepo_FindAllRelayerForCache(t *testing.T) {
	data, err := new(IbcRelayerRepo).FindAllRelayerForCache()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
}
