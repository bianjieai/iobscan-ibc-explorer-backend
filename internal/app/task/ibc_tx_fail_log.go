package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/distributiontask"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IBCTxFailLogTask struct {
}

var _ibcTxFailLogTask distributiontask.CronTask = new(IBCTxFailLogTask)

func (t *IBCTxFailLogTask) Name() string {
	return "ibc_tx_fail_log_task"
}
func (t *IBCTxFailLogTask) Cron() string {
	if taskConf.IBCTxFailLogTask != "" {
		return taskConf.IBCTxFailLogTask
	}

	return "0 0 2 * * ?"
}

func (t *IBCTxFailLogTask) BeforeHook() error {
	return nil
}

func (t *IBCTxFailLogTask) Run() {
	startTime, endTime := utils.YesterdayUnix()
	seg := &segment{
		StartTime: startTime,
		EndTime:   endTime,
	}

	_ = t.deal(seg, true, opUpdate)
	_ = t.deal(seg, false, opUpdate)
}

func (t *IBCTxFailLogTask) RunWithParam(startTime, endTime int64, isTargetHistory bool) int {
	const step = 3600 * 24
	var segs []*segment
	var err error
	if startTime > 0 && endTime > 0 {
		segs = segmentTool(step, startTime, endTime)
	} else {
		segs, err = t.getTxSegments(step, isTargetHistory)
		if err != nil {
			return -1
		}
	}

	for _, v := range segs {
		_ = t.deal(v, isTargetHistory, opInsert)
	}

	return 0
}

func (t *IBCTxFailLogTask) getTxSegments(step int64, isTargetHistory bool) ([]*segment, error) {
	startTime, err := ibcTxRepo.GetMinTxTime(isTargetHistory)
	if err != nil {
		return nil, err
	}

	segs := segmentTool(step, startTime, time.Now().Unix()-86400)
	return segs, nil
}

func (t *IBCTxFailLogTask) deal(seg *segment, isTargetHistory bool, op int) error {
	const limit = 2000
	var skip int64 = 0
	knownTypeFailLogMap := make(map[string]*entity.IBCTxFailLog)
	knownTypeFailLogCountMap := make(map[string]int64)
	var otherTypeFailLogList []*entity.IBCTxFailLog
	nowTime := time.Now().Unix()

	aggrFunc := func(failCode entity.TxFailCode, chain, log string) {
		if chain == "" {
			return
		}

		failLogEntity := &entity.IBCTxFailLog{
			Chain:            chain,
			Log:              log,
			Code:             failCode,
			TxsNumber:        1,
			SegmentStartTime: seg.StartTime,
			SegmentEndTime:   seg.EndTime,
			CreateAt:         nowTime,
			UpdateAt:         nowTime,
		}
		if failCode == entity.TxFailCodeOther {
			otherTypeFailLogList = append(otherTypeFailLogList, failLogEntity)
		} else {
			key := fmt.Sprintf("%s:%s", failLogEntity.Chain, failLogEntity.Code)
			if _, ok := knownTypeFailLogMap[key]; !ok {
				knownTypeFailLogMap[key] = failLogEntity
			}
			knownTypeFailLogCountMap[key] += 1
		}
	}

	logrus.Infof("task %s deal isTargetHistory: %t, segment: %d-%d", t.Name(), isTargetHistory, seg.StartTime, seg.EndTime)
	for {
		txs, err := ibcTxRepo.FindFailLog(seg.StartTime, seg.EndTime, skip, limit, isTargetHistory)
		if err != nil {
			logrus.Errorf("task %s ibcTxRepo.FindFailLog err, %v", t.Name(), err)
			return err
		}

		for _, tx := range txs {
			if tx.Status == entity.IbcTxStatusFailed {
				failCode := t.failType(tx.ScTxInfo.Log)
				aggrFunc(failCode, tx.ScChain, tx.ScTxInfo.Log)
				aggrFunc(failCode, tx.DcChain, tx.ScTxInfo.Log)
			} else { // refund
				var failCode entity.TxFailCode
				var failLog string
				if tx.AckTimeoutTxInfo.Msg.Type == string(entity.TxTypeTimeoutPacket) {
					failCode = entity.TxFailCodeTimeout
				} else {
					failLog = tx.AckTimeoutTxInfo.Msg.AckPacketMsg().Acknowledgement
					failCode = t.failType(failLog)
				}

				aggrFunc(failCode, tx.ScChain, failLog)
				aggrFunc(failCode, tx.DcChain, failLog)
			}
		}

		if len(txs) < limit {
			break
		}

		skip += limit
	}

	var allFailLogList []*entity.IBCTxFailLog
	if len(knownTypeFailLogMap) > 0 {
		allFailLogList = make([]*entity.IBCTxFailLog, 0, len(knownTypeFailLogMap))
		for k, v := range knownTypeFailLogMap {
			v.TxsNumber = knownTypeFailLogCountMap[k]
			allFailLogList = append(allFailLogList, v)
		}
	}
	if len(otherTypeFailLogList) > 0 {
		allFailLogList = append(allFailLogList, otherTypeFailLogList...)
	}

	if op == opInsert {
		if err := ibcTxFailLogRepo.BatchInsert(allFailLogList); err != nil {
			logrus.Errorf("task %s BatchInsert err, %v", t.Name(), err)
		}
	} else {
		if err := ibcTxFailLogRepo.BatchSwap(seg.StartTime, seg.EndTime, allFailLogList); err != nil {
			logrus.Errorf("task %s BatchSwap err, %v", t.Name(), err)
		}
	}

	logrus.Infof("task %s deal end isTargetHistory: %t, segment: %d-%d", t.Name(), isTargetHistory, seg.StartTime, seg.EndTime)
	return nil
}

func (t *IBCTxFailLogTask) failType(log string) entity.TxFailCode {
	if strings.Contains(log, "packet timeout") {
		return entity.TxFailCodeTimeout
	} else if strings.Contains(log, "out of gas") {
		return entity.TxFailCodeOutOfGas
	} else if strings.Contains(log, "insufficient funds") {
		return entity.TxFailCodeInsufficientFunds
	} else if strings.Contains(log, "client is not active") {
		return entity.TxFailCodeClientNotActive
	} else if strings.Contains(log, "cannot parse packet fowrading information") {
		return entity.TxFailCodeParsePacketFowradingInfoErr
	} else {
		return entity.TxFailCodeOther
	}
}
