package task

import (
	"fmt"
	"math"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type TokenStatisticsTask struct {
	segmentMinTime    int64
	segmentRecvTxsMap map[string][]*dto.CountIBCTokenRecvTxsDTO
}

var tokenStatisticsTask TokenStatisticsTask

func (t *TokenStatisticsTask) Name() string {
	return "ibc_token_statistics_task"
}

func (t *TokenStatisticsTask) Switch() bool {
	return global.Config.Task.SwitchIbcTokenStatisticsTask
}

func TokenIncrementStatistics(segments []*segment) error {
	return tokenStatisticsTask.deal(segments, false)
}

func (t *TokenStatisticsTask) Run() int {
	t.segmentMinTime = math.MaxInt64
	t.segmentRecvTxsMap = make(map[string][]*dto.CountIBCTokenRecvTxsDTO)

	if err := tokenTraceStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s tokenTraceStatisticsRepo.CreateNew err, %v", t.Name(), err)
		return -1
	}

	//if err := tokenStatisticsRepo.CreateNew(); err != nil {
	//	logrus.Errorf("task %s tokenStatisticsRepo.CreateNew err, %v", t.Name(), err)
	//	return -1
	//}

	segments, err := getTxTimeSegment(false, segmentStepLatest)
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}

	historySegments, err := getTxTimeSegment(true, segmentStepHistory)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}

	t.segmentMinTime = segments[0].StartTime
	// 优先处理历史分段
	logrus.Infof("task %s deal history segment total: %d", t.Name(), len(historySegments))
	if err = t.dealHistory(historySegments); err != nil {
		logrus.Errorf("task %s dealHistory err, %v", t.Name(), err)
		return -1
	}

	logrus.Infof("task %s deal segment total: %d", t.Name(), len(segments))
	if err = t.deal(segments, true); err != nil {
		logrus.Errorf("task %s deal err, %v", t.Name(), err)
		return -1
	}

	if err = tokenTraceStatisticsRepo.SwitchColl(); err != nil {
		logrus.Errorf("task %s tokenTraceStatisticsRepo.SwitchColl err, %v", t.Name(), err)
		return -1
	}

	//if err = tokenStatisticsRepo.SwitchColl(); err != nil {
	//	logrus.Errorf("task %s tokenStatisticsRepo.SwitchColl err, %v", t.Name(), err)
	//	return -1
	//}

	return 1
}

// dealHistory 处理历史记录，针对ex_ibc_tx
func (t *TokenStatisticsTask) dealHistory(segments []*segment) error {
	for _, v := range segments {
		logrus.Infof("task %s dealHistory segment [%d, %d]", t.Name(), v.StartTime, v.EndTime)
		//transferTxs, err := ibcTxRepo.CountBaseDenomHistoryTransferTxs(v.StartTime, v.EndTime)
		//if err != nil {
		//	logrus.Errorf("task %s CountBaseDenomHistoryTransferTxs err, %v", t.Name(), err)
		//	return err
		//}
		//
		//if len(transferTxs) > 0 {
		//	if err = t.saveTokenTransferData(transferTxs, v.StartTime, v.EndTime, opInsert); err != nil {
		//		logrus.Errorf("task %s dealHistory saveTokenTransferData err, %v", t.Name(), err)
		//	}
		//}

		traceReceiveTxs, err := ibcTxRepo.CountIBCTokenHistoryRecvTxs(v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s CountIBCTokenHistoryRecvTxs err, %v", t.Name(), err)
			return err
		}

		if len(traceReceiveTxs) > 0 {
			// 将新老表重叠的分段数据记录到map
			if v.StartTime >= t.segmentMinTime {
				t.segmentRecvTxsMap[fmt.Sprintf("%d-%d", v.StartTime, v.EndTime)] = traceReceiveTxs
			}
			if err = t.saveTraceReceiveData(traceReceiveTxs, v, true, true); err != nil {
				logrus.Errorf("task %s dealHistory saveTraceReceiveData err, %v", t.Name(), err)
			}
		}
		logrus.Debugf("dealHistory task %s scan ex_ibc_tx finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
	return nil
}

// deal 处理最新的记录，针对ex_ibc_tx_latest
func (t *TokenStatisticsTask) deal(segments []*segment, fullStatistics bool) error {
	for _, v := range segments {
		logrus.Infof("task %s deal segment [%d, %d]", t.Name(), v.StartTime, v.EndTime)
		//transferTxs, err := ibcTxRepo.CountBaseDenomTransferTxs(v.StartTime, v.EndTime)
		//if err != nil {
		//	logrus.Errorf("task %s CountBaseDenomTransferTxs err, %v", t.Name(), err)
		//	return err
		//}
		//
		//if len(transferTxs) > 0 {
		//	if err = t.saveTokenTransferData(transferTxs, v.StartTime, v.EndTime, op); err != nil {
		//		logrus.Errorf("task %s deal saveTokenTransferData err, %v", t.Name(), err)
		//	}
		//}

		traceReceiveTxs, err := ibcTxRepo.CountIBCTokenRecvTxs(v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s CountIBCTokenRecvTxs err, %v", t.Name(), err)
			return err
		}

		if len(traceReceiveTxs) > 0 {
			if fullStatistics {
				traceReceiveTxs = t.integrationStatisticsData(traceReceiveTxs, v)
			}
			if err = t.saveTraceReceiveData(traceReceiveTxs, v, false, fullStatistics); err != nil {
				logrus.Errorf("task %s deal saveTraceReceiveData err, %v", t.Name(), err)
			}
		}
		logrus.Debugf("deal task %s scan ex_ibc_tx_latest finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
	return nil
}

func (t *TokenStatisticsTask) integrationStatisticsData(aggrRes []*dto.CountIBCTokenRecvTxsDTO, seg *segment) []*dto.CountIBCTokenRecvTxsDTO {
	// 新表中的段与历史表重和，需要整合数据
	hirtoryAggrRes, ok := t.segmentRecvTxsMap[fmt.Sprintf("%d-%d", seg.StartTime, seg.EndTime)]
	if !ok {
		return aggrRes
	}

	integrationDataMap := make(map[string]*dto.CountIBCTokenRecvTxsDTO, len(aggrRes))
	for _, v := range aggrRes {
		key := fmt.Sprintf("%s%s", v.Chain, v.Denom)
		if data, ok := integrationDataMap[key]; !ok {
			integrationDataMap[key] = v
		} else {
			data.Count += v.Count
		}
	}

	for _, v := range hirtoryAggrRes {
		key := fmt.Sprintf("%s%s", v.Chain, v.Denom)
		if data, ok := integrationDataMap[key]; !ok {
			integrationDataMap[key] = v
		} else {
			data.Count += v.Count
		}
	}

	integrationDataList := make([]*dto.CountIBCTokenRecvTxsDTO, 0, len(integrationDataMap))
	for _, v := range integrationDataMap {
		integrationDataList = append(integrationDataList, v)
	}
	return integrationDataList
}

//func (t *TokenStatisticsTask) saveTokenTransferData(dtoList []*dto.CountBaseDenomTxsDTO, segmentStart, segmentEnd int64, op int) error {
//	var statistics = make([]*entity.IBCTokenStatistics, 0, len(dtoList))
//	for _, v := range dtoList {
//		statistics = append(statistics, &entity.IBCTokenStatistics{
//			BaseDenom:        v.BaseDenom,
//			BaseDenomChain:   v.BaseDenomChain,
//			TransferTxs:      v.Count,
//			SegmentStartTime: segmentStart,
//			SegmentEndTime:   segmentEnd,
//			CreateAt:         time.Now().Unix(),
//			UpdateAt:         time.Now().Unix(),
//		})
//	}
//
//	var err error
//	if op == opInsert {
//		if err = tokenStatisticsRepo.BatchInsertToNew(statistics); err != nil {
//			logrus.Errorf("task %s tokenStatisticsRepo.BatchInsertToNew err, %v", t.Name(), err)
//		}
//	} else {
//		if err = tokenStatisticsRepo.BatchSwap(segmentStart, segmentEnd, statistics); err != nil {
//			logrus.Errorf("task %s tokenStatisticsRepo.BatchSwap err, %v", t.Name(), err)
//		}
//	}
//
//	return err
//}

func (t *TokenStatisticsTask) saveTraceReceiveData(dtoList []*dto.CountIBCTokenRecvTxsDTO, seg *segment, targetHistory, fullStatistics bool) error {
	var statistics = make([]*entity.IBCTokenTraceStatistics, 0, len(dtoList))
	timeUnix := time.Now().Unix()
	for _, v := range dtoList {
		statistics = append(statistics, &entity.IBCTokenTraceStatistics{
			Denom:            v.Denom,
			Chain:            v.Chain,
			ReceiveTxs:       v.Count,
			SegmentStartTime: seg.StartTime,
			SegmentEndTime:   seg.EndTime,
			CreateAt:         timeUnix,
			UpdateAt:         timeUnix,
		})
	}

	var err error
	if fullStatistics {
		if targetHistory {
			err = tokenTraceStatisticsRepo.BatchInsertToNew(statistics)
		} else {
			err = tokenTraceStatisticsRepo.BatchSwapNew(seg.StartTime, seg.EndTime, statistics)
		}
	} else {
		err = tokenTraceStatisticsRepo.BatchSwap(seg.StartTime, seg.EndTime, statistics)
	}

	return err
}
