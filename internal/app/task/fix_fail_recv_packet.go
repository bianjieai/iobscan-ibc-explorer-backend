package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
	segments, err := getSegment()
	if err != nil {
		logrus.Errorf("task %s getSegment error, %v", t.Name(), err)
		return -1
	}

	historySegments, err := getHistorySegment()
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

	for _, v := range segments {
		logrus.Infof("task %s fix %s %d-%d", t.Name(), target, v.StartTime, v.EndTime)
		var skip int64 = 0
		for {
			txs, err := ibcTxRepo.FindRecvPacketTxsEmptyTxs(v.StartTime, v.EndTime, skip, limit, isTargetHistory)
			if err != nil {
				logrus.Errorf("task %s FindRecvPacketTxsEmptyTxs %s %d-%d err, %v", t.Name(), target, v.StartTime, v.EndTime, err)
				return err
			}

			for _, val := range txs {
				err := SaveFailRecvPacketTx(val, isTargetHistory)
				if err != nil && err != qmgo.ErrNoSuchDocuments {
					logrus.Errorf("task %s SaveFailRecvPacketTx %s err, chain_id: %s, packet_id: %s, %v", t.Name(), target, val.ScChainId, val.ScTxInfo.Msg.CommonMsg().PacketId, err)
					return err
				}
			}

			if int64(len(txs)) < limit {
				break
			}
			skip += limit
		}
	}
	return nil
}

func SaveFailRecvPacketTx(ibcTx *entity.ExIbcTx, history bool) error {
	recvTxs, err := txRepo.GetRecvPacketTxs(ibcTx.DcChainId, ibcTx.ScTxInfo.Msg.CommonMsg().PacketId)
	if err != nil {
		return err
	}
	var recvTx *entity.Tx
	for _, val := range recvTxs {
		if val.Status == entity.TxStatusSuccess {
			recvTx = val
			for index, msg := range val.DocTxMsgs {
				if msg.Type == constant.MsgTypeRecvPacket {
					ibcTx.DcConnectionId = getConnectByRecvPacketEventsNews(val.EventsNew, index)
				}
			}
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
			MsgAmount: nil,
			Msg:       getMsgByType(*recvTx, constant.MsgTypeRecvPacket),
		}
		return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, bson.M{
			"$set": bson.M{
				"dc_tx_info": ibcTx.DcTxInfo,
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
