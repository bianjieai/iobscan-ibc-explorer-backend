package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var _ OneOffTask = new(FixAckTxPacketIdTask)

type FixAckTxPacketIdTask struct {
}

func (f FixAckTxPacketIdTask) Name() string {
	return "fix_ack_tx_packet_id"
}

func (f FixAckTxPacketIdTask) Switch() bool {
	return false
}

func (f FixAckTxPacketIdTask) RunWithParam(chainsStr string) int {
	return f.handle(chainsStr)
}
func (f FixAckTxPacketIdTask) Run() int {
	return 1
	//return f.handle(global.Config.ChainConfig.FixAckTxPacketIdChains)
}

func (f FixAckTxPacketIdTask) handle(chainsStr string) int {
	newChainIds := strings.Split(chainsStr, ",")
	if len(newChainIds) == 0 {
		logrus.Errorf("task %s don't have fix ack packet_id chains", f.Name())
		return 1
	}
	// init queue
	chainQueue := new(utils.QueueString)
	for _, v := range newChainIds {
		chainQueue.Push(v)
	}
	fixChainCoordinator := &chainQueueCoordinator{
		chainQueue: chainQueue,
	}

	//handle chain
	handleChain := func(workerName string) int {
		chainId, err := fixChainCoordinator.getChain()
		if err != nil {
			logrus.Infof("chain_id %s worker %s exit", chainId, workerName)
			return 1
		}
		defer printExectime(fmt.Sprintf("chain:%s", chainId), time.Now().Unix())
		minHTx, err := txRepo.FindHeight(chainId, true)
		if err != nil {
			logrus.Errorf("find minHeight err chain_id:%s err:%v", chainId, err.Error())
			return -1
		}
		maxETx, err := txRepo.FindHeight(chainId, false)
		if err != nil {
			logrus.Errorf("skip fix ack_tx packet_id chain_id:%s err:%v", chainId, err.Error())
			return -1
		}
		logrus.Infof("start fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
			minHTx.Height-1, maxETx.Height, chainId)
		ret := NewfixAckTxTask(chainId, minHTx.Height-1, maxETx.Height).Run()
		if ret > 0 {
			logrus.Infof("finish fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
				minHTx.Height-1, maxETx.Height, chainId)
		} else {
			logrus.Errorf("fail fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
				minHTx.Height-1, maxETx.Height, chainId)
		}
		return 1
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(5)
	for i := 1; i <= 5; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			handleChain(wn)
			waitGroup.Done()
		}(workName)
	}
	waitGroup.Wait()

	return 1
}

//===========worker=================================

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
			logrus.Errorf("fix ack packet_id task in height:%d-%d chain_id:%s err:%s",
				height, height+constant.IncreHeight, t.ChainId, err.Error())
			return -1
		}
		if len(txs) > 0 {
			if err := t.doTask(txs); err != nil {
				logrus.Errorf("fix ack packet_id task height:%d-%d chain_id:%s err:%s",
					height, height+constant.IncreHeight, t.ChainId, err.Error())
				return -1
			}
			logrus.Infof("finish fix ack packet_id txs in height:%d-%d chain_id:%s",
				height, height+constant.IncreHeight, t.ChainId)
		}
		height += constant.IncreHeight
		logrus.Infof("finish scan %d-%d txs:%d chain_id:%s",
			height-constant.IncreHeight, height, len(txs), t.ChainId)
		if height > t.EndHeight {
			break
		}
	}
	return 1
}

func (t *fixAckTxTask) doTask(txs []*entity.Tx) error {
	type TxMsg struct {
		Type string      `bson:"type"`
		Msg  interface{} `bson:"msg"`
	}
	for _, val := range txs {
		valMsgs := make([]interface{}, 0, len(val.DocTxMsgs))
		msgsChange := false
		for _, msg := range val.DocTxMsgs {
			if msg.Type == constant.MsgTypeAcknowledgement {
				ackmsg := msg.AckPacketMsg()
				if ackmsg.PacketId == "" {
					msgsChange = true
					utils.UnmarshalJsonIgnoreErr(utils.MarshalJsonIgnoreErr(msg.Msg), &ackmsg)
					ackmsg.PacketId = fmt.Sprintf("%v%v%v%v%v", ackmsg.Packet.SourcePort, ackmsg.Packet.SourceChannel,
						ackmsg.Packet.DestinationPort, ackmsg.Packet.DestinationChannel, ackmsg.Packet.Sequence)
					valMsgs = append(valMsgs, TxMsg{Type: msg.Type, Msg: &ackmsg})
					continue
				}
			}
			valMsgs = append(valMsgs, msg)
		}
		if msgsChange {
			if err := txRepo.UpdateAckPacketId(t.ChainId, val.Height, val.TxHash, valMsgs); err != nil {
				return err
			}
		}

	}
	return nil
}
