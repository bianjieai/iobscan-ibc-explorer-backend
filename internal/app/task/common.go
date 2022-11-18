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
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type segment struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

func getHistorySegment(step int64) ([]*segment, error) {
	first, err := ibcTxRepo.FirstHistory()
	if err != nil {
		return nil, err
	}

	latest, err := ibcTxRepo.LatestHistory()
	if err != nil {
		return nil, err
	}

	start := time.Unix(first.CreateAt, 0)
	startUnix := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	end := time.Unix(latest.CreateAt, 0)
	endUnix := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 59, time.Local).Unix()

	var segments []*segment
	for temp := startUnix; temp < endUnix; temp += step {
		segments = append(segments, &segment{
			StartTime: temp,
			EndTime:   temp + step - 1,
		})
	}
	return segments, nil
}

func getSegment(step int64) ([]*segment, error) {
	first, err := ibcTxRepo.First()
	if err != nil {
		return nil, err
	}

	start := time.Unix(first.CreateAt, 0)
	startUnix := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local).Unix()
	end := time.Now()
	endUnix := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 59, time.Local).Unix()

	var segments []*segment
	for temp := startUnix; temp < endUnix; temp += step {
		segments = append(segments, &segment{
			StartTime: temp,
			EndTime:   temp + step - 1,
		})
	}

	return segments, nil
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

// todayUnix 获取今日第一秒和最后一秒的时间戳
func todayUnix() (int64, int64) {
	now := time.Now()
	startUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}

// yesterdayUnix 获取昨日第一秒和最后一秒的时间戳
func yesterdayUnix() (int64, int64) {
	date := time.Now().AddDate(0, 0, -1)
	startUnix := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}

func isConnectionErr(err error) bool {
	return true // 直接return true, 避免task被各种奇怪的返回值问题卡死
	//return strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "i/o timeout") ||
	//	strings.Contains(err.Error(), "unsupported protocol scheme")
}

func getAllChainMap() (map[string]*entity.ChainConfig, error) {
	allChainList, err := chainConfigRepo.FindAll()
	if err != nil {
		return nil, err
	}

	allChainMap := make(map[string]*entity.ChainConfig)
	for _, v := range allChainList {
		allChainMap[v.ChainName] = v
	}

	return allChainMap, err
}

// getChainIdNameMap, map key: chain id, value: chain name
func getChainIdNameMap() (map[string]string, error) {
	allChainList, err := chainConfigRepo.FindAllChainInfos()
	if err != nil {
		return nil, err
	}

	allChainMap := make(map[string]string)
	for _, v := range allChainList {
		allChainMap[v.CurrentChainId] = v.ChainName
	}

	return allChainMap, err
}

func matchDcInfo(scChainId, scPort, scChannel string, allChainMap map[string]*entity.ChainConfig) (dcChainId, dcPort, dcChannel string) {
	if allChainMap == nil || allChainMap[scChainId] == nil {
		return
	}

	for _, ibcInfo := range allChainMap[scChainId].IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.PortId == scPort && path.ChannelId == scChannel {
				dcChainId = path.Chain
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
func traceDenom(fullDenomPath, chainId string, allChainMap map[string]*entity.ChainConfig) *entity.IBCDenom {
	unix := time.Now().Unix()
	denom := calculateIbcHash(fullDenomPath)
	rootDenom := getRootDenom(fullDenomPath)
	if !strings.HasPrefix(denom, constant.IBCTokenPrefix) { // base denom
		return &entity.IBCDenom{
			Chain:          chainId,
			Denom:          denom,
			PrevDenom:      "",
			PrevChain:      "",
			BaseDenom:      denom,
			BaseDenomChain: chainId,
			DenomPath:      "",
			RootDenom:      rootDenom,
			IsBaseDenom:    true,
			CreateAt:       unix,
			UpdateAt:       unix,
		}
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("trace denom: %s, chain: %s, full path: %s, error. %v ", denom, chainId, fullDenomPath, err)
		}
	}()

	var currentChainId string
	var isBaseDenom bool
	currentChainId = chainId
	pathSplits := strings.Split(fullDenomPath, "/")
	denomPath := strings.Join(pathSplits[0:len(pathSplits)-1], "/")
	var TraceDenomList []*dto.DenomSimpleDTO
	TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
		Denom:   denom,
		ChainId: chainId,
	})

	for {
		if len(pathSplits) <= 1 {
			break
		}

		currentPort, currentChannel := pathSplits[0], pathSplits[1]
		tempPrevChainId, tempPrevPort, tempPrevChannel := matchDcInfo(currentChainId, currentPort, currentChannel, allChainMap)
		if tempPrevChainId == "" { // trace to end
			break
		} else {
			TraceDenomList = append(TraceDenomList, &dto.DenomSimpleDTO{
				Denom:   calculateIbcHash(strings.Join(pathSplits[2:], "/")),
				ChainId: tempPrevChainId,
			})
		}

		currentChainId, currentPort, currentChannel = tempPrevChainId, tempPrevPort, tempPrevChannel
		pathSplits = pathSplits[2:]
	}

	var prevDenom, prevChainId, baseDenom, baseDenomChainId string
	if len(TraceDenomList) == 1 { // denom is base denom
		isBaseDenom = true
		baseDenom = denom
		baseDenomChainId = chainId
	} else {
		isBaseDenom = false
		prevDenom = TraceDenomList[1].Denom
		prevChainId = TraceDenomList[1].ChainId
		baseDenom = TraceDenomList[len(TraceDenomList)-1].Denom
		baseDenomChainId = TraceDenomList[len(TraceDenomList)-1].ChainId
	}

	return &entity.IBCDenom{
		Chain:          chainId,
		Denom:          denom,
		PrevDenom:      prevDenom,
		PrevChain:      prevChainId,
		BaseDenom:      baseDenom,
		BaseDenomChain: baseDenomChainId,
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

// queryClientState 查询lcd client_state_path接口
func queryClientState(lcd, apiPath, port, channel string) (*vo.ClientStateResp, error) {
	apiPath = strings.ReplaceAll(apiPath, replaceHolderChannel, channel)
	apiPath = strings.ReplaceAll(apiPath, replaceHolderPort, port)
	url := fmt.Sprintf("%s%s", lcd, apiPath)

	if state, err := lcdTxDataCacheRepo.GetClientState(utils.Md5(url)); err == nil {
		return state, nil
	}

	bz, err := utils.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp vo.ClientStateResp
	err = json.Unmarshal(bz, &resp)
	if err != nil {
		return nil, err
	}

	_ = lcdTxDataCacheRepo.SetClientState(utils.Md5(url), &resp)
	return &resp, nil
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
