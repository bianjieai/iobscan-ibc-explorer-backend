package task

import (
	"encoding/json"
	"fmt"
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
	"go.mongodb.org/mongo-driver/mongo"
)

type IbcSyncTransferTxTask struct {
}

type syncTransferTxCoordinator struct {
	allChainMap map[string]*entity.ChainConfig
	chainQueue  *utils.QueueString
}

var _ Task = new(IbcSyncTransferTxTask)
var coordinator *syncTransferTxCoordinator

func (t *IbcSyncTransferTxTask) Name() string {
	return "ibc_sync_transfer_tx_task"
}

func (t *IbcSyncTransferTxTask) Cron() int {
	if taskConf.CronTimeSyncTransferTxTask > 0 {
		return taskConf.CronTimeSyncTransferTxTask
	}
	return ThreeMinute
}

func (t *IbcSyncTransferTxTask) Run() int {
	if err := t.initCoordinator(); err != nil {
		return -1
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(syncTransferTxTaskWorkerQuantity)
	for i := 1; i <= syncTransferTxTaskWorkerQuantity; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			newSyncTransferTxWorker(t.Name(), wn).exec()
			waitGroup.Done()
		}(workName)
	}
	waitGroup.Wait()

	return 1
}

func (t *IbcSyncTransferTxTask) initCoordinator() error {
	allChainList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s chainConfigRepo.FindAll error, %v", t.Name(), err)
		return err
	}

	allChainMap := make(map[string]*entity.ChainConfig)
	chainQueue := new(utils.QueueString)
	for _, v := range allChainList {
		allChainMap[v.ChainId] = v
		chainQueue.Push(v.ChainId)
	}

	coordinator = &syncTransferTxCoordinator{
		allChainMap: allChainMap,
		chainQueue:  chainQueue,
	}
	return nil
}

// =========================================================================
// =========================================================================
// worker

func newSyncTransferTxWorker(taskName, workerName string) *syncTransferTxWorker {
	return &syncTransferTxWorker{
		taskName:   taskName,
		workerName: workerName,
	}
}

type syncTransferTxWorker struct {
	taskName   string
	workerName string
}

func (w *syncTransferTxWorker) exec() {
	logrus.Infof("task %s worker %s start", w.taskName, w.workerName)
	for {
		chainId, err := w.getChain()
		logrus.Infof("task %s worker %s get chain: %v", w.taskName, w.workerName, chainId)
		if err != nil {
			break
		}

		startTime := time.Now().Unix()
		if err = w.parseChainIbcTx(chainId); err != nil {
			logrus.Errorf("task %s worker %s parse chain %s tx error,time use: %d(s), %v", w.taskName, w.workerName, chainId, time.Now().Unix()-startTime, err)
		} else {
			logrus.Infof("task %s worker %s parse chain %s tx end,time use: %d(s)", w.taskName, w.workerName, chainId, time.Now().Unix()-startTime)
		}
	}
}

func (w *syncTransferTxWorker) getChain() (string, error) {
	if coordinator == nil || coordinator.chainQueue == nil {
		return "", fmt.Errorf("coordinator or chain queue is nil")
	}

	return coordinator.chainQueue.Pop()
}

func (w *syncTransferTxWorker) parseChainIbcTx(chainId string) error {
	totalParseTx := 0
	limit := 500

	taskRecord, err := w.checkTaskRecord(chainId)
	if err != nil {
		return err
	}
	if taskRecord.Status == entity.TaskRecordStatusClose {
		return nil
	}

	denomMap, err := w.getChainDenomMap(chainId)
	if err != nil {
		return err
	}

	for {
		checkFollowingStatus, err := w.checkFollowingStatus(chainId)
		if err != nil {
			return err
		}
		if !checkFollowingStatus {
			logrus.Warningf("chain %s is not follow status", chainId)
			return nil
		}

		txList, err := w.getTxList(chainId, taskRecord.Height, int64(limit))
		if err != nil {
			return err
		}

		if len(txList) == 0 {
			return nil
		}

		ibcTxList, ibcDenomList := w.handleSourceTx(chainId, txList, denomMap)
		if len(ibcDenomList) > 0 {
			if err = denomRepo.InsertBatch(ibcDenomList); err != nil {
				logrus.Errorf("task %s worker %s denomRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chainId, err)
				return err
			}
		}
		if len(ibcTxList) > 0 {
			if err = ibcTxRepo.InsertBatch(ibcTxList); err != nil {
				logrus.Errorf("task %s worker %s ibcTxRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chainId, err)
				return err
			}
		}

		taskRecord.Height = txList[len(txList)-1].Height
		if err = taskRecordRepo.UpdateHeight(taskRecord.TaskName, taskRecord.Height); err != nil {
			logrus.Errorf("task %s worker %s taskRecordRepo.UpdateHeight %s error, %v", w.taskName, w.workerName, chainId, err)
			return err
		}

		totalParseTx += len(txList)
		if totalParseTx >= global.Config.Task.SingleChainSyncTransferTxMax {
			break
		}
	}

	return nil
}

func (w *syncTransferTxWorker) handleSourceTx(chainId string, txList []*entity.Tx, denomMap map[string]*entity.IBCDenom) ([]*entity.ExIbcTx, entity.IBCDenomList) {
	var ibcTxList []*entity.ExIbcTx
	var ibcDenomList entity.IBCDenomList
	for _, tx := range txList {
		for msgIndex, msg := range tx.DocTxMsgs {
			if msg.Type != constant.MsgTypeTransfer {
				continue
			}

			var ibcTxStatus entity.IbcTxStatus
			switch tx.Status {
			case entity.TxStatusSuccess:
				ibcTxStatus = entity.IbcTxStatusProcessing
			case entity.TxStatusFailed:
				ibcTxStatus = entity.IbcTxStatusFailed
			}

			transferTxMsg := msg.TransferMsg()
			scPort := transferTxMsg.SourcePort
			scChannel := transferTxMsg.SourceChannel
			scDenom := transferTxMsg.Token.Denom
			dcChainId, dcPort, dcChannel := w.matchDcInfo(chainId, scPort, scChannel)

			var fullDenomPath, sequence string
			ibcDenom, isExisted := denomMap[scDenom]
			if ibcTxStatus == entity.IbcTxStatusFailed { // get base_denom info from ibc_denom collection

			} else {
				dcPort, dcChannel, fullDenomPath, sequence = w.getIbcInfoFromEventsMsg(msgIndex, tx)
				if !isExisted { // denom 不存在
					ibcDenom = w.traceDenom(scDenom, fullDenomPath, chainId)
				}
			}

			if dcChainId == "" && ibcTxStatus != entity.IbcTxStatusFailed {
				ibcTxStatus = entity.IbcTxStatusSetting
			}

			if ibcTxStatus == entity.IbcTxStatusProcessing && !isExisted && ibcDenom != nil {
				ibcDenomList = append(ibcDenomList, ibcDenom)
				denomMap[ibcDenom.Denom] = ibcDenom
			}

			var baseDemom, baseDenomChainId string
			if ibcDenom != nil {
				baseDemom = ibcDenom.BaseDenom
				baseDenomChainId = ibcDenom.BaseDenomChainId
			}
			recordId := fmt.Sprintf("%s%s%s%s%s%s%s%d", scPort, scChannel, dcPort, dcChannel, sequence, chainId, tx.TxHash, msgIndex)
			ibcTxList = append(ibcTxList, &entity.ExIbcTx{
				RecordId:  recordId,
				ScAddr:    transferTxMsg.Sender,
				DcAddr:    transferTxMsg.Receiver,
				ScPort:    scPort,
				ScChannel: scChannel,
				ScChainId: chainId,
				DcPort:    dcPort,
				DcChannel: dcChannel,
				DcChainId: dcChainId,
				Sequence:  sequence,
				Status:    ibcTxStatus,
				ScTxInfo: &entity.TxInfo{
					Hash:      tx.TxHash,
					Status:    tx.Status,
					Time:      tx.Time,
					Height:    tx.Height,
					Fee:       tx.Fee,
					MsgAmount: transferTxMsg.Token,
					Msg:       msg,
				},
				DcTxInfo:       nil,
				RefundedTxInfo: nil,
				Log: &entity.Log{
					ScLog: tx.Log,
				},
				Denoms: &entity.Denoms{
					ScDenom: scDenom,
					DcDenom: "",
				},
				BaseDenom:        baseDemom,
				BaseDenomChainId: baseDenomChainId,
				TxTime:           tx.Time,
				CreateAt:         time.Now().Unix(),
				UpdateAt:         time.Now().Unix(),
			})
		}
	}

	return ibcTxList, ibcDenomList
}

// traceDenom 通过追踪denom的溯源路径，解析出当前denom的相关信息
//   - fullDenomPath denom完全路径，例："transfer/channel-1/uiris"
func (w *syncTransferTxWorker) traceDenom(denom, fullDenomPath, chainId string) *entity.IBCDenom {
	unix := time.Now().Unix()
	if !strings.HasPrefix(denom, constant.IBCTokenPreFix) { // base denom
		return &entity.IBCDenom{
			ChainId:          chainId,
			Denom:            denom,
			PrevDenom:        "",
			PrevChainId:      "",
			BaseDenom:        denom,
			BaseDenomChainId: chainId,
			DenomPath:        "",
			IsSourceChain:    false,
			IsBaseDenom:      true,
			CreateAt:         unix,
			UpdateAt:         unix,
		}
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("trace denom: %s, chain: %s, full path: %s, error. %v ", denom, chainId, fullDenomPath, err)
		}
	}()

	var currentChainId string
	var isBaseDenom bool
	currentChainId = chainId
	pathSplits := strings.Split(fullDenomPath, "/")
	denomPath := strings.Join(pathSplits[0:len(pathSplits)-1], "/")
	var TraceDenomList []*dto.DenomSimpleDTO
	TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
		Denom:   denom,
		ChainId: chainId,
	})

	for {
		if len(pathSplits) <= 1 {
			break
		}

		currentPort, currentChannel := pathSplits[0], pathSplits[1]
		tempPrevChainId, tempPrevPort, tempPrevChannel := w.matchDcInfo(currentChainId, currentPort, currentChannel)
		if tempPrevChainId == "" { // 无法向前溯源了
			break
		} else {
			TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
				Denom:   w.calculateIbcHash(strings.Join(pathSplits[2:], "/")),
				ChainId: tempPrevChainId,
			})
		}

		currentChainId, currentPort, currentChannel = tempPrevChainId, tempPrevPort, tempPrevChannel
		pathSplits = pathSplits[2:]
	}

	var prevDenom, prevChainId, baseDenom, baseDenomChainId string
	if len(TraceDenomList) == 1 { // denom 本身就是base denom
		isBaseDenom = true
		baseDenom = denom
		baseDenomChainId = chainId
	} else {
		isBaseDenom = false
		prevDenom = TraceDenomList[1].Denom
		prevChainId = TraceDenomList[1].ChainId
		baseDenom = TraceDenomList[len(TraceDenomList)-1].Denom
		baseDenomChainId = TraceDenomList[len(TraceDenomList)-1].ChainId
	}

	return &entity.IBCDenom{
		ChainId:          chainId,
		Denom:            denom,
		PrevDenom:        prevDenom,
		PrevChainId:      prevChainId,
		BaseDenom:        baseDenom,
		BaseDenomChainId: baseDenomChainId,
		DenomPath:        denomPath,
		IsSourceChain:    false,
		IsBaseDenom:      isBaseDenom,
		CreateAt:         unix,
		UpdateAt:         unix,
	}
}

func (w *syncTransferTxWorker) calculateIbcHash(fullPath string) string {
	if len(strings.Split(fullPath, "/")) == 1 {
		return fullPath
	}

	hash := utils.Sha256(fullPath)
	return fmt.Sprintf("%s/%s", constant.IBCTokenPreFix, strings.ToUpper(hash))
}

func (w *syncTransferTxWorker) getIbcInfoFromEventsMsg(msgIndex int, tx *entity.Tx) (dcPort, dcChannel, denomPath, sequence string) {
	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if evt.Type == "send_packet" {
				for _, attr := range evt.Attributes {
					switch attr.Key {
					case "packet_dst_port":
						dcPort = attr.Value
					case "packet_dst_channel":
						dcChannel = attr.Value
					case "packet_sequence":
						sequence = attr.Value
					case "packet_data":
						var data model.TransferTxPacketData
						_ = json.Unmarshal([]byte(attr.Value), &data)
						denomPath = data.Denom
					default:
					}
				}
			}
		}
	}

	return
}

func (w *syncTransferTxWorker) matchDcInfo(scChainId, scPort, scChannel string) (dcChainId, dcPort, dcChannel string) {
	if coordinator == nil || coordinator.allChainMap == nil || coordinator.allChainMap[scChainId] == nil {
		return
	}

	for _, ibcInfo := range coordinator.allChainMap[scChainId].IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.PortId == scPort && path.ChannelId == scChannel {
				dcChainId = path.ChainId
				dcPort = path.Counterparty.PortId
				dcChannel = path.Counterparty.ChannelId
				return
			}
		}
	}

	return
}

func (w *syncTransferTxWorker) checkFollowingStatus(chainId string) (bool, error) {
	status, err := syncTaskRepo.CheckFollowingStatus(chainId)
	if err != nil {
		logrus.Errorf("task %s worker %s checkFollowingStatus %s error, %v", w.taskName, w.workerName, chainId, err)
		return false, err
	}

	return status, nil
}

// checkTaskRecord 检查task_record的状态，如果不存在task_record 记录，则新增
func (w *syncTransferTxWorker) checkTaskRecord(chainId string) (*entity.IbcTaskRecord, error) {
	taskRecord, err := taskRecordRepo.FindByTaskName(fmt.Sprintf(entity.TaskNameFmt, chainId))
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s worker %s checkTaskRecord %s error, %v", w.taskName, w.workerName, chainId, err)
			return nil, err
		}

		taskRecord = &entity.IbcTaskRecord{
			TaskName: fmt.Sprintf(entity.TaskNameFmt, chainId),
			Height:   0,
			Status:   entity.TaskRecordStatusOpen,
			CreateAt: time.Now().Unix(),
			UpdateAt: time.Now().Unix(),
		}

		if err := taskRecordRepo.Insert(taskRecord); err != nil {
			logrus.Errorf("task %s worker %s checkTaskRecord %s error, %v", w.taskName, w.workerName, chainId, err)
			return nil, err
		} else {
			return taskRecord, nil
		}
	}

	return taskRecord, nil
}

func (w *syncTransferTxWorker) getChainDenomMap(chainId string) (map[string]*entity.IBCDenom, error) {
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

func (w *syncTransferTxWorker) getTxList(chainId string, height, limit int64) ([]*entity.Tx, error) {
	transferTxList, err := txRepo.GetTransferTx(chainId, height, limit)
	if err != nil {
		logrus.Errorf("task %s worker %s GetTransferTx %s error, %v", w.taskName, w.workerName, chainId, err)
		return nil, err
	}

	if len(transferTxList) < int(limit) {
		return transferTxList, nil
	}

	txHashMap := make(map[string]string)
	for _, v := range transferTxList {
		txHashMap[v.TxHash] = ""
	}

	heightTxList, err := txRepo.FindByTypeAndHeight(chainId, constant.MsgTypeTransfer, height)
	if err != nil {
		logrus.Errorf("task %s worker %s FindByTypeAndHeight %s error, %v", w.taskName, w.workerName, chainId, err)
		return nil, err
	}

	for _, v := range heightTxList {
		_, ok := txHashMap[v.TxHash]
		if !ok {
			transferTxList = append(transferTxList, v)
		}
	}

	return transferTxList, nil
}
