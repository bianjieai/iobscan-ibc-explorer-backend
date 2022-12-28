package task

import (
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type TokenStatisticsTask struct {
}

var tokenStatisticsTask TokenStatisticsTask

func (t *TokenStatisticsTask) Name() string {
	return "ibc_token_statistics_task"
}

func (t *TokenStatisticsTask) Switch() bool {
	return global.Config.Task.SwitchIbcTokenStatisticsTask
}

func (t *TokenStatisticsTask) Run() int {
	if err := tokenTraceStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s tokenTraceStatisticsRepo.CreateNew err, %v", t.Name(), err)
		return -1
	}

	//if err := tokenStatisticsRepo.CreateNew(); err != nil {
	//	logrus.Errorf("task %s tokenStatisticsRepo.CreateNew err, %v", t.Name(), err)
	//	return -1
	//}

	historySegments, err := getHistorySegment(segmentStepHistory)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}
	logrus.Infof("task %s deal history segment total: %d", t.Name(), len(historySegments))
	if err = t.dealHistory(historySegments); err != nil {
		logrus.Errorf("task %s dealHistory err, %v", t.Name(), err)
		return -1
	}

	segments, err := getSegment(segmentStepLatest)
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}
	logrus.Infof("task %s deal segment total: %d", t.Name(), len(segments))
	if err = t.deal(segments, opInsert); err != nil {
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
			if err = t.saveTraceReceiveData(traceReceiveTxs, v.StartTime, v.EndTime, opInsert); err != nil {
				logrus.Errorf("task %s dealHistory saveTraceReceiveData err, %v", t.Name(), err)
			}
		}
		logrus.Debugf("dealHistory task %s scan ex_ibc_tx finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
	return nil
}

// deal 处理最新的记录，针对ex_ibc_tx_latest
func (t *TokenStatisticsTask) deal(segments []*segment, op int) error {
	for _, v := range segments {
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
			if err = t.saveTraceReceiveData(traceReceiveTxs, v.StartTime, v.EndTime, op); err != nil {
				logrus.Errorf("task %s deal saveTraceReceiveData err, %v", t.Name(), err)
			}
		}
		logrus.Debugf("deal task %s scan ex_ibc_tx_latest finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
	return nil
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

func (t *TokenStatisticsTask) saveTraceReceiveData(dtoList []*dto.CountIBCTokenRecvTxsDTO, segmentStart, segmentEnd int64, op int) error {
	var statistics = make([]*entity.IBCTokenTraceStatistics, 0, len(dtoList))
	for _, v := range dtoList {
		statistics = append(statistics, &entity.IBCTokenTraceStatistics{
			Denom:            v.Denom,
			Chain:            v.Chain,
			ReceiveTxs:       v.Count,
			SegmentStartTime: segmentStart,
			SegmentEndTime:   segmentEnd,
			CreateAt:         time.Now().Unix(),
			UpdateAt:         time.Now().Unix(),
		})
	}

	var err error
	if op == opInsert {
		if err = tokenTraceStatisticsRepo.BatchInsertToNew(statistics); err != nil {
			logrus.Errorf("task %s tokenTraceStatisticsRepo.BatchInsertToNew err, %v", t.Name(), err)
		}
	} else {
		if err = tokenTraceStatisticsRepo.BatchSwap(segmentStart, segmentEnd, statistics); err != nil {
			logrus.Errorf("task %s tokenTraceStatisticsRepo.BatchSwap err, %v", t.Name(), err)
		}
	}

	return err
}
