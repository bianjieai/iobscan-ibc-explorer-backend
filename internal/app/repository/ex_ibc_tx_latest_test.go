package repository

import (
	"encoding/json"
	"testing"
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
	data, err := new(ExIbcTxRepo).GetRelayerInfo(1638979200, 16390655990)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func TestExIbcTxRepo_GetOneRelayerScTxPacketId(t *testing.T) {
	data, err := new(ExIbcTxRepo).GetRelayerInfo(1638979200, 1639065599)
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

func TestExIbcTxRepo_GetHistoryRelayerSuccessPacketTxs(t *testing.T) {
	data, err := new(ExIbcTxRepo).CountHistoryRelayerSuccessPacketTxs(1638979200, 1639065599)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
	data1, err1 := new(ExIbcTxRepo).CountRelayerSuccessPacketTxs(1638979200, 1639065599)
	if err1 != nil {
		t.Fatal(err1.Error())
	}
	ret1, _ := json.Marshal(data1)
	t.Log(string(ret1))
}

func TestExIbcTxRepo_GetRelayerPacketTxs(t *testing.T) {
	data, err := new(ExIbcTxRepo).CountRelayerPacketTxs(1638979200, 1639065599)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
	data1, err1 := new(ExIbcTxRepo).CountHistoryRelayerPacketTxs(1638979200, 1639065599)
	if err1 != nil {
		t.Fatal(err1.Error())
	}
	ret1, _ := json.Marshal(data1)
	t.Log(string(ret1))
}

func TestExIbcTxRepo_GetRelayerPacketAmount(t *testing.T) {
	data, err := new(ExIbcTxRepo).CountRelayerPacketAmount(1638979200, 1639065599)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
	data1, err1 := new(ExIbcTxRepo).CountHistoryRelayerPacketAmount(1638979200, 1639065599)
	if err1 != nil {
		t.Fatal(err1.Error())
	}
	ret1, _ := json.Marshal(data1)
	t.Log(string(ret1))
}
