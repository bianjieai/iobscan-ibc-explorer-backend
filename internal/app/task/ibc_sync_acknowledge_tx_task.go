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

	syncLatestAcknowledge := func() error {
		txs, err := ibcTxRepo.GetNeedAcknowledgeTxs(false)
		if err != nil {
			return err
		}
		for _, val := range txs {
			err := t.SaveAcknowledgeTx(val, false)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				logrus.Warnf("task %s SaveAcknowledgeTx failed %s, chain_id:%s packet_id:%s",
					t.Name(),
					err.Error(),
					val.ScChainId,
					val.ScTxInfo.Msg.CommonMsg().PacketId)
			}
		}
		return nil
	}

	syncLatestRecvPacket := func() error {
		txs, err := ibcTxRepo.GetNeedRecvPacketTxs(false)
		if err != nil {
			return err
		}
		for _, val := range txs {
			err := SaveRecvPacketTx(val, false)
			if err != nil && err != qmgo.ErrNoSuchDocuments {
				logrus.Warnf("task %s SaveRecvPacketTx failed %s, chain_id:%s packet_id:%s",
					t.Name(),
					err.Error(),
					val.ScChainId,
					val.ScTxInfo.Msg.CommonMsg().PacketId)
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := syncLatestAcknowledge()
		logrus.Infof("task %s fix Acknowledge latest end, %v", t.Name(), err)
	}()

	go func() {
		defer wg.Done()
		err := syncLatestRecvPacket()
		logrus.Infof("task %s fix RecvPacket latest end, %v", t.Name(), err)
	}()

	wg.Wait()
	return 1
}

func (t *IbcSyncAcknowledgeTxTask) SaveAcknowledgeTx(ibcTx *entity.ExIbcTx, history bool) error {
	ackTxs, err := txRepo.GetAcknowledgeTxs(ibcTx.ScChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)
	if err != nil {
		return err
	}
	if len(ackTxs) > 0 {
		//"成功"状态IBC，第三段取最新的ack tx交易
		ackTx := ackTxs[0]
		ibcTx.RefundedTxInfo = &entity.TxInfo{
			Hash:      ackTx.TxHash,
			Height:    ackTx.Height,
			Time:      ackTx.Time,
			Status:    ackTx.Status,
			Fee:       ackTx.Fee,
			Memo:      ackTx.Memo,
			Signers:   ackTx.Signers,
			MsgAmount: nil,
			Msg:       getMsgByType(*ackTx, constant.MsgTypeAcknowledgement),
		}
		return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, bson.M{
			"$set": bson.M{
				"refunded_tx_info": ibcTx.RefundedTxInfo,
			},
		})
	}
	return nil
}

func getMsgByType(tx entity.Tx, msgType string) *model.TxMsg {
	for _, msg := range tx.DocTxMsgs {
		if msg.Type == msgType {
			return msg
		}
	}
	return nil
}

func SaveRecvPacketTx(ibcTx *entity.ExIbcTx, history bool) error {
	recvTxs, err := txRepo.GetRecvPacketTxs(ibcTx.DcChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)
	if err != nil {
		return err
	}
	var recvTx *entity.Tx
	for _, val := range recvTxs {
		if val.Status == entity.TxStatusSuccess {
			recvTx = val
			if ibcTx.DcConnectionId == "" {
				for index, msg := range val.DocTxMsgs {
					if msg.Type == constant.MsgTypeRecvPacket {
						ibcTx.DcConnectionId = getConnectByRecvPacketEventsNews(val.EventsNew, index)
					}
				}
			}
			break
		}

	}
	//没有匹配成功，取最新recv_packet
	if recvTx == nil && len(recvTxs) > 0 {
		recvTx = recvTxs[0]
	}
	if recvTx != nil {
		ibcTx.DcTxInfo = &entity.TxInfo{
			Hash:      recvTx.TxHash,
			Height:    recvTx.Height,
			Time:      recvTx.Time,
			Status:    recvTx.Status,
			Fee:       recvTx.Fee,
			Memo:      recvTx.Memo,
			Signers:   recvTx.Signers,
			ErrLog:    recvTx.Log,
			MsgAmount: nil,
			Msg:       getMsgByType(*recvTx, constant.MsgTypeRecvPacket),
		}
		return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, bson.M{
			"$set": bson.M{
				"dc_tx_info":       ibcTx.DcTxInfo,
				"dc_connection_id": ibcTx.DcConnectionId,
			},
		})
	}
	return nil
}

func getConnectByRecvPacketEventsNews(eventNews []entity.EventNew, msgIndex int) string {
	var connect string
	for _, item := range eventNews {
		if item.MsgIndex == uint32(msgIndex) {
			for _, val := range item.Events {
				if val.Type == "write_acknowledgement" || val.Type == "recv_packet" {
					for _, attribute := range val.Attributes {
						switch attribute.Key {
						case "packet_connection":
							connect = attribute.Value
							//case "packet_ack":
							//	ackData = attribute.Value
						}
					}
				}
			}
		}
	}
	return connect
}
