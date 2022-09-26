package task

import (
	"fmt"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
)

var relayerDataTask RelayerDataTask
var _ OneOffTask = new(RelayerDataTask)

type RelayerDataTask struct {
	distRelayerMap map[string]bool
}

func (t *RelayerDataTask) Name() string {
	return "ibc_relayer_data_task"
}

func (t *RelayerDataTask) Switch() bool {
	return global.Config.Task.SwitchOnlyInitRelayerData
}

func (t *RelayerDataTask) Run() int {
	startTime := time.Now().Unix()
	historySegments, err := getHistorySegment(segmentStepHistory)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}
	//insert relayer data
	t.handleNewRelayerOnce(historySegments, true)

	segments, err := getSegment(segmentStepLatest)
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}
	//insert relayer data
	t.handleNewRelayerOnce(segments, false)
	logrus.Infof("task %s finish deal, time use %d(s)", t.Name(), time.Now().Unix()-startTime)
	return 1
}

func (t *RelayerDataTask) initdistRelayerMap() {
	t.distRelayerMap = make(map[string]bool, 20)
	skip := int64(0)
	limit := int64(1000)
	for {
		dbRelayers, err := relayerRepo.FindAll(skip, limit)
		if err != nil {
			logrus.Error("find relayer by page fail, ", err.Error())
			return
		}

		for _, val := range dbRelayers {
			key := fmt.Sprintf("%s:%s:%s", val.ChainA, val.ChainAAddress, val.ChannelA)
			key1 := fmt.Sprintf("%s:%s:%s", val.ChainB, val.ChainBAddress, val.ChannelB)
			t.distRelayerMap[key] = true
			t.distRelayerMap[key1] = true
		}
		if len(dbRelayers) < int(limit) {
			break
		}
		skip += limit
	}

	return
}

func (t *RelayerDataTask) handleNewRelayerOnce(segments []*segment, historyData bool) {
	t.initdistRelayerMap()
	for _, v := range segments {
		var relayersData []entity.IBCRelayer
		if historyData {
			relayersData = t.handleIbcTxHistory(v.StartTime, v.EndTime)
		} else {
			relayersData = t.handleIbcTxLatest(v.StartTime, v.EndTime)
		}
		if len(relayersData) > 0 {
			relayersData = distinctRelayer(relayersData, t.distRelayerMap)
			relayersData = filterDbExist(relayersData, t.distRelayerMap)
			if len(relayersData) == 0 {
				continue
			}
			if err := relayerRepo.Insert(relayersData); err != nil && !qmgo.IsDup(err) {
				logrus.Error("insert relayer data fail, ", err.Error())
			}
		}
		logrus.Debugf("task %s find relayer finish segment [%v:%v], relayers:%v", t.Name(), v.StartTime, v.EndTime, len(relayersData))
	}
}

func (t *RelayerDataTask) handleIbcTxLatest(startTime, endTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, val := range relayerDtos {
		relayers = append(relayers, createRelayerData(val))
	}
	return relayers
}

func (t *RelayerDataTask) handleIbcTxHistory(startTime, endTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetHistoryRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, val := range relayerDtos {
		relayers = append(relayers, createRelayerData(val))
	}
	return relayers
}
