package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"sync"
)

type FixFailRecvPacketTask struct {
}

var _ OneOffTask = new(FixFailRecvPacketTask)

func (t *FixFailRecvPacketTask) Name() string {
	return "fix_fail_recv_packet_task"
}

func (t *FixFailRecvPacketTask) Switch() bool {
	return global.Config.Task.SwitchFixFailRecvPacketTask
}

func (t *FixFailRecvPacketTask) Run() int {
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

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := t.fixFailRecvPacketTxs(ibcTxTargetLatest, segments)
		logrus.Infof("task %s fix latest end, %v", t.Name(), err)
	}()

	go func() {
		defer wg.Done()
		err := t.fixFailRecvPacketTxs(ibcTxTargetHistory, historySegments)
		logrus.Infof("task %s fix history end, %v", t.Name(), err)
	}()

	wg.Wait()
	return 1
}

func (t *FixFailRecvPacketTask) fixFailRecvPacketTxs(target string, segments []*segment) error {
	const limit int64 = 1000
	isTargetHistory := false
	if target == ibcTxTargetHistory {
		isTargetHistory = true
	}
	doHandleSegments(t.Name(), 3, segments, isTargetHistory, func(seg *segment, isTargetHistory bool) {
		var skip int64 = 0
		for {
			txs, err := ibcTxRepo.FindRecvPacketTxsEmptyTxs(seg.StartTime, seg.EndTime, skip, limit, isTargetHistory)
			if err != nil {
				logrus.Errorf("task %s FindRecvPacketTxsEmptyTxs %s %d-%d err, %v", t.Name(), target, seg.StartTime, seg.EndTime, err)
				return
			}

			for _, val := range txs {
				err := SaveRecvPacketTx(val, isTargetHistory)
				if err != nil && err != qmgo.ErrNoSuchDocuments {
					logrus.Errorf("task %s SaveRecvPacketTx %s err, chain_id: %s, packet_id: %s, %v", t.Name(), target, val.ScChain, val.ScTxInfo.Msg.CommonMsg().PacketId, err)
					return
				}
			}

			if int64(len(txs)) < limit {
				break
			}
			skip += limit
		}
	})

	return nil
}
