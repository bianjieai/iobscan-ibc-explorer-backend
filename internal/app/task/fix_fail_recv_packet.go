package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
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
				err := t.SaveFailRecvPacketTx(val, isTargetHistory)
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

func (t *FixFailRecvPacketTask) SaveFailRecvPacketTx(ibcTx *entity.ExIbcTx, history bool) error {
	relayers, err := relayerRepo.FindRelayer(ibcTx.ScChainId, ibcTx.RefundedTxInfo.Msg.CommonMsg().Signer, ibcTx.ScChannel)
	if err != nil {
		return err
	}
	dcAddrMap := make(map[string]struct{}, len(relayers))
	for _, val := range relayers {
		if val.ChainAAddress == ibcTx.RefundedTxInfo.Msg.CommonMsg().Signer && val.ChainBAddress != "" {
			dcAddrMap[val.ChainBAddress] = struct{}{}
		} else if val.ChainBAddress == ibcTx.RefundedTxInfo.Msg.CommonMsg().Signer && val.ChainAAddress != "" {
			dcAddrMap[val.ChainAAddress] = struct{}{}
		}
	}
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
					ibcTx.DcConnectionId = t.getConnectByRecvPacketEventsNews(val.EventsNew, index)
				}
			}
		} else {
			//失败的recv_packet交易
			if len(val.Signers) > 0 {
				_, ok := dcAddrMap[val.Signers[0]]
				if ok {
					recvTx = val
					break
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
		return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, ibcTx)
	}
	return nil
}

func (t *FixFailRecvPacketTask) getConnectByRecvPacketEventsNews(eventNews []entity.EventNew, msgIndex int) string {
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
