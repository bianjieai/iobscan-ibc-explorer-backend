package task

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type segment struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

func segmentTool(step int64, startTime, endTime int64) []*segment {
	start := time.Unix(startTime, 0)
	startUnix := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	end := time.Unix(endTime, 0)
	endUnix := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 59, time.Local).Unix()

	var segments []*segment
	for temp := startUnix; temp < endUnix; temp += step {
		segments = append(segments, &segment{
			StartTime: temp,
			EndTime:   temp + step - 1,
		})
	}

	return segments
}

func isConnectionErr(err error) bool {
	return true // 直接return true, 避免task被各种奇怪的返回值问题卡死
	//return strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "i/o timeout") ||
	//	strings.Contains(err.Error(), "unsupported protocol scheme")
}

func matchDcInfo(scChain, scPort, scChannel string, allChainMap map[string]*entity.ChainConfig) (dcChain, dcPort, dcChannel string) {
	if allChainMap == nil || allChainMap[scChain] == nil {
		return
	}

	for _, ibcInfo := range allChainMap[scChain].IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.PortId == scPort && path.ChannelId == scChannel {
				dcChain = path.Chain
				dcPort = path.Counterparty.PortId
				dcChannel = path.Counterparty.ChannelId
				return
			}
		}
	}

	return
}

// getRootDenom get root denom by denom path
//   - fullPath full fullPath, eg："transfer/channel-1/uiris", "uatom"
func getRootDenom(fullPath string) string {
	split := strings.Split(fullPath, "/")
	return split[len(split)-1]
}

// splitFullPath get denom path and root denom from denom path
//   - fullPath full fullPath, eg："transfer/channel-1/uiris", "uatom"
func splitFullPath(fullPath string) (denomPath, rootDenom string) {
	pathSplits := strings.Split(fullPath, "/")
	denomPath = strings.Join(pathSplits[0:len(pathSplits)-1], "/")
	rootDenom = pathSplits[len(pathSplits)-1]
	return
}

// calculateIbcHash calculate denom hash by denom path
//   - fullPath full fullPath, eg："transfer/channel-1/uiris", "uatom"
func calculateIbcHash(fullPath string) string {
	if len(strings.Split(fullPath, "/")) == 1 {
		return fullPath
	}

	hash := utils.Sha256(fullPath)
	return fmt.Sprintf("%s/%s", constant.IBCTokenPrefix, strings.ToUpper(hash))
}

// traceDenom trace denom path, parse denom info
//   - fullDenomPath denom full path，eg："transfer/channel-1/uiris", "uatom"
func traceDenom(fullDenomPath, chain string, allChainMap map[string]*entity.ChainConfig) *entity.IBCDenom {
	unix := time.Now().Unix()
	denom := calculateIbcHash(fullDenomPath)
	rootDenom := getRootDenom(fullDenomPath)
	if !strings.HasPrefix(denom, constant.IBCTokenPrefix) { // base denom
		return &entity.IBCDenom{
			Chain:          chain,
			Denom:          denom,
			PrevDenom:      "",
			PrevChain:      "",
			BaseDenom:      denom,
			BaseDenomChain: chain,
			DenomPath:      "",
			RootDenom:      rootDenom,
			IsBaseDenom:    true,
			CreateAt:       unix,
			UpdateAt:       unix,
		}
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("trace denom: %s, chain: %s, full path: %s, error. %v ", denom, chain, fullDenomPath, err)
		}
	}()

	var currentChain string
	var isBaseDenom bool
	currentChain = chain
	pathSplits := strings.Split(fullDenomPath, "/")
	denomPath := strings.Join(pathSplits[0:len(pathSplits)-1], "/")
	var TraceDenomList []*dto.DenomSimpleDTO
	TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
		Denom: denom,
		Chain: chain,
	})

	for {
		if len(pathSplits) <= 1 {
			break
		}

		currentPort, currentChannel := pathSplits[0], pathSplits[1]
		tempPrevChain, tempPrevPort, tempPrevChannel := matchDcInfo(currentChain, currentPort, currentChannel, allChainMap)
		if tempPrevChain == "" { // trace to end
			break
		} else {
			TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
				Denom: calculateIbcHash(strings.Join(pathSplits[2:], "/")),
				Chain: tempPrevChain,
			})
		}

		currentChain, currentPort, currentChannel = tempPrevChain, tempPrevPort, tempPrevChannel
		pathSplits = pathSplits[2:]
	}

	var prevDenom, prevChain, baseDenom, baseDenomChain string
	if len(TraceDenomList) == 1 { // denom is base denom
		isBaseDenom = true
		baseDenom = denom
		baseDenomChain = chain
	} else {
		isBaseDenom = false
		prevDenom = TraceDenomList[1].Denom
		prevChain = TraceDenomList[1].Chain
		baseDenom = TraceDenomList[len(TraceDenomList)-1].Denom
		baseDenomChain = TraceDenomList[len(TraceDenomList)-1].Chain
	}

	return &entity.IBCDenom{
		Chain:          chain,
		Denom:          denom,
		PrevDenom:      prevDenom,
		PrevChain:      prevChain,
		BaseDenom:      baseDenom,
		BaseDenomChain: baseDenomChain,
		DenomPath:      denomPath,
		RootDenom:      rootDenom,
		IsBaseDenom:    isBaseDenom,
		CreateAt:       unix,
		UpdateAt:       unix,
	}
}

// calculateNextDenomPath calculate full denom path of next hop.
// return full denom path and cross back identification
func calculateNextDenomPath(packet model.Packet) (string, bool) {
	prefixSc := fmt.Sprintf("%s/%s/", packet.SourcePort, packet.SourceChannel)
	prefixDc := fmt.Sprintf("%s/%s/", packet.DestinationPort, packet.DestinationChannel)
	denomPath := packet.Data.Denom
	if strings.HasPrefix(denomPath, prefixSc) { // transfer to prev chain
		denomPath = strings.Replace(denomPath, prefixSc, "", 1)
		return denomPath, true
	} else {
		denomPath = fmt.Sprintf("%s%s", prefixDc, denomPath)
		return denomPath, false
	}
}

// parseTransferTxEvents parse ibc info from events of transfer tx
func parseTransferTxEvents(msgIndex int, tx *entity.Tx) (dcPort, dcChannel, denomFullPath, sequence, scConnection string) {
	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if evt.Type == "send_packet" {
				for _, attr := range evt.Attributes {
					switch attr.Key {
					case "packet_dst_port":
						dcPort = attr.Value
					case "packet_dst_channel":
						dcChannel = attr.Value
					case "packet_sequence":
						sequence = attr.Value
					case "packet_data":
						var data model.TransferTxPacketData
						_ = json.Unmarshal([]byte(attr.Value), &data)
						denomFullPath = data.Denom
					case "packet_connection":
						scConnection = attr.Value
					default:
					}
				}
			}
		}
	}

	return
}

// parseRecvPacketTxEvents parse ibc info from events of recv packet tx
func parseRecvPacketTxEvents(msgIndex int, tx *entity.Tx) (dcConnection, packetAck string, existPacketAck bool) {
	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if evt.Type == "recv_packet" {
				for _, attr := range evt.Attributes {
					switch attr.Key {
					case "packet_connection":
						dcConnection = attr.Value
					default:
					}
				}
			}

			if evt.Type == "write_acknowledgement" {
				for _, attr := range evt.Attributes {
					switch attr.Key {
					case "packet_ack":
						packetAck = attr.Value
						existPacketAck = true
					default:
					}
				}
			}
		}
	}

	return
}

// parseAckPacketTxEvents parse ibc info from events of ack packet tx
func parseAckPacketTxEvents(msgIndex int, tx *entity.Tx) (existTransferEvent bool) {
	if len(tx.EventsNew) > msgIndex {
		for _, evt := range tx.EventsNew[msgIndex].Events {
			if evt.Type == "transfer" {
				existTransferEvent = true
				return
			}
		}
	}
	return
}

//并发处理全量数据
func doHandleSegments(taskName string, workNum int, segments []*segment, isTargetHistory bool, dowork WorkerExecHandler) {
	if workNum <= 0 {
		return
	}
	st := time.Now().Unix()
	logrus.Infof("task %s worker group start, target history: %t", taskName, isTargetHistory)
	defer func() {
		logrus.Infof("task %s worker group end, target history: %t, time use: %d(s)", taskName, isTargetHistory, time.Now().Unix()-st)
	}()
	var wg sync.WaitGroup
	wg.Add(workNum)
	for i := 0; i < workNum; i++ {
		num := i
		go func(num int) {
			defer wg.Done()

			for id, v := range segments {
				if id%workNum != num {
					continue
				}
				logrus.Infof("task %s worker %d fix %d-%d, target history: %t", taskName, num, v.StartTime, v.EndTime, isTargetHistory)
				dowork(v, isTargetHistory)
			}
		}(num)
	}
	wg.Wait()
}

type WorkerExecHandler func(seg *segment, isTargetHistory bool)
