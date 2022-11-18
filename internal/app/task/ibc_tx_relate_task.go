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
var relateCoordinator *stringQueueCoordinator

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
		chainQueue.Push(v.ChainName)
	}
	relateCoordinator = &stringQueueCoordinator{
		stringQueue: chainQueue,
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
		chain, err := w.getChain()
		if err != nil {
			logrus.Infof("task %s worker %s exit", w.taskName, w.workerName)
			break
		}

		if cf, ok := w.chainMap[chain]; ok && cf.Status == entity.ChainStatusClosed {
			logrus.Infof("task %s worker %s chain %s is closed", w.taskName, w.workerName, chain)
			continue
		}

		logrus.Infof("task %s worker %s get chain: %v", w.taskName, w.workerName, chain)
		startTime := time.Now().Unix()
		if err = w.relateTx(chain); err != nil {
			logrus.Errorf("task %s worker %s relate chain %s tx error,time use: %d(s), %v", w.taskName, w.workerName, chain, time.Now().Unix()-startTime, err)
		} else {
			logrus.Infof("task %s worker %s relate chain %s tx end,time use: %d(s)", w.taskName, w.workerName, chain, time.Now().Unix()-startTime)
		}
	}
}

func (w *ibcTxRelateWorker) getChain() (string, error) {
	if w.target == ibcTxTargetHistory {
		return relateHistoryCoordinator.getOne()
	}
	return relateCoordinator.getOne()
}

func (w *ibcTxRelateWorker) relateTx(chain string) error {
	totalRelateTx := 0
	//const limit = 500
	maxParseTx := global.Config.Task.SingleChainIbcTxRelateMax
	if maxParseTx <= 0 {
		maxParseTx = defaultMaxHandlerTx
	}

	denomMap, err := w.getChainDenomMap(chain)
	if err != nil {
		return err
	}

	for {
		txList, err := w.getToBeRelatedTxs(chain, constant.DefaultLimit)
		if err != nil {
			logrus.Errorf("task %s worker %s chain %s getToBeRelatedTxs error, %v", w.taskName, w.workerName, chain, err)
			return err
		}

		if len(txList) == 0 {
			return nil
		}

		w.handlerIbcTxs(chain, txList, denomMap)

		totalRelateTx += len(txList)
		if len(txList) < constant.DefaultLimit || totalRelateTx >= maxParseTx {
			break
		}
		time.Sleep(200 * time.Millisecond) // avoid master-slave delay problem
	}

	return nil
}

func (w *ibcTxRelateWorker) handlerIbcTxs(scChainId string, ibcTxList []*entity.ExIbcTx, denomMap map[string]*entity.IBCDenom) {
	recvPacketTxMap, ackTxMap, refundedTxMap, timeoutIbcTxMap, noFoundAckMap := w.packetIdTx(scChainId, ibcTxList)

	var ibcDenomNewList entity.IBCDenomList
	for _, ibcTx := range ibcTxList {
		if ibcTx.DcChain == "" || ibcTx.ScTxInfo == nil || ibcTx.ScTxInfo.Msg == nil {
			w.setNextTryTime(ibcTx)
		} else {
			packetId := ibcTx.ScTxInfo.Msg.CommonMsg().PacketId
			if syncTxs, ok := recvPacketTxMap[w.genPacketTxMapKey(ibcTx.DcChain, packetId)]; ok {
				ackSyncTxs := ackTxMap[w.genPacketTxMapKey(ibcTx.ScChain, packetId)]
				ibcDenom := w.loadRecvPacketTx(ibcTx, syncTxs, ackSyncTxs)
				if ibcDenom != nil && denomMap[ibcDenom.Denom] == nil {
					denomMap[ibcDenom.Denom] = ibcDenom
					ibcDenomNewList = append(ibcDenomNewList, ibcDenom)
				}
			}

			if syncTxs, ok := refundedTxMap[w.genPacketTxMapKey(ibcTx.ScChain, packetId)]; ok && ibcTx.Status == entity.IbcTxStatusProcessing {
				recvSyncTxs := recvPacketTxMap[w.genPacketTxMapKey(ibcTx.DcChain, packetId)]
				w.loadTimeoutPacketTx(ibcTx, syncTxs, recvSyncTxs)
			}
		}

		if ibcTx.Status == entity.IbcTxStatusProcessing {
			w.setNextTryTime(ibcTx)
			//记录"处理中"状态
			ibcTx = w.updateProcessInfo(ibcTx, timeoutIbcTxMap, noFoundAckMap)
		}
		var repaired bool
		ibcTx, repaired = w.repairTxInfo(ibcTx)
		if err := w.updateIbcTx(ibcTx, repaired); err != nil {
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
		if ibcTx.DcChain == "" {
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

func (w *ibcTxRelateWorker) loadRecvPacketTx(ibcTx *entity.ExIbcTx, txs, ackTxs []*entity.Tx) *entity.IBCDenom {
	refundedMatchAckTx := func() *entity.Tx {
		var matchTx *entity.Tx
		for _, ackTx := range ackTxs {
			for msgIndex, msg := range ackTx.DocTxMsgs {
				if msg.Type != constant.MsgTypeAcknowledgement || msg.CommonMsg().PacketId != ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
					continue
				}
				existTransferEvent := parseAckPacketTxEvents(msgIndex, ackTx)
				if existTransferEvent {
					matchTx = ackTx
				}
			}
		}
		return matchTx
	}

	successMatchAckTx := func() *entity.Tx {
		var matchTx *entity.Tx
		for _, ackTx := range ackTxs {
			for _, msg := range ackTx.DocTxMsgs {
				if msg.Type != constant.MsgTypeAcknowledgement || msg.CommonMsg().PacketId != ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
					continue
				}
				if matchTx == nil {
					matchTx = ackTx
				} else {
					if ackTx.Time > matchTx.Time {
						matchTx = ackTx
					}
				}
			}
		}
		return matchTx
	}

	var ibcDenom *entity.IBCDenom
	var matchAckTx *entity.Tx
	for _, tx := range txs {
		if tx.Status == entity.TxStatusFailed {
			continue
		}
		for msgIndex, msg := range tx.DocTxMsgs {
			if msg.Type != constant.MsgTypeRecvPacket || msg.CommonMsg().PacketId != ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
				continue
			}
			dcConnection, packetAck, existPacketAck := parseRecvPacketTxEvents(msgIndex, tx)
			if !existPacketAck {
				continue
			}

			if strings.Contains(packetAck, "error") {
				matchAckTx = refundedMatchAckTx()
				if matchAckTx == nil { // 改为refunded状态时，必须要ack_packet交易
					return nil
				}
				ibcTx.Status = entity.IbcTxStatusRefunded
			} else {
				matchAckTx = successMatchAckTx()
				ibcTx.Status = entity.IbcTxStatusSuccess
			}

			ibcTx.DcConnectionId = dcConnection
			ibcTx.DcTxInfo = &entity.TxInfo{
				Hash:      tx.TxHash,
				Status:    tx.Status,
				Time:      tx.Time,
				Height:    tx.Height,
				Fee:       tx.Fee,
				MsgAmount: nil,
				Msg:       msg,
				Memo:      tx.Memo,
				Signers:   tx.Signers,
				Log:       tx.Log,
			}
			ibcTx.UpdateAt = time.Now().Unix()

			dcDenomFullPath, isCrossBack := calculateNextDenomPath(msg.RecvPacketMsg().Packet)
			dcDenom := calculateIbcHash(dcDenomFullPath)
			ibcTx.Denoms.DcDenom = dcDenom // set ibc tx dc denom
			if ibcTx.Status == entity.IbcTxStatusSuccess {
				if !isCrossBack {
					dcDenomPath, rootDenom := splitFullPath(dcDenomFullPath)
					ibcDenom = &entity.IBCDenom{
						Symbol:         "",
						Chain:          ibcTx.DcChain,
						Denom:          dcDenom,
						PrevDenom:      ibcTx.Denoms.ScDenom,
						PrevChain:      ibcTx.ScChain,
						BaseDenom:      ibcTx.BaseDenom,
						BaseDenomChain: ibcTx.BaseDenomChain,
						DenomPath:      dcDenomPath,
						IsBaseDenom:    false,
						RootDenom:      rootDenom,
						CreateAt:       time.Now().Unix(),
						UpdateAt:       time.Now().Unix(),
					}
				}
			}
		}
	}

	if matchAckTx != nil && ibcTx.Status != entity.IbcTxStatusProcessing {
		for _, msg := range matchAckTx.DocTxMsgs {
			if msg.Type != constant.MsgTypeAcknowledgement || msg.CommonMsg().PacketId != ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
				continue
			}
			ibcTx.RefundedTxInfo = &entity.TxInfo{
				Hash:      matchAckTx.TxHash,
				Status:    matchAckTx.Status,
				Time:      matchAckTx.Time,
				Height:    matchAckTx.Height,
				Fee:       matchAckTx.Fee,
				MsgAmount: nil,
				Msg:       msg,
				Memo:      matchAckTx.Memo,
				Signers:   matchAckTx.Signers,
				Log:       matchAckTx.Log,
			}
		}
	}
	return ibcDenom
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
					Memo:      tx.Memo,
					Signers:   tx.Signers,
					Log:       tx.Log,
				}
				ibcTx.UpdateAt = time.Now().Unix()
			} else {
				w.setNextTryTime(ibcTx)
			}
			return
		}
	}
}

func (w *ibcTxRelateWorker) loadTimeoutPacketTx(ibcTx *entity.ExIbcTx, tx *entity.Tx, recvSyncTxs []*entity.Tx) {
	packetId := ibcTx.ScTxInfo.Msg.CommonMsg().PacketId
	for _, msg := range tx.DocTxMsgs {
		if msg.Type == constant.MsgTypeTimeoutPacket && msg.CommonMsg().PacketId == packetId {
			ibcTx.Status = entity.IbcTxStatusRefunded
			ibcTx.RefundedTxInfo = &entity.TxInfo{
				Hash:      tx.TxHash,
				Status:    tx.Status,
				Time:      tx.Time,
				Height:    tx.Height,
				Fee:       tx.Fee,
				MsgAmount: nil,
				Msg:       msg,
				Memo:      tx.Memo,
				Signers:   tx.Signers,
				Log:       tx.Log,
			}
		}
	}

	if ibcTx.Status != entity.IbcTxStatusProcessing {
		var matchRecvTx *entity.Tx
		var matchRecvTxMsg *model.TxMsg
		for _, recvTx := range recvSyncTxs {
			for _, msg := range recvTx.DocTxMsgs {
				if msg.Type != constant.MsgTypeRecvPacket || msg.CommonMsg().PacketId != ibcTx.ScTxInfo.Msg.CommonMsg().PacketId {
					continue
				}

				if matchRecvTx == nil {
					matchRecvTx = recvTx
					matchRecvTxMsg = msg
				} else {
					if recvTx.Time > matchRecvTx.Time {
						matchRecvTx = recvTx
						matchRecvTxMsg = msg
					}
				}
			}
		}

		if matchRecvTx != nil {
			ibcTx.DcTxInfo = &entity.TxInfo{
				Hash:      matchRecvTx.TxHash,
				Status:    matchRecvTx.Status,
				Time:      matchRecvTx.Time,
				Height:    matchRecvTx.Height,
				Fee:       matchRecvTx.Fee,
				MsgAmount: nil,
				Msg:       matchRecvTxMsg,
				Memo:      matchRecvTx.Memo,
				Signers:   matchRecvTx.Signers,
				Log:       matchRecvTx.Log,
			}
		}
	}
}

func (w *ibcTxRelateWorker) setNextTryTime(ibcTx *entity.ExIbcTx) {
	now := time.Now().Unix()
	ibcTx.RetryTimes += 1
	ibcTx.NextTryTime = now + (ibcTx.RetryTimes * 2)
	ibcTx.UpdateAt = time.Now().Unix()
}

func (w *ibcTxRelateWorker) repairTxInfo(ibcTx *entity.ExIbcTx) (*entity.ExIbcTx, bool) {
	var repaired bool
	if ibcTx.DcClientId == "" {
		if cf, ok := w.chainMap[ibcTx.DcChain]; ok {
			ibcTx.DcClientId = cf.GetChannelClient(ibcTx.DcPort, ibcTx.DcChannel)
			repaired = true
		}
	}

	if ibcTx.ScClientId == "" {
		if cf, ok := w.chainMap[ibcTx.ScChain]; ok {
			ibcTx.ScClientId = cf.GetChannelClient(ibcTx.ScPort, ibcTx.ScChannel)
			repaired = true
		}
	}

	if ibcTx.ScConnectionId == "" {
		packetId := ibcTx.ScTxInfo.Msg.CommonMsg().PacketId
		status := entity.TxStatusSuccess
		if transferTxs, err := txRepo.FindByPacketIds(ibcTx.ScChain, constant.MsgTypeTransfer, []string{packetId}, &status); err == nil && len(transferTxs) > 0 {
			for msgIndex, msg := range transferTxs[0].DocTxMsgs {
				if msg.Type == constant.MsgTypeTransfer && msg.CommonMsg().PacketId == packetId {
					_, _, _, _, ibcTx.ScConnectionId = parseTransferTxEvents(msgIndex, transferTxs[0])
					repaired = true
					break
				}
			}
		}
	}
	return ibcTx, repaired
}

func (w *ibcTxRelateWorker) genPacketTxMapKey(chainId, packetId string) string {
	return fmt.Sprintf("%s_%s", chainId, packetId)
}

func (w *ibcTxRelateWorker) packetIdTx(scChainId string, ibcTxList []*entity.ExIbcTx) (recvPacketTxMap, ackTxMap map[string][]*entity.Tx, refundedTxMap map[string]*entity.Tx, timeoutIbcTxMap, noFoundAckMap map[string]struct{}) {
	packetIdsMap := w.packetIdsMap(ibcTxList)
	chainLatestBlockMap := w.findLatestBlock(scChainId, ibcTxList)
	var refundedTxPacketIds, ackPacketIds []string
	recvPacketTxMap = make(map[string][]*entity.Tx)
	ackTxMap = make(map[string][]*entity.Tx)
	refundedTxMap = make(map[string]*entity.Tx)
	timeoutIbcTxMap = make(map[string]struct{})
	noFoundAckMap = make(map[string]struct{})
	packetIdRecordMap := make(map[string]string, len(packetIdsMap))
	status := entity.TxStatusSuccess

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
					//logrus.Warningf("unkonwn timeout time %s, chain: %s, packet id: %s", timeoutStr, dcChainId, packet.PacketId)
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
						mk := w.genPacketTxMapKey(dcChainId, recvMsg.PacketId)
						recvPacketTxMap[mk] = append(recvPacketTxMap[mk], tx)
						// recv_packet成功时查询ack
						if tx.Status == entity.TxStatusSuccess {
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
		ackTxList, err := txRepo.FindByPacketIds(scChainId, constant.MsgTypeAcknowledgement, ackPacketIds, &status)
		if err == nil {
			for _, tx := range ackTxList {
				for _, msg := range tx.DocTxMsgs {
					if msg.Type == constant.MsgTypeAcknowledgement {
						if ackMsg := msg.AckPacketMsg(); ackMsg.PacketId != "" {
							mk := w.genPacketTxMapKey(scChainId, ackMsg.PacketId)
							ackTxMap[mk] = append(ackTxMap[mk], tx)
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
	return
}

func (w *ibcTxRelateWorker) packetIdsMap(ibcTxList []*entity.ExIbcTx) map[string][]*dto.PacketIdDTO {
	res := make(map[string][]*dto.PacketIdDTO)
	for _, tx := range ibcTxList {
		if tx.DcChain == "" || tx.ScTxInfo == nil || tx.ScTxInfo.Msg == nil {
			logrus.Warningf("ibc tx dc_chain_id or sc_tx_info exception, record_id: %s", tx.RecordId)
			continue
		}

		transferMsg := tx.ScTxInfo.Msg.TransferMsg()
		if transferMsg.PacketId == "" {
			logrus.Warningf("ibc tx packet_id is empty, hash: %s", tx.ScTxInfo.Hash)
			continue
		}

		res[tx.DcChain] = append(res[tx.DcChain], &dto.PacketIdDTO{
			DcChainId:     tx.DcChain,
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
		if tx.DcChain == "" {
			continue
		}

		if _, ok := blockMap[tx.DcChain]; !ok {
			findFunc(tx.DcChain)
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

func (w *ibcTxRelateWorker) updateIbcTx(ibcTx *entity.ExIbcTx, repaired bool) error {
	if w.target == ibcTxTargetHistory {
		return ibcTxRepo.UpdateIbcHistoryTx(ibcTx, repaired)
	}
	return ibcTxRepo.UpdateIbcTx(ibcTx, repaired)
}

func (w *ibcTxRelateWorker) getChainDenomMap(chain string) (map[string]*entity.IBCDenom, error) {
	denomList, err := denomRepo.FindByChain(chain)
	if err != nil {
		logrus.Errorf("task %s worker %s getChainDenomMap %s error, %v", w.taskName, w.workerName, chain, err)
		return nil, err
	}

	denomMap := make(map[string]*entity.IBCDenom, len(denomList))
	for _, v := range denomList {
		denomMap[v.Denom] = v
	}
	return denomMap, nil
}
