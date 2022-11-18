package task

import (
	"encoding/json"
	"testing"
)

var (
	_fixChainCfg          *modifyChainConfigTask
	_fixChannelIdTask     *fixChannelIdTask
	_fixIbcRelayerTask    *fixIbcRelayerTask
	_fixIbcTaskRecordTask *fixIbcTaskRecordTask
)

func TestNewModifyChainIdContent(t *testing.T) {
	client := NewModifyChainIdContent("chain_config_copy")
	if client != nil {
		client.Run()
	}
}

func TestGetIbcRelayerData(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixIbcRelayerTask = NewModifyIbcRelayer(mapData)
	datas, err := _fixIbcRelayerTask.GetIbcRelayerData(0, 100)
	if err != nil {
		t.Fatal(err.Error())
	}
	bytesdata, _ := json.Marshal(datas)
	t.Log(string(bytesdata))

	_fixIbcRelayerTask.UpdateIbcRelayerData(*datas[0])
}

func TestFixIbcRelayerTaskRun(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixIbcRelayerTask = NewModifyIbcRelayer(mapData)
	_fixIbcRelayerTask.Run()
}

func TestFixChannelIdTaskGetIbcChannelData(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixChannelIdTask = NewModifyChannalIdTask(mapData)
	datas, err := _fixChannelIdTask.GetIbcChannelData(0, 100)
	if err != nil {
		t.Fatal(err.Error())
	}
	bytesdata, _ := json.Marshal(datas)
	t.Log(string(bytesdata))

	_fixChannelIdTask.UpdateIbcChannel(*datas[0])
}

func TestGetIbcChannelStatisticData(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixChannelIdTask = NewModifyChannalIdTask(mapData)
	datas, err := _fixChannelIdTask.GetIbcChannelStatisticData(0, 11)
	if err != nil {
		t.Fatal(err.Error())
	}
	bytesdata, _ := json.Marshal(datas)
	t.Log(string(bytesdata))

	for _, val := range datas {
		_fixChannelIdTask.UpdateIbcChannelStatistic(*val)
	}
}

func TestFixChannelIdTaskRun(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixChannelIdTask = NewModifyChannalIdTask(mapData)
	_fixChannelIdTask.Run()
}

func TestGetAllChainConigs(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixChainCfg = NewModifyChainConfig(mapData)
	datas, err := _fixChainCfg.GetAllChainConigs()
	if err != nil {
		t.Fatal(err.Error())
	}
	bytesdata, _ := json.Marshal(datas)
	t.Log(string(bytesdata))
	//_fixChainCfg.UpdateChainConfig(*datas[0])

}

func TestModifyChainConfigRun(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixChainCfg = NewModifyChainConfig(mapData)
	_fixChainCfg.Run()
}

func TestGetIbcTaskRecordData(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixIbcTaskRecordTask = NewModifyIbcTaskRecord(mapData)
	datas, err := _fixIbcTaskRecordTask.GetIbcTaskRecordData()
	if err != nil {
		t.Fatal(err.Error())
	}
	bytesdata, _ := json.Marshal(datas)
	t.Log(string(bytesdata))

	for _, val := range datas {
		_fixIbcTaskRecordTask.UpdateIbcTaskRecord(*val)
	}

}

func TestFixIbcTaskRecordTaskRun(t *testing.T) {
	mapData, err := _initChainVerCfgMap()
	if err != nil {
		t.Fatal(err.Error())
	}
	_fixIbcTaskRecordTask = NewModifyIbcTaskRecord(mapData)
	_fixIbcTaskRecordTask.Run()
}
