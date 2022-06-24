package repository

import (
	"encoding/json"
	"testing"
)

func TestIbcRelayerRepo_CountChannelRelayers(t *testing.T) {
	data, err := new(IbcRelayerRepo).CountChannelRelayers()
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}
