package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

type IbcSyncAcknowledgeTxTask struct {
}

var _ Task = new(IbcSyncAcknowledgeTxTask)

func (t *IbcSyncAcknowledgeTxTask) Name() string {
	return "ibc_sync_acknowledge_tx_task"
}

func (t *IbcSyncAcknowledgeTxTask) Cron() int {
	if taskConf.CronTimeSyncAckTxTask > 0 {
		return taskConf.CronTimeSyncAckTxTask
	}
	return ThreeMinute
}

func (t *IbcSyncAcknowledgeTxTask) Run() int {

	syncAcknowledge := func(history bool) error {
		txs, err := ibcTxRepo.GetNeedAcknowledgeTxs(history)
		if err != nil {
			return err
		}
		for _, val := range txs {
			err := t.SaveAcknowledgeTx(val, history)
			if err != nil {
				logrus.Warn("SaveAcknowledgeTx failed, "+err.Error(),
					"chain_id:", val.ScChainId,
					"packet_id:", val.ScTxInfo.Msg.CommonMsg().PacketId)
			}
		}
		return nil
	}

	syncRecvPacket := func(history bool) error {
		txs, err := ibcTxRepo.GetNeedFailRecvPacketTxs(history)
		if err != nil {
			return err
		}
		for _, val := range txs {
			err := SaveFailRecvPacketTx(val, history)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				logrus.Warn("SaveFailRecvPacketTx failed, "+err.Error(),
					"chain_id:", val.ScChainId,
					"packet_id:", val.ScTxInfo.Msg.CommonMsg().PacketId)
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := syncAcknowledge(false)
		logrus.Infof("task %s fix Acknowledge latest end, %v", t.Name(), err)
	}()

	go func() {
		defer wg.Done()
		err := syncRecvPacket(false)
		logrus.Infof("task %s fix RecvPacket latest end, %v", t.Name(), err)
	}()

	wg.Wait()
	return 1
}

func (t *IbcSyncAcknowledgeTxTask) SaveAcknowledgeTx(ibcTx *entity.ExIbcTx, history bool) error {
	ackTx, err := txRepo.GetAcknowledgeTxs(ibcTx.ScChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)
	if err != nil {
		return err
	}
	ibcTx.RefundedTxInfo = &entity.TxInfo{
		Hash:      ackTx.TxHash,
		Height:    ackTx.Height,
		Time:      ackTx.Time,
		Status:    ackTx.Status,
		Fee:       ackTx.Fee,
		Memo:      ackTx.Memo,
		Signers:   ackTx.Signers,
		MsgAmount: nil,
		Msg:       getMsgByType(ackTx, constant.MsgTypeAcknowledgement),
	}
	return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, bson.M{
		"$set": bson.M{
			"refunded_tx_info": ibcTx.RefundedTxInfo,
		},
	})
}

func getMsgByType(tx entity.Tx, msgType string) *model.TxMsg {
	for _, msg := range tx.DocTxMsgs {
		if msg.Type == msgType {
			return msg
		}
	}
	return nil
}
