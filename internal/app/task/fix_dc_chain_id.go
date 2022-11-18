package task

import (
	"sync"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type FixDcChainIdTask struct {
	chainMap map[string]*entity.ChainConfig
}

var _ OneOffTask = new(FixDcChainIdTask)

func (t *FixDcChainIdTask) Name() string {
	return "fix_dc_chain_id_task"
}

func (t *FixDcChainIdTask) Switch() bool {
	return global.Config.Task.SwitchFixDcChainIdTask
}

func (t *FixDcChainIdTask) Run() int {
	segments, err := getSegment(segmentStepLatest)
	if err != nil {
		logrus.Errorf("task %s getSegment error, %v", t.Name(), err)
		return -1
	}

	historySegments, err := getHistorySegment(segmentStepHistory)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment error, %v", t.Name(), err)
		return -1
	}

	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", t.Name(), err)
		return -1
	}
	t.chainMap = chainMap

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		t.fixDcChainId(ibcTxTargetLatest, segments)
		logrus.Infof("task %s fix latest end", t.Name())
	}()

	go func() {
		defer wg.Done()
		t.fixDcChainId(ibcTxTargetHistory, historySegments)
		logrus.Infof("task %s fix history end", t.Name())
	}()

	wg.Wait()
	return 1
}

func (t *FixDcChainIdTask) fixDcChainId(target string, segments []*segment) {
	const limit int64 = 1000
	isTargetHistory := false
	if target == ibcTxTargetHistory {
		isTargetHistory = true
	}

	for _, v := range segments {
		logrus.Infof("task %s fix %s %d-%d", t.Name(), target, v.StartTime, v.EndTime)
		var skip int64 = 0
		var toBeFixedTxs []*fixDcChainIdTx
		for {
			txs, err := ibcTxRepo.FindDcChainIdEmptyTxs(v.StartTime, v.EndTime, skip, limit, isTargetHistory)
			if err != nil {
				logrus.Errorf("task %s FindDcChainIdEmptyTxs %s %d-%d err, %v", t.Name(), target, v.StartTime, v.EndTime, err)
				break
			}

			for _, tx := range txs {
				dcChainId, _, dcChannel := matchDcInfo(tx.ScChain, constant.PortTransfer, tx.ScChannel, t.chainMap)
				toBeFixedTxs = append(toBeFixedTxs, &fixDcChainIdTx{
					RecordId:  tx.RecordId,
					DcChainId: dcChainId,
					DcChannel: dcChannel,
					Status:    tx.Status,
				})
			}

			if int64(len(txs)) < limit {
				break
			}
			skip += limit
		}

		for _, tx := range toBeFixedTxs {
			if err := ibcTxRepo.FixDcChainId(tx.RecordId, tx.DcChainId, tx.DcChannel, tx.Status, isTargetHistory); err != nil {
				logrus.Errorf("task %s FixDcChainId(%s) %s err, dcChainId: %s, dcChannel: %s, %v", t.Name(), tx.RecordId, target, tx.DcChainId, tx.DcChannel, err)
			}
		}
	}
}

type fixDcChainIdTx struct {
	RecordId  string
	DcChainId string
	DcChannel string
	Status    entity.IbcTxStatus
}
