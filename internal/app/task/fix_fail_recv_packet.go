package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
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

	syncFailRecvPacket := func(history bool) error {
		txs, err := ibcTxRepo.GetNeedFailRecvPacketTxs(history)
		if err != nil {
			return err
		}
		for _, val := range txs {
			err := t.SaveFailRecvPacketTx(val, history)
			if err != nil {
				logrus.Warn("SaveFailRecvPacketTx failed, "+err.Error(),
					" chain_id: ", val.ScChainId,
					" packet_id: ", val.ScTxInfo.Msg.CommonMsg().PacketId)
			}
		}
		return nil
	}

	if err := syncFailRecvPacket(false); err != nil {
		logrus.Error(err.Error())
		return -1
	}
	if err := syncFailRecvPacket(true); err != nil {
		logrus.Error(err.Error())
		return -1
	}
	return 1
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
		if len(val.Signers) > 0 {
			_, ok := dcAddrMap[val.Signers[0]]
			if ok {
				recvTx = val
				break
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
			Log:       recvTx.Log,
			MsgAmount: nil,
			Msg:       getMsgByType(*recvTx, constant.MsgTypeRecvPacket),
		}
		return ibcTxRepo.UpdateOne(ibcTx.RecordId, history, ibcTx)
	}
	return nil
}
