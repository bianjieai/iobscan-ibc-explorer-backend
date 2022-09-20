package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"strings"
)

var _ OneOffTask = new(FixAckTxPacketIdTask)

type FixAckTxPacketIdTask struct {
}

func (f FixAckTxPacketIdTask) Name() string {
	return "fix_ack_tx_packet_id"
}

func (f FixAckTxPacketIdTask) Switch() bool {
	return global.Config.Task.SwitchFixAckTxPacketIdTask
}

func (f FixAckTxPacketIdTask) RunWithParam(chainsStr string) int {
	return f.handle(chainsStr)
}
func (f FixAckTxPacketIdTask) Run() int {
	return f.handle(global.Config.ChainConfig.FixAckTxPacketIdChains)
}

func (f FixAckTxPacketIdTask) handle(chainsStr string) int {
	newChainIds := strings.Split(chainsStr, ",")
	if len(newChainIds) == 0 {
		logrus.Errorf("task %s don't have fix ack packet_id chains", f.Name())
		return 1
	}

	for _, val := range newChainIds {
		minHTx, err := txRepo.FindHeight(val, true)
		if err != nil {
			logrus.Errorf("find minHeight err chain_id:%s err:%v", val, err.Error())
			return -1
		}
		maxETx, err := txRepo.FindHeight(val, false)
		if err != nil {
			logrus.Errorf("skip fix ack_tx packet_id chain_id:%s err:%v", val, err.Error())
			return -1
		}
		logrus.Debugf("start fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
			minHTx.Height-1, maxETx.Height, val)
		ret := NewfixAckTxTask(val, minHTx.Height-1, maxETx.Height).Run()
		if ret > 0 {
			logrus.Infof("finish fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
				minHTx.Height-1, maxETx.Height, val)
		} else {
			logrus.Errorf("fail fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
				minHTx.Height-1, maxETx.Height, val)
		}

	}
	return 1
}

type fixAckTxTask struct {
	StartHeight int64
	EndHeight   int64
	ChainId     string
}

func NewfixAckTxTask(chainId string, startH, endH int64) *fixAckTxTask {
	return &fixAckTxTask{
		StartHeight: startH,
		EndHeight:   endH,
		ChainId:     chainId,
	}
}

func (t *fixAckTxTask) Run() int {
	if t.EndHeight < t.StartHeight || t.StartHeight < 0 {
		logrus.Println("EndHeight or StartHeight is invalid")
		return -1
	}
	height := t.StartHeight
	for {
		txs, err := txRepo.FindAllAckTxs(t.ChainId, height)
		if err != nil {
			logrus.Error(err.Error())
			return -1
		}
		if len(txs) > 0 {
			if err := t.doTask(t.ChainId, txs); err != nil {
				logrus.Error(err.Error())
			}
			logrus.Debugf("finish fix ack txs in height:%d-%d\n", height, height+constant.IncreHeight)
		}
		height += constant.IncreHeight
		if height > t.EndHeight {
			break
		}
	}
	return 1
}

func (t *fixAckTxTask) doTask(chainId string, txs []*entity.Tx) error {
	for _, val := range txs {
		valMsgs := make([]interface{}, 0, len(val.DocTxMsgs))
		for _, msg := range val.DocTxMsgs {
			if msg.Type == constant.MsgTypeAcknowledgement {
				ackmsg := msg.AckPacketMsg()
				utils.UnmarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(msg.Msg), &ackmsg)
				ackmsg.PacketId = fmt.Sprintf("%v%v%v%v%v", ackmsg.Packet.SourcePort, ackmsg.Packet.SourceChannel,
					ackmsg.Packet.DestinationPort, ackmsg.Packet.DestinationChannel, ackmsg.Packet.Sequence)
				valMsgs = append(valMsgs, ackmsg)
			} else {
				valMsgs = append(valMsgs, msg)
			}
		}
		if err := txRepo.UpdateAckPacketId(chainId, val.Height, val.TxHash, valMsgs); err != nil {
			return err
		}

	}
	return nil
}
