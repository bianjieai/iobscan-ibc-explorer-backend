package task

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IbcTxRelateTask struct {
}

var _ Task = new(IbcTxRelateTask)
var relateCoordinator *chainQueueCoordinator

func (t *IbcTxRelateTask) Name() string {
	return "ibc_tx_relate_task"
}

func (t *IbcTxRelateTask) Cron() int {
	if taskConf.CronTimeIbcTxRelateTask > 0 {
		return taskConf.CronTimeIbcTxRelateTask
	}
	return ThreeMinute
}

func (t *IbcTxRelateTask) workerNum() int {
	if global.Config.Task.IbcTxRelateWorkerNum > 0 {
		return global.Config.Task.IbcTxRelateWorkerNum
	}
	return ibcTxRelateTaskWorkerNum
}

func (t *IbcTxRelateTask) Run() int {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", t.Name(), err)
		return -1
	}

	// init coordinator
	chainQueue := new(utils.QueueString)
	for _, v := range chainMap {
		chainQueue.Push(v.ChainId)
	}
	relateCoordinator = &chainQueueCoordinator{
		chainQueue: chainQueue,
	}

	workerNum := t.workerNum()
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	for i := 1; i <= workerNum; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			newIbcTxRelateWorker(t.Name(), wn, ibcTxTargetLatest, chainMap).exec()
			waitGroup.Done()
		}(workName)
	}
	waitGroup.Wait()

	return 1
}

// =========================================================================
// =========================================================================
// worker

func newIbcTxRelateWorker(taskName, workerName, target string, chainMap map[string]*entity.ChainConfig) *ibcTxRelateWorker {
	return &ibcTxRelateWorker{
		taskName:   taskName,
		workerName: workerName,
		target:     target,
		chainMap:   chainMap,
	}
}

type ibcTxRelateWorker struct {
	taskName   string
	workerName string
	target     string
	chainMap   map[string]*entity.ChainConfig
}

func (w *ibcTxRelateWorker) exec() {
	logrus.Infof("task %s worker %s start", w.taskName, w.workerName)
	for {
		chainId, err := w.getChain()
		if err != nil {
			logrus.Infof("task %s worker %s exit", w.taskName, w.workerName)
			break
		}

		logrus.Infof("task %s worker %s get chain: %v", w.taskName, w.workerName, chainId)
		startTime := time.Now().Unix()
		if err = w.relateTx(chainId); err != nil {
			logrus.Errorf("task %s worker %s relate chain %s tx error,time use: %d(s), %v", w.taskName, w.workerName, chainId, time.Now().Unix()-startTime, err)
		} else {
			logrus.Infof("task %s worker %s relate chain %s tx end,time use: %d(s)", w.taskName, w.workerName, chainId, time.Now().Unix()-startTime)
		}
	}
}

func (w *ibcTxRelateWorker) getChain() (string, error) {
	if w.target == ibcTxTargetHistory {
		return relateHistoryCoordinator.getChain()
	}
	return relateCoordinator.getChain()
}

func (w *ibcTxRelateWorker) relateTx(chainId string) error {
	totalRelateTx := 0
	//const limit = 500
	maxParseTx := global.Config.Task.SingleChainIbcTxRelateMax
	if maxParseTx <= 0 {
		maxParseTx = defaultMaxHandlerTx
	}

	denomMap, err := w.getChainDenomMap(chainId)
	if err != nil {
		return err
	}

	for {
		txList, err := w.getToBeRelatedTxs(chainId, constant.DefaultLimit)
		if err != nil {
			logrus.Errorf("task %s worker %s chain %s getToBeRelatedTxs error, %v", w.taskName, w.workerName, chainId, err)
			return err
		}

		if len(txList) == 0 {
			return nil
		}

		w.handlerIbcTxs(chainId, txList, denomMap)

		totalRelateTx += len(txList)
		if len(txList) < constant.DefaultLimit || totalRelateTx >= maxParseTx {
			break
		}
		time.Sleep(1 * time.Second) // avoid master-slave delay problem
	}

	return nil
}

func (w *ibcTxRelateWorker) handlerIbcTxs(scChainId string, ibcTxList []*entity.ExIbcTx, denomMap map[string]*entity.IBCDenom) {
	recvPacketTxMap, refundedTxMap, ackTxMap, timeoutIbcTxMap, noFoundAckMap := w.packetIdTx(scChainId, ibcTxList)

	var ibcDenomNewList entity.IBCDenomList
	for _, ibcTx := range ibcTxList {
		if ibcTx.DcChainId == "" || ibcTx.ScTxInfo == nil || ibcTx.ScTxInfo.Msg == nil {
			w.setNextTryTime(ibcTx)
		} else if syncTx, ok := recvPacketTxMap[w.genPacketTxMapKey(ibcTx.DcChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)]; ok {
			ibcDenom := w.loadRecvPacketTx(ibcTx, syncTx)
			if ibcDenom != nil && denomMap[ibcDenom.Denom] == nil {
				denomMap[ibcDenom.Denom] = ibcDenom
				ibcDenomNewList = append(ibcDenomNewList, ibcDenom)
			}
		} else if syncTx, ok = ackTxMap[w.genPacketTxMapKey(ibcTx.ScChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)]; ok {
			w.loadAckPacketTx(ibcTx, syncTx)
		} else if syncTx, ok = refundedTxMap[w.genPacketTxMapKey(ibcTx.ScChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)]; ok {
			w.loadTimeoutPacketTx(ibcTx, syncTx)
		} else {
			w.setNextTryTime(ibcTx)
		}

		//记录"处理中"状态
		ibcTx = w.updateProcessInfo(ibcTx, timeoutIbcTxMap, noFoundAckMap)
		if err := w.updateIbcTx(ibcTx); err != nil {
			logrus.Errorf("task %s worker %s chain %s updateIbcTx error, record_id: %s, %v", w.taskName, w.workerName, scChainId, ibcTx.RecordId, err)
		}
	}

	// add denom
	if len(ibcDenomNewList) > 0 {
		if err := denomRepo.InsertBatch(ibcDenomNewList); err != nil {
			logrus.Errorf("task %s worker %s chain %s insert denoms error, %v", w.taskName, w.workerName, scChainId, err)
		}
	}
}

func (w *ibcTxRelateWorker) updateProcessInfo(ibcTx *entity.ExIbcTx, timeOutMap map[string]struct{}, noFoundAckMap map[string]struct{}) *entity.ExIbcTx {
	if ibcTx.Status == entity.IbcTxStatusProcessing {
		if ibcTx.DcChainId == "" {
			ibcTx.ProcessInfo = constant.NoFoundDcChainId
		} else {
			if _, ok := timeOutMap[ibcTx.RecordId]; ok {
				ibcTx.ProcessInfo = constant.NoFoundSuccessTimeoutPacket
			} else if _, ok := noFoundAckMap[ibcTx.RecordId]; ok {
				ibcTx.ProcessInfo = constant.NoFoundSuccessAcknowledgePacket
			} else {
				ibcTx.ProcessInfo = constant.NoFoundSuccessRecvPacket
			}
		}
	} else {
		ibcTx.ProcessInfo = ""
	}
	return ibcTx
}

func (w *ibcTxRelateWorker) loadRecvPacketTx(ibcTx *entity.ExIbcTx, tx *entity.Tx) *entity.IBCDenom {
	var writeAckRes bool
	var recvMsg *model.TxMsg
	for msgIndex, msg := range tx.DocTxMsgs {
		if msg.Type != constant.MsgTypeRecvPacket || msg.CommonMsg().PacketId != ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
			continue
		}

		recvMsg = msg
		for _, event := range tx.EventsNew {
			if event.MsgIndex != uint32(msgIndex) {
				continue
			}

			for _, ee := range event.Events {
				if ee.Type == "write_acknowledgement" {
					for _, attr := range ee.Attributes {
						if attr.Key == "packet_ack" {
							writeAckRes = !strings.Contains(attr.Value, "error")
							break
						}
					}
					break
				}
			}
		}

		if writeAckRes {
			ibcTx.Status = entity.IbcTxStatusSuccess
		} else {
			ibcTx.Status = entity.IbcTxStatusFailed
		}
		ibcTx.DcTxInfo = &entity.TxInfo{
			Hash:      tx.TxHash,
			Status:    tx.Status,
			Time:      tx.Time,
			Height:    tx.Height,
			Fee:       tx.Fee,
			MsgAmount: nil,
			Msg:       recvMsg,
		}
		ibcTx.UpdateAt = time.Now().Unix()

		dcDenomPath, isCrossBack := calculateNextDenomPath(recvMsg.RecvPacketMsg().Packet)
		dcDenom := calculateIbcHash(dcDenomPath)
		ibcTx.Denoms.DcDenom = dcDenom // set ibc tx dc denom

		if ibcTx.Status == entity.IbcTxStatusSuccess {
			if !isCrossBack {
				rootDenom := getRootDenom(dcDenomPath)
				return &entity.IBCDenom{
					Symbol:           "",
					ChainId:          ibcTx.DcChainId,
					Denom:            dcDenom,
					PrevDenom:        ibcTx.Denoms.ScDenom,
					PrevChainId:      ibcTx.ScChainId,
					BaseDenom:        ibcTx.BaseDenom,
					BaseDenomChainId: ibcTx.BaseDenomChainId,
					DenomPath:        dcDenomPath,
					IsBaseDenom:      false,
					RootDenom:        rootDenom,
					CreateAt:         time.Now().Unix(),
					UpdateAt:         time.Now().Unix(),
				}
			}
		}
	} // for
	return nil
}

func (w *ibcTxRelateWorker) loadAckPacketTx(ibcTx *entity.ExIbcTx, tx *entity.Tx) {
	for _, msg := range tx.DocTxMsgs {
		if msg.Type == constant.MsgTypeAcknowledgement && msg.CommonMsg().PacketId == ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
			if strings.Contains(msg.AckPacketMsg().Acknowledgement, "error") { // ack error
				ibcTx.Status = entity.IbcTxStatusRefunded
				ibcTx.RefundedTxInfo = &entity.TxInfo{
					Hash:      tx.TxHash,
					Status:    tx.Status,
					Time:      tx.Time,
					Height:    tx.Height,
					Fee:       tx.Fee,
					MsgAmount: nil,
					Msg:       msg,
				}
				ibcTx.UpdateAt = time.Now().Unix()
			} else {
				w.setNextTryTime(ibcTx)
			}
			return
		}
	}
}

func (w *ibcTxRelateWorker) loadTimeoutPacketTx(ibcTx *entity.ExIbcTx, tx *entity.Tx) {
	for _, msg := range tx.DocTxMsgs {
		if msg.Type == constant.MsgTypeTimeoutPacket && msg.CommonMsg().PacketId == ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
			ibcTx.Status = entity.IbcTxStatusRefunded
			ibcTx.RefundedTxInfo = &entity.TxInfo{
				Hash:      tx.TxHash,
				Status:    tx.Status,
				Time:      tx.Time,
				Height:    tx.Height,
				Fee:       tx.Fee,
				MsgAmount: nil,
				Msg:       msg,
			}
		}
	}
}

func (w *ibcTxRelateWorker) setNextTryTime(ibcTx *entity.ExIbcTx) {
	now := time.Now().Unix()
	ibcTx.RetryTimes += 1
	ibcTx.NextTryTime = now + (ibcTx.RetryTimes * ThreeMinute)
	ibcTx.UpdateAt = time.Now().Unix()
}

func (w *ibcTxRelateWorker) genPacketTxMapKey(chainId, packetId string) string {
	return fmt.Sprintf("%s_%s", chainId, packetId)
}

func (w *ibcTxRelateWorker) packetIdTx(scChainId string, ibcTxList []*entity.ExIbcTx) (map[string]*entity.Tx, map[string]*entity.Tx, map[string]*entity.Tx, map[string]struct{}, map[string]struct{}) {
	packetIdsMap := w.packetIdsMap(ibcTxList)
	chainLatestBlockMap := w.findLatestBlock(scChainId, ibcTxList)
	var refundedTxPacketIds, ackPacketIds []string
	recvPacketTxMap := make(map[string]*entity.Tx)
	refundedTxMap := make(map[string]*entity.Tx)
	ackTxMap := make(map[string]*entity.Tx)
	timeoutIbcTxMap := make(map[string]struct{})
	noFoundAckMap := make(map[string]struct{})
	packetIdRecordMap := make(map[string]string, len(packetIdsMap))

	for dcChainId, packetIds := range packetIdsMap {
		latestBlock := chainLatestBlockMap[dcChainId]
		var recvPacketIds []string
		for _, packet := range packetIds { // recv && refunded
			packetIdRecordMap[dcChainId+packet.PacketId] = packet.RecordId
			recvPacketIds = append(recvPacketIds, packet.PacketId)
			timeoutStr := strconv.FormatInt(packet.TimeOutTime, 10)
			if len(timeoutStr) > 10 { // 非秒级时间
				if len(timeoutStr) == 19 && time.Now().UnixNano() > packet.TimeOutTime { // Nano
					refundedTxPacketIds = append(refundedTxPacketIds, packet.PacketId)
					timeoutIbcTxMap[packet.RecordId] = struct{}{}
				} else {
					logrus.Warningf("unkonwn timeout time %s, chain: %s, packet id: %s", timeoutStr, dcChainId, packet.PacketId)
					refundedTxPacketIds = append(refundedTxPacketIds, packet.PacketId)
				}
			} else if latestBlock != nil {
				if latestBlock.Height > packet.TimeoutHeight || latestBlock.Time > packet.TimeOutTime {
					refundedTxPacketIds = append(refundedTxPacketIds, packet.PacketId)
					timeoutIbcTxMap[packet.RecordId] = struct{}{}
				}
			}
		}

		// 处理 recv_packet tx
		recvTxList, err := txRepo.FindByPacketIds(dcChainId, constant.MsgTypeRecvPacket, recvPacketIds, nil)
		if err != nil {
			logrus.Errorf("task %s worker %s dc chain %s find recv txs error, %v", w.taskName, w.workerName, dcChainId, err)
			continue
		}

		for _, tx := range recvTxList {
			for _, msg := range tx.DocTxMsgs {
				if msg.Type == constant.MsgTypeRecvPacket {
					if recvMsg := msg.RecvPacketMsg(); recvMsg.PacketId != "" {
						if tx.Status == entity.TxStatusSuccess {
							recvPacketTxMap[w.genPacketTxMapKey(dcChainId, recvMsg.PacketId)] = tx
						} else {
							ackPacketIds = append(ackPacketIds, recvMsg.PacketId)
							if recordId, ok := packetIdRecordMap[dcChainId+recvMsg.PacketId]; ok {
								noFoundAckMap[recordId] = struct{}{}
							}
						}
					} else {
						logrus.Errorf("%s recv packet tx(%s) packet id is empty", dcChainId, tx.TxHash)
					}
				}
			}
		}
	}

	if len(refundedTxPacketIds) > 0 {
		status := entity.TxStatusSuccess
		refundedTxList, err := txRepo.FindByPacketIds(scChainId, constant.MsgTypeTimeoutPacket, refundedTxPacketIds, &status)
		if err == nil {
			for _, tx := range refundedTxList {
				for _, msg := range tx.DocTxMsgs {
					if msg.Type == constant.MsgTypeTimeoutPacket {
						if timeoutMsg := msg.TimeoutPacketMsg(); timeoutMsg.PacketId != "" {
							refundedTxMap[w.genPacketTxMapKey(scChainId, timeoutMsg.PacketId)] = tx
						} else {
							logrus.Errorf("%s timeout packet tx(%s) packet id is empty", scChainId, tx.TxHash)
						}
					}
				}
			}
		} else {
			logrus.Errorf("task %s worker %s sc chain %s find timeout txs error, %v", w.taskName, w.workerName, scChainId, err)
		}
	}

	if len(ackPacketIds) > 0 {
		status := entity.TxStatusSuccess
		ackTxList, err := txRepo.FindByPacketIds(scChainId, constant.MsgTypeAcknowledgement, ackPacketIds, &status)
		if err == nil {
			for _, tx := range ackTxList {
				for _, msg := range tx.DocTxMsgs {
					if msg.Type == constant.MsgTypeAcknowledgement {
						if ackMsg := msg.AckPacketMsg(); ackMsg.PacketId != "" {
							ackTxMap[w.genPacketTxMapKey(scChainId, ackMsg.PacketId)] = tx
						} else {
							logrus.Errorf("%s ack packet tx(%s) packet id is empty", scChainId, tx.TxHash)
						}
					}
				}
			}
		} else {
			logrus.Errorf("task %s worker %s sc chain %s find ack txs error, %v", w.taskName, w.workerName, scChainId, err)
		}
	}

	return recvPacketTxMap, refundedTxMap, ackTxMap, timeoutIbcTxMap, noFoundAckMap
}

func (w *ibcTxRelateWorker) packetIdsMap(ibcTxList []*entity.ExIbcTx) map[string][]*dto.PacketIdDTO {
	res := make(map[string][]*dto.PacketIdDTO)
	for _, tx := range ibcTxList {
		if tx.DcChainId == "" || tx.ScTxInfo == nil || tx.ScTxInfo.Msg == nil {
			logrus.Warningf("ibc tx dc_chain_id or sc_tx_info exception, record_id: %s", tx.RecordId)
			continue
		}

		transferMsg := tx.ScTxInfo.Msg.TransferMsg()
		if transferMsg.PacketId == "" {
			logrus.Warningf("ibc tx packet_id is empty, hash: %s", tx.ScTxInfo.Hash)
			continue
		}

		res[tx.DcChainId] = append(res[tx.DcChainId], &dto.PacketIdDTO{
			DcChainId:     tx.DcChainId,
			TimeoutHeight: transferMsg.TimeoutHeight.RevisionHeight,
			PacketId:      transferMsg.PacketId,
			TimeOutTime:   transferMsg.TimeoutTimestamp,
			RecordId:      tx.RecordId,
		})
	}
	return res
}

func (w *ibcTxRelateWorker) findLatestBlock(scChainId string, ibcTxList []*entity.ExIbcTx) map[string]*dto.HeightTimeDTO {
	blockMap := make(map[string]*dto.HeightTimeDTO)

	findFunc := func(chainId string) {
		block, err := syncBlockRepo.FindLatestBlock(chainId)
		if err != nil {
			logrus.Errorf("task %s worker %s chain %s findLatestBlock error, %v", w.taskName, w.workerName, chainId, err)
		} else {
			blockMap[chainId] = &dto.HeightTimeDTO{
				Height: block.Height,
				Time:   block.Time,
			}
		}
	}

	findFunc(scChainId)
	for _, tx := range ibcTxList {
		if tx.DcChainId == "" {
			continue
		}

		if _, ok := blockMap[tx.DcChainId]; !ok {
			findFunc(tx.DcChainId)
		}
	}

	return blockMap
}

func (w *ibcTxRelateWorker) getToBeRelatedTxs(chainId string, limit int64) ([]*entity.ExIbcTx, error) {
	if w.target == ibcTxTargetHistory {
		return ibcTxRepo.FindProcessingHistoryTxs(chainId, limit)
	}
	return ibcTxRepo.FindProcessingTxs(chainId, limit)
}

func (w *ibcTxRelateWorker) updateIbcTx(ibcTx *entity.ExIbcTx) error {
	if w.target == ibcTxTargetHistory {
		return ibcTxRepo.UpdateIbcHistoryTx(ibcTx)
	}
	return ibcTxRepo.UpdateIbcTx(ibcTx)
}

func (w *ibcTxRelateWorker) getChainDenomMap(chainId string) (map[string]*entity.IBCDenom, error) {
	denomList, err := denomRepo.FindByChainId(chainId)
	if err != nil {
		logrus.Errorf("task %s worker %s getChainDenomMap %s error, %v", w.taskName, w.workerName, chainId, err)
		return nil, err
	}

	denomMap := make(map[string]*entity.IBCDenom, len(denomList))
	for _, v := range denomList {
		denomMap[v.Denom] = v
	}
	return denomMap, nil
}
