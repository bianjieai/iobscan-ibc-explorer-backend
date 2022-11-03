package repository

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExIbcTxRepo_FindAll(t *testing.T) {
	data, err := new(ExIbcTxRepo).FindAll(0, 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func TestExIbcTxRepo_FindAllHistory(t *testing.T) {
	data, err := new(ExIbcTxRepo).FindAllHistory(0, 20)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func TestExIbcTxRepo_GetRelayerInfo(t *testing.T) {
	now := time.Now().Unix()
	data, err := new(ExIbcTxRepo).GetRelayerInfo(now-86400, now)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func TestExIbcTxRepo_GetOneRelayerScTxPacketId(t *testing.T) {
	now := time.Now().Unix()
	data, err := new(ExIbcTxRepo).GetRelayerInfo(now-86400, now)
	if err != nil {
		t.Fatal(err.Error())
	}

	data1, err1 := new(ExIbcTxRepo).GetOneRelayerScTxPacketId(data[0])
	if err1 != nil {
		t.Fatal(err1.Error())
	}
	ret1, _ := json.Marshal(data1)
	t.Log(string(ret1))
}
