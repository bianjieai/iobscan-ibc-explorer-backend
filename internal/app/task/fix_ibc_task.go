package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type FixIbxTxTask struct {
	chainMap map[string]*entity.ChainConfig
	domain   string // all, partly
}

var _ OneOffTask = new(FixIbxTxTask)

const (
	domainAll    = "all"
	domainPartly = "partly"
)

func (t *FixIbxTxTask) Name() string {
	return "fix_ibc_tx_task"
}

func (t *FixIbxTxTask) Switch() bool {
	return false
}

func (t *FixIbxTxTask) Run() int {
	t.domain = domainAll
	return t.handle()
}

func (t *FixIbxTxTask) RunWithParam(domain string) int {
	t.domain = domain
	return t.handle()
}

func (t *FixIbxTxTask) handle() int {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}
	t.chainMap = chainMap

	historySegment, err := getHistorySegment(4 * 3600)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}
	t.workerGroup(true, historySegment)

	segments, err := getSegment(8 * 3600)
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}
	t.workerGroup(false, segments)

	return 1
}

func (t *FixIbxTxTask) workerGroup(isTargetHistory bool, segments []*segment) {
	st := time.Now().Unix()
	logrus.Infof("task %s worker group start, target hirtoty: %t", t.Name(), isTargetHistory)
	var wg sync.WaitGroup
	wg.Add(fixIbxTxWorkerNum)
	for i := 0; i < fixIbxTxWorkerNum; i++ {
		num := i
		go func() {
			defer wg.Done()
			t.workerExec(isTargetHistory, segments, num)
		}()
	}
	wg.Wait()
	logrus.Infof("task %s worker group end, target hirtoty: %t, time use: %d(s)", t.Name(), isTargetHistory, time.Now().Unix()-st)
}

func (t *FixIbxTxTask) workerExec(isTargetHistory bool, segments []*segment, workerNum int) {
	for i, v := range segments {
		if i%fixIbxTxWorkerNum != workerNum {
			continue
		}
		logrus.Infof("task %s worker %d fix %d-%d, target history: %t", t.Name(), workerNum, v.StartTime, v.EndTime, isTargetHistory)
		t.fixSegment(v, isTargetHistory)
	}
}

func (t *FixIbxTxTask) fixSegment(seg *segment, isTargetHistory bool) {
	const limit int64 = 2000
	var skip int64 = 0

	for {
		var txs []*entity.ExIbcTx
		var err error
		if t.domain == domainAll {
			txs, err = ibcTxRepo.FindByCreateAt(seg.StartTime, seg.EndTime, skip, limit, isTargetHistory)
		} else {
			txs, err = ibcTxRepo.FindEmptyDcConnTxs(seg.StartTime, seg.EndTime, skip, limit, isTargetHistory)
		}
		if err != nil {
			logrus.Errorf("task %s fixSegment %t %d-%d", t.Name(), isTargetHistory, seg.StartTime, seg.EndTime)
			return
		}

		t.fixTxs(txs, isTargetHistory)

		if int64(len(txs)) < limit {
			break
		}
		skip += limit
	}
}

func (t *FixIbxTxTask) fixTxs(ibcTxs []*entity.ExIbcTx, isTargetHistory bool) {
	if len(ibcTxs) == 0 {
		return
	}

	// 找出需要用到的sync tx
	chainHashMap := make(map[string][]string)
	for _, v := range ibcTxs {
		if v.DcTxInfo != nil {
			chainHashMap[v.DcChainId] = append(chainHashMap[v.DcChainId], v.DcTxInfo.Hash)
		}
		if v.RefundedTxInfo != nil {
			chainHashMap[v.ScChainId] = append(chainHashMap[v.ScChainId], v.RefundedTxInfo.Hash)
		}
		if v.ScTxInfo != nil {
			chainHashMap[v.ScChainId] = append(chainHashMap[v.ScChainId], v.ScTxInfo.Hash)
		}
	}

	chainHashTxMap := make(map[string]*entity.Tx)
	for k, v := range chainHashMap {
		txs, err := txRepo.GetTxByHashes(k, v)
		if err != nil {
			logrus.Errorf("task %s GetTxByHashes err, %v", t.Name(), err)
			continue
		}

		for _, tx := range txs {
			chainHashTxMap[fmt.Sprintf("%s%d%s", k, tx.Height, tx.TxHash)] = tx
		}
	}

	// 补充数据
	for _, v := range ibcTxs {
		if v.ScTxInfo != nil {
			tx, ok := chainHashTxMap[fmt.Sprintf("%s%d%s", v.ScChainId, v.ScTxInfo.Height, v.ScTxInfo.Hash)]
			if ok {
				v.ScTxInfo.Memo = tx.Memo
				v.ScTxInfo.Signers = tx.Signers
				v.ScTxInfo.Log = tx.Log
				if tx.Status == entity.TxStatusSuccess {
					for msgIndex, msg := range tx.DocTxMsgs {
						if msg.Type == constant.MsgTypeTransfer && msg.CommonMsg().PacketId == v.ScTxInfo.Msg.CommonMsg().PacketId {
							_, _, _, _, scConnection := parseTransferTxEvents(msgIndex, tx)
							v.ScConnectionId = scConnection
						}
					}
				}
			}
		}

		if v.DcTxInfo != nil {
			tx, ok := chainHashTxMap[fmt.Sprintf("%s%d%s", v.DcChainId, v.DcTxInfo.Height, v.DcTxInfo.Hash)]
			if ok {
				v.DcTxInfo.Memo = tx.Memo
				v.DcTxInfo.Signers = tx.Signers
				v.DcTxInfo.Log = tx.Log
				if tx.Status == entity.TxStatusSuccess {
					for msgIndex, msg := range tx.DocTxMsgs {
						if msg.Type == constant.MsgTypeRecvPacket && msg.CommonMsg().PacketId == v.DcTxInfo.Msg.CommonMsg().PacketId {
							dcConnection, _, _ := parseRecvPacketTxEvents(msgIndex, tx)
							v.DcConnectionId = dcConnection
						}
					}
				}
			}
		}

		if v.RefundedTxInfo != nil {
			tx, ok := chainHashTxMap[fmt.Sprintf("%s%d%s", v.ScChainId, v.RefundedTxInfo.Height, v.RefundedTxInfo.Hash)]
			if ok {
				v.RefundedTxInfo.Memo = tx.Memo
				v.RefundedTxInfo.Signers = tx.Signers
				v.RefundedTxInfo.Log = tx.Log
			}
		}

		if cf, ok := t.chainMap[v.ScChainId]; ok {
			v.ScClientId = cf.GetChannelClient(v.ScPort, v.ScChannel)
		}
		if cf, ok := t.chainMap[v.DcChainId]; ok {
			v.DcClientId = cf.GetChannelClient(v.DcPort, v.DcChannel)
		}

		if err := ibcTxRepo.FixIbxTx(v, isTargetHistory); err != nil {
			logrus.Errorf("task %s FixIbxTx(%s) err, %v", t.Name(), v.RecordId, err)
		}
	}
}
