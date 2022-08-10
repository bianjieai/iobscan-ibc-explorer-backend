package task

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type IbcSyncTransferTxTask struct {
}

var _ Task = new(IbcSyncTransferTxTask)
var transferTxCoordinator *chainQueueCoordinator

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
	transferTxCoordinator = &chainQueueCoordinator{
		chainQueue: chainQueue,
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(syncTransferTxTaskWorkerQuantity)
	for i := 1; i <= syncTransferTxTaskWorkerQuantity; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			newSyncTransferTxWorker(t.Name(), wn, chainMap).exec()
			waitGroup.Done()
		}(workName)
	}
	waitGroup.Wait()

	return 1
}

// =========================================================================
// =========================================================================
// worker

func newSyncTransferTxWorker(taskName, workerName string, chainMap map[string]*entity.ChainConfig) *syncTransferTxWorker {
	return &syncTransferTxWorker{
		taskName:   taskName,
		workerName: workerName,
		chainMap:   chainMap,
	}
}

type syncTransferTxWorker struct {
	taskName   string
	workerName string
	chainMap   map[string]*entity.ChainConfig
}

func (w *syncTransferTxWorker) exec() {
	logrus.Infof("task %s worker %s start", w.taskName, w.workerName)
	for {
		chainId, err := transferTxCoordinator.getChain()
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

func (w *syncTransferTxWorker) parseChainIbcTx(chainId string) error {
	totalParseTx := 0
	const limit = 500
	maxParseTx := global.Config.Task.SingleChainSyncTransferTxMax
	if maxParseTx <= 0 {
		maxParseTx = defaultMaxHandlerTx
	}

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
		if len(txList) < limit || totalParseTx >= maxParseTx {
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
			dcChainId, dcPort, dcChannel := matchDcInfo(chainId, scPort, scChannel, w.chainMap)

			var fullDenomPath, sequence string
			ibcDenom, isExisted := denomMap[scDenom]
			if ibcTxStatus == entity.IbcTxStatusFailed { // get base_denom info from ibc_denom collection

			} else {
				dcPort, dcChannel, fullDenomPath, sequence = w.getIbcInfoFromEventsMsg(msgIndex, tx)
				if !isExisted { // denom 不存在
					ibcDenom = traceDenom(fullDenomPath, chainId, w.chainMap)
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
			recordIdStr := fmt.Sprintf("%s%s%s%s%s%s%s%d", scPort, scChannel, dcPort, dcChannel, sequence, chainId, tx.TxHash, msgIndex)
			recordId := utils.Md5(recordIdStr)
			nowUnix := time.Now().Unix()
			ibcTxList = append(ibcTxList, &entity.ExIbcTx{
				RecordId:  recordId,
				TxTime:    tx.Time,
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
				RetryTimes:       0,
				NextTryTime:      nowUnix,
				CreateAt:         nowUnix,
				UpdateAt:         nowUnix,
			})
		}
	}

	return ibcTxList, ibcDenomList
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
