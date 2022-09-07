package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
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

	if err := syncAcknowledge(false); err != nil {
		logrus.Error(err.Error())
		return -1
	}
	if err := syncAcknowledge(true); err != nil {
		logrus.Error(err.Error())
		return -1
	}
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
	return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, ibcTx)
}

func getMsgByType(tx entity.Tx, msgType string) *model.TxMsg {
	for _, msg := range tx.DocTxMsgs {
		if msg.Type == msgType {
			return msg
		}
	}
	return nil
}
