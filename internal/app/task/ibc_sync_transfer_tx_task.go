package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/ibctool"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IbcSyncTransferTxTask struct {
}

var _ Task = new(IbcSyncTransferTxTask)
var transferTxCoordinator *stringQueueCoordinator

func (t *IbcSyncTransferTxTask) Name() string {
	return "ibc_sync_transfer_tx_task"
}

func (t *IbcSyncTransferTxTask) Cron() int {
	if taskConf.CronTimeSyncTransferTxTask > 0 {
		return taskConf.CronTimeSyncTransferTxTask
	}
	return ThreeMinute
}

func (t *IbcSyncTransferTxTask) workerNum() int {
	if global.Config.Task.SyncTransferTxWorkerNum > 0 {
		return global.Config.Task.SyncTransferTxWorkerNum
	}
	return syncTransferTxTaskWorkerNum
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
		chainQueue.Push(v.ChainName)
	}
	transferTxCoordinator = &stringQueueCoordinator{
		stringQueue: chainQueue,
	}

	workerNum := t.workerNum()
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	for i := 1; i <= workerNum; i++ {
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
		chain, err := transferTxCoordinator.getOne()
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
		if err = w.parseChainIbcTx(chain); err != nil {
			logrus.Errorf("task %s worker %s parse chain %s tx error,time use: %d(s), %v", w.taskName, w.workerName, chain, time.Now().Unix()-startTime, err)
		} else {
			logrus.Infof("task %s worker %s parse chain %s tx end,time use: %d(s)", w.taskName, w.workerName, chain, time.Now().Unix()-startTime)
		}
	}
}

func (w *syncTransferTxWorker) parseChainIbcTx(chain string) error {
	totalParseTx := 0
	//const limit = 500
	maxParseTx := global.Config.Task.SingleChainSyncTransferTxMax
	if maxParseTx <= 0 {
		maxParseTx = defaultMaxHandlerTx
	}

	taskRecord, err := w.checkTaskRecord(chain)
	if err != nil {
		return err
	}
	if taskRecord.Status == entity.TaskRecordStatusClose {
		return nil
	}

	denomMap, err := w.getChainDenomMap(chain)
	if err != nil {
		return err
	}

	for {
		checkFollowingStatus, err := w.checkFollowingStatus(chain)
		if err != nil {
			return err
		}
		if !checkFollowingStatus {
			logrus.Warningf("chain %s is not follow status", chain)
			return nil
		}

		txList, err := w.getTxList(chain, taskRecord.Height, int64(constant.DefaultLimit))
		if err != nil {
			return err
		}

		if len(txList) == 0 {
			return nil
		}

		ibcTxList, ibcDenomList := w.handleSourceTx(chain, txList, denomMap)
		if len(ibcDenomList) > 0 {
			if err = denomRepo.InsertBatch(ibcDenomList); err != nil {
				logrus.Errorf("task %s worker %s denomRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chain, err)
				return err
			}
		}
		if len(ibcTxList) > 0 {
			if err = ibcTxRepo.InsertBatch(ibcTxList); err != nil {
				logrus.Errorf("task %s worker %s ibcTxRepo.InsertBatch %s error, %v", w.taskName, w.workerName, chain, err)
				return err
			}
		}

		taskRecord.Height = txList[len(txList)-1].Height
		if err = taskRecordRepo.UpdateHeight(taskRecord.TaskName, taskRecord.Height); err != nil {
			logrus.Errorf("task %s worker %s taskRecordRepo.UpdateHeight %s error, %v", w.taskName, w.workerName, chain, err)
			return err
		}

		totalParseTx += len(txList)
		if len(txList) < constant.DefaultLimit || totalParseTx >= maxParseTx {
			break
		}
	}

	return nil
}

func (w *syncTransferTxWorker) handleSourceTx(chain string, txList []*entity.Tx, denomMap map[string]*entity.IBCDenom) ([]*entity.ExIbcTx, entity.IBCDenomList) {
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
			dcChain, dcPort, dcChannel := ibctool.MatchDcInfo(chain, scPort, scChannel, w.chainMap)

			var fullDenomPath, sequence, scConnection string
			ibcDenom, isExisted := denomMap[scDenom]
			if ibcTxStatus != entity.IbcTxStatusFailed {
				dcPort, dcChannel, fullDenomPath, sequence, scConnection = parseTransferTxEvents(msgIndex, tx)
				if !isExisted { // denom 不存在
					ibcDenom = ibctool.TraceDenom(fullDenomPath, chain, w.chainMap)
				}
			}

			if dcChain == "" && ibcTxStatus != entity.IbcTxStatusFailed {
				ibcTxStatus = entity.IbcTxStatusSetting
			}

			if ibcTxStatus == entity.IbcTxStatusProcessing && !isExisted && ibcDenom != nil {
				ibcDenomList = append(ibcDenomList, ibcDenom)
				denomMap[ibcDenom.Denom] = ibcDenom
			}

			var baseDemom, baseDenomChain string
			if ibcDenom != nil {
				baseDemom = ibcDenom.BaseDenom
				baseDenomChain = ibcDenom.BaseDenomChain
			}
			recordIdStr := fmt.Sprintf("%s%s%s%s%s%s%s%d", scPort, scChannel, dcPort, dcChannel, sequence, chain, tx.TxHash, msgIndex)
			recordId := utils.Md5(recordIdStr)
			nowUnix := time.Now().Unix()

			exIbcTx := &entity.ExIbcTx{
				Id:             primitive.NewObjectID(),
				RecordId:       recordId,
				TxTime:         tx.Time,
				ScAddr:         transferTxMsg.Sender,
				DcAddr:         transferTxMsg.Receiver,
				ScPort:         scPort,
				ScChannel:      scChannel,
				ScConnectionId: scConnection,
				ScClientId:     "",
				ScChain:        chain,
				DcPort:         dcPort,
				DcChannel:      dcChannel,
				DcConnectionId: "",
				DcClientId:     "",
				DcChain:        dcChain,
				Sequence:       sequence,
				Status:         ibcTxStatus,
				ScTxInfo: &entity.TxInfo{
					Hash:      tx.TxHash,
					Status:    tx.Status,
					Time:      tx.Time,
					Height:    tx.Height,
					Fee:       tx.Fee,
					MsgAmount: transferTxMsg.Token,
					Msg:       msg,
					Memo:      tx.Memo,
					Signers:   tx.Signers,
					Log:       tx.Log,
				},
				DcTxInfo:         nil,
				AckTimeoutTxInfo: nil,
				Denoms: &entity.Denoms{
					ScDenom: scDenom,
					DcDenom: "",
				},
				BaseDenom:      baseDemom,
				BaseDenomChain: baseDenomChain,
				RetryTimes:     0,
				NextTryTime:    nowUnix,
				CreateAt:       nowUnix,
				UpdateAt:       nowUnix,
			}
			w.setClientId(exIbcTx) // set ScClientId, DcClientId
			ibcTxList = append(ibcTxList, exIbcTx)
		}
	}
	return ibcTxList, ibcDenomList
}

func (w *syncTransferTxWorker) setClientId(ibcTx *entity.ExIbcTx) {
	chainConf, ok := w.chainMap[ibcTx.ScChain]
	if ok {
		client := chainConf.GetChannelClient(ibcTx.ScPort, ibcTx.ScChannel)
		ibcTx.ScClientId = client
	}

	if ibcTx.DcChain != "" {
		chainConf, ok = w.chainMap[ibcTx.DcChain]
		if ok {
			client := chainConf.GetChannelClient(ibcTx.DcPort, ibcTx.DcChannel)
			ibcTx.DcClientId = client
		}
	}
}

func (w *syncTransferTxWorker) checkFollowingStatus(chain string) (bool, error) {
	status, err := syncTaskRepo.CheckFollowingStatus(chain)
	if err != nil {
		logrus.Errorf("task %s worker %s checkFollowingStatus %s error, %v", w.taskName, w.workerName, chain, err)
		return false, err
	}

	return status, nil
}

// checkTaskRecord 检查task_record的状态，如果不存在task_record 记录，则新增
func (w *syncTransferTxWorker) checkTaskRecord(chain string) (*entity.IbcTaskRecord, error) {
	taskRecord, err := taskRecordRepo.FindByTaskName(fmt.Sprintf(entity.TaskNameFmt, chain))
	if err != nil {
		if err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s worker %s checkTaskRecord %s error, %v", w.taskName, w.workerName, chain, err)
			return nil, err
		}

		taskRecord = &entity.IbcTaskRecord{
			TaskName: fmt.Sprintf(entity.TaskNameFmt, chain),
			Height:   0,
			Status:   entity.TaskRecordStatusOpen,
			CreateAt: time.Now().Unix(),
			UpdateAt: time.Now().Unix(),
		}

		if err := taskRecordRepo.Insert(taskRecord); err != nil {
			logrus.Errorf("task %s worker %s checkTaskRecord %s error, %v", w.taskName, w.workerName, chain, err)
			return nil, err
		} else {
			return taskRecord, nil
		}
	}

	return taskRecord, nil
}

func (w *syncTransferTxWorker) getChainDenomMap(chain string) (map[string]*entity.IBCDenom, error) {
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

func (w *syncTransferTxWorker) getTxList(chain string, height, limit int64) ([]*entity.Tx, error) {
	transferTxList, err := txRepo.GetTransferTx(chain, height, limit)
	if err != nil {
		logrus.Errorf("task %s worker %s GetTransferTx %s error, %v", w.taskName, w.workerName, chain, err)
		return nil, err
	}

	if len(transferTxList) < int(limit) {
		return transferTxList, nil
	}

	maxHeight := transferTxList[len(transferTxList)-1].Height
	txHashMap := make(map[string]string)
	for _, v := range transferTxList {
		if v.Height == maxHeight {
			txHashMap[v.TxHash] = ""
		}
	}

	heightTxList, err := txRepo.FindByTypeAndHeight(chain, constant.MsgTypeTransfer, maxHeight)
	if err != nil {
		logrus.Errorf("task %s worker %s FindByTypeAndHeight %s error, %v", w.taskName, w.workerName, chain, err)
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
