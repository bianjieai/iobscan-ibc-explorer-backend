package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
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

func (f FixAckTxPacketIdTask) RunWithParam(chainsStr string, endHeightStr string) int {
	return f.handle(chainsStr, endHeightStr)
}
func (f FixAckTxPacketIdTask) Run() int {
	return 1
	//return f.handle(global.Config.ChainConfig.FixAckTxPacketIdChains)
}

func (f FixAckTxPacketIdTask) handle(chainsStr string, endHeightStr string) int {
	newChainIds := strings.Split(chainsStr, ",")
	if len(newChainIds) == 0 {
		logrus.Errorf("task %s don't have valid chains(%s)", f.Name(), chainsStr)
		return 1
	}
	chainEndHeights := strings.Split(endHeightStr, ",")
	if len(chainEndHeights) == 0 {
		logrus.Errorf("task %s don't have valid end_height(%s)", f.Name(), endHeightStr)
		return 1
	}
	if len(chainEndHeights) != len(newChainIds) {
		logrus.Errorf("task %s chains(%s) don't match end_height(%s)", f.Name(), chainsStr, endHeightStr)
		return 1
	}
	chainHeightMap := make(map[string]int64, len(chainEndHeights))
	for _, val := range chainEndHeights {
		datas := strings.Split(val, ":")
		if len(datas) != 2 {
			logrus.Errorf("task %s don't have valid end_height(<chain_id:end_height>(%s))", f.Name(), val)
			return 1
		}
		endHeight, err := strconv.ParseInt(datas[1], 10, 64)
		if err != nil {
			logrus.Errorf("TaskController run %s err, %v", f.Name(), err)
			return 1
		}
		chainHeightMap[datas[0]] = endHeight
	}

	for _, chain := range newChainIds {
		if _, ok := chainHeightMap[chain]; !ok {
			logrus.Errorf("task %s chains(%s) don't match end_height", f.Name(), chain)
			return 1
		}
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
	handleChain := func(workerName, chainId string, maxEndHeight int64) int {
		minHTx, err := txRepo.FindHeight(chainId, true)
		if err != nil {
			logrus.Errorf("task_name:%s find minHeight err chain_id:%s err:%v", f.Name(), chainId, err.Error())
			return -1
		}
		//maxETx, err := txRepo.FindHeight(chainId, false)
		//if err != nil {
		//	logrus.Errorf("task_name:%s skip fix ack_tx packet_id chain_id:%s err:%v", f.Name(), chainId, err.Error())
		//	return -1
		//}
		logrus.Infof("task_name:%s fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
			f.Name(), minHTx.Height-1, maxEndHeight, chainId)
		ret := NewfixAckTxTask(chainId, f.Name(), minHTx.Height-1, maxEndHeight).Run()
		if ret > 0 {
			logrus.Infof("task_name:%s finish fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
				f.Name(), minHTx.Height-1, maxEndHeight, chainId)
		} else {
			logrus.Errorf("task_name:%s fail fix ack_tx packet_id,start-end:%v-%v chain_id:%s",
				f.Name(), minHTx.Height-1, maxEndHeight, chainId)
		}
		return 1
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(5)
	for i := 1; i <= 5; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			defer waitGroup.Done()
			for {
				chainId, err := fixChainCoordinator.getChain()
				if err != nil {
					logrus.Infof("task_name:%s chain_id %s worker %s exit", f.Name(), chainId, wn)
					return
				}
				handleChain(wn, chainId, chainHeightMap[chainId])
			}

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
	TaskName    string
}

func NewfixAckTxTask(chainId, taskName string, startH, endH int64) *fixAckTxTask {
	return &fixAckTxTask{
		StartHeight: startH,
		EndHeight:   endH,
		ChainId:     chainId,
		TaskName:    taskName,
	}
}

func (t *fixAckTxTask) Run() int {
	if t.EndHeight < t.StartHeight || t.StartHeight < 0 {
		logrus.Println("EndHeight or StartHeight is invalid")
		return -1
	}
	height := t.StartHeight
	endTime := int64(1640966400) //修补数据的截止时间2022-01-01 00:00:00
	for {
		txs, err := txRepo.FindAllAckTxs(t.ChainId, height)
		if err != nil {
			logrus.Errorf("fix ack packet_id task in height:%d-%d chain_id:%s err:%s",
				height, height+constant.IncreHeight, t.ChainId, err.Error())
			return -1
		}
		curTxTime := int64(0)
		if len(txs) > 0 {
			curTxTime = txs[0].Time
			if err := t.doTask(txs); err != nil {
				logrus.Errorf("task_name:%s fix ack packet_id task height:%d-%d chain_id:%s err:%s",
					t.TaskName, height, height+constant.IncreHeight, t.ChainId, err.Error())
				return -1
			}
			logrus.Infof("task_name:%s finish fix ack packet_id txs in height:%d-%d chain_id:%s",
				t.TaskName, height, height+constant.IncreHeight, t.ChainId)
		}
		height += constant.IncreHeight
		logrus.Infof("task_name:%s finish scan %d-%d txs:%d chain_id:%s",
			t.TaskName, height-constant.IncreHeight, height, len(txs), t.ChainId)
		if height > t.EndHeight || curTxTime > endTime {
			if curTxTime > endTime {
				logrus.Infof("task_name:%s finish for reach end_time(2022-01-01 00:00:00) chain_id:%s", t.TaskName, t.ChainId)
			}
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
