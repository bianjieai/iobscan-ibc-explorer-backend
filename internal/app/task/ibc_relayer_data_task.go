package task

/***
  ibc_relayer_data_task 一次性执行任务
  功能范围：
      1.根据已注册的relayer的地址、链信息，更新channel_pair_info字段。
      2.统计已注册的relayer相关信息（交易总数、成功交易总数、relayer费用总价值、交易总价值）。
      3.全量数据扫描发现未注册的relayer的数据，并保存到ibc_relayer表中。
      4.全量统计更新(包括已注册,未注册)relayer相关信息(交易总数、成功交易总数、relayer费用总价值、交易总价值)。
*/

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

var relayerDataTask RelayerDataTask
var _ OneOffTask = new(RelayerDataTask)

type RelayerDataTask struct {
	distRelayerMap map[string]string
	denomPriceMap  map[string]CoinItem
}

func (t *RelayerDataTask) Name() string {
	return "ibc_relayer_data_task"
}

func (t *RelayerDataTask) Switch() bool {
	return global.Config.Task.SwitchOnlyInitRelayerData
}

func (t *RelayerDataTask) Run() int {
	t.denomPriceMap = getTokenPriceMap()
	startTime := time.Now().Unix()

	doRegisterRelayer(t.denomPriceMap)

	historySegments, err := getHistorySegment(segmentStepHistory)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}
	//insert relayer data
	t.handleNewRelayerOnce(historySegments, true)

	segments, err := getSegment(segmentStepLatest)
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}
	//insert relayer data
	t.handleNewRelayerOnce(segments, false)

	t.aggrUnknowRelayerChannelPair()
	logrus.Infof("task %s finish deal, time use %d(s)", t.Name(), time.Now().Unix()-startTime)
	return 1
}

func (t *RelayerDataTask) initdistRelayerMap() {
	t.distRelayerMap = make(map[string]string, 20)
	skip := int64(0)
	limit := int64(1000)
	for {
		dbRelayers, err := relayerRepo.FindAll(skip, limit, repository.RelayerAllType)
		if err != nil {
			logrus.Error("find relayer by page fail, ", err.Error())
			return
		}

		for _, val := range dbRelayers {
			for _, pairInfo := range val.ChannelPairInfo {
				key := GenerateDistRelayerId(pairInfo.ChainA, pairInfo.ChainAAddress, pairInfo.ChainB, pairInfo.ChainBAddress)
				t.distRelayerMap[key] = val.RelayerId
				if val.RelayerName == "" {
					break
				}
			}
		}
		if len(dbRelayers) < int(limit) {
			break
		}
		skip += limit
	}

	return
}

func getRelayerPairIds(relayerChannelPairInfo []entity.ChannelPairInfo) []string {
	pairIds := make([]string, 0, len(relayerChannelPairInfo))
	for i := range relayerChannelPairInfo {
		val := relayerChannelPairInfo[i]
		if val.ChannelB != "" && val.ChannelA != "" {
			pairId := entity.GenerateRelayerPairId(val.ChainA, val.ChannelA, val.ChainAAddress,
				val.ChainB, val.ChannelB, val.ChainBAddress)
			pairIds = append(pairIds, pairId)
		}
	}
	return pairIds
}

func doRegisterRelayer(denomPriceMap map[string]CoinItem) {

	skip := int64(0)
	limit := int64(1000)
	for {
		relayers, err := relayerRepo.FindAll(skip, limit, repository.RelayerRegisterType)
		if err != nil {
			logrus.Error("find unknow relayer by page fail, ", err.Error())
			return
		}
		for _, relayer := range relayers {
			handleRegisterRelayerChannelPair(relayer)
			handleRegisterStatistic(denomPriceMap, relayer)
		}
		if len(relayers) < int(limit) {
			break
		}
		skip += limit
	}
}

func handleRegisterRelayerChannelPair(relayer *entity.IBCRelayerNew) {

	channelPairs, change, err := matchRegisterRelayerChannelPairInfo(relayer.ChannelPairInfo)
	if err != nil {
		logrus.Error("match register relayer channel pair fail, ", err.Error())
		return
	}
	if change {
		relayer.ChannelPairInfo = channelPairs
		if err := relayerRepo.UpdateChannelPairInfo(relayer.RelayerId, relayer.ChannelPairInfo); err != nil {
			logrus.Error("update register relayer statistic fail, ", err.Error())
		}
	}
	return
}

func handleRegisterStatistic(denomPriceMap map[string]CoinItem, relayer *entity.IBCRelayerNew) {
	item := getRelayerStatisticData(denomPriceMap, relayer)
	if err := relayerRepo.UpdateTxsInfo(item.RelayerId, item.RelayedTotalTxs, item.RelayedSuccessTxs,
		item.RelayedTotalTxsValue, item.TotalFeeValue); err != nil {
		logrus.Error("update register relayer statistic fail, ", err.Error())
	}
}

func (t *RelayerDataTask) handleNewRelayerOnce(segments []*segment, historyData bool) {
	t.initdistRelayerMap()
	doHandleSegments(t.Name(), 3, segments, historyData, t.doOneSegment)
}

func (t *RelayerDataTask) doOneSegment(v *segment, historyData bool) {
	var channelPairInfos []entity.ChannelPairInfo
	if historyData {
		channelPairInfos = t.handleIbcTxHistory(v.StartTime, v.EndTime)
	} else {
		channelPairInfos = t.handleIbcTxLatest(v.StartTime, v.EndTime)
	}
	if len(channelPairInfos) > 0 {
		newRelayerMap, dbRelayerMap := t.handleChannelPairInfo(channelPairInfos)
		newRelayer := distinctRelayerArr(newRelayerMap, false)
		dbRelayer := distinctRelayerArr(dbRelayerMap, true)
		if err := t.saveAndUpdateRelayer(newRelayer, dbRelayer); err != nil {
			logrus.Error("save and update relayer data error: ", err.Error(),
				fmt.Sprintf(" segment [%v:%v] ", v.StartTime, v.EndTime))
		}
		logrus.Infof("task %s find relayer finish segment [%v:%v]", t.Name(), v.StartTime, v.EndTime)
	}
}

func distinctRelayerArr(data map[string]entity.IBCRelayerNew, existInDb bool) []entity.IBCRelayerNew {

	retData := make([]entity.IBCRelayerNew, 0, len(data))
	if existInDb { //exist data in db
		relayerInDbSetMap := make(map[string]entity.IBCRelayerNew, len(data))

		for _, val := range data {
			value, ok := relayerInDbSetMap[val.RelayerId]
			if !ok {
				relayerInDbSetMap[val.RelayerId] = val
			} else {
				pairIds := make([]string, 0, len(value.ChannelPairInfo))
				for i := range value.ChannelPairInfo {
					item := value.ChannelPairInfo[i]
					pairId := entity.GenerateRelayerPairId(item.ChainA, item.ChannelA, item.ChainAAddress,
						item.ChainB, item.ChannelB, item.ChainBAddress)
					if !utils.InArray(pairIds, pairId) {
						pairIds = append(pairIds, pairId)
					}
				}

				//获取差异的channel_pair合并到一起
				for i := range val.ChannelPairInfo {
					item := val.ChannelPairInfo[i]
					pairId := entity.GenerateRelayerPairId(item.ChainA, item.ChannelA, item.ChainAAddress,
						item.ChainB, item.ChannelB, item.ChainBAddress)
					if !utils.InArray(pairIds, pairId) {
						pairIds = append(pairIds, pairId)
						value.ChannelPairInfo = append(value.ChannelPairInfo, val.ChannelPairInfo[i])
					}
				}
				relayerInDbSetMap[val.RelayerId] = value
			}
		}

		for _, val := range relayerInDbSetMap {
			retData = append(retData, val)
		}
	} else {
		for _, val := range data {
			retData = append(retData, val)
		}
	}
	return retData
}

func removeEmptyChannelData(addrPairInfo []entity.ChannelPairInfo) []entity.ChannelPairInfo {
	pairIds := getRelayerPairIds(addrPairInfo)
	addrPairInfoPairIdsLen := len(pairIds)
	if addrPairInfoPairIdsLen != len(addrPairInfo) {
		//需要清除chain,address在列表里已经存在但channel为空的数据
		dataMap := make(map[string]struct{}, len(addrPairInfo))
		channelPairArrs := make([]entity.ChannelPairInfo, 0, len(addrPairInfo))
		for _, val := range addrPairInfo {
			if val.PairId != "" {
				key := GenerateDistRelayerId(val.ChainA, val.ChainAAddress, val.ChainB, val.ChainBAddress)
				dataMap[key] = struct{}{}
				channelPairArrs = append(channelPairArrs, val)
			}
		}

		for _, val := range addrPairInfo {
			if val.PairId == "" {
				key := GenerateDistRelayerId(val.ChainA, val.ChainAAddress, val.ChainB, val.ChainBAddress)
				if _, ok := dataMap[key]; !ok {
					channelPairArrs = append(channelPairArrs, val)
				}
			}
		}

		return channelPairArrs
	}
	return addrPairInfo
}

func (t *RelayerDataTask) handleIbcTxLatest(startTime, endTime int64) []entity.ChannelPairInfo {
	relayerDtos, err := ibcTxRepo.GetRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.ChannelPairInfo
	for _, val := range relayerDtos {
		item := t.createChannelPairInfoData(val)
		if item.Valid() {
			relayers = append(relayers, item)
		}
	}
	return relayers
}

func (t *RelayerDataTask) handleIbcTxHistory(startTime, endTime int64) []entity.ChannelPairInfo {
	relayerDtos, err := ibcTxRepo.GetHistoryRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.ChannelPairInfo
	for _, val := range relayerDtos {
		item := t.createChannelPairInfoData(val)
		if item.Valid() {
			relayers = append(relayers, item)
		}
	}
	return relayers
}

func (t *RelayerDataTask) createChannelPairInfoData(dto *dto.GetRelayerInfoDTO) entity.ChannelPairInfo {
	channelPairInfo := entity.ChannelPairInfo{
		ChainA:        dto.ScChainId,
		ChainB:        dto.DcChainId,
		ChannelA:      dto.ScChannel,
		ChannelB:      dto.DcChannel,
		ChainBAddress: dto.DcChainAddress,
	}
	return channelPairInfo
}

func checkNoExist(pairId string, data entity.IBCRelayerNew) bool {
	for _, val := range data.ChannelPairInfo {
		if val.PairId == pairId {
			return false
		}
	}
	return true
}

func (t *RelayerDataTask) handleChannelPairInfo(channelPairInfos []entity.ChannelPairInfo) (map[string]entity.IBCRelayerNew, map[string]entity.IBCRelayerNew) {
	newRelayerMap := make(map[string]entity.IBCRelayerNew, 20)
	dbRelayerMap := make(map[string]entity.IBCRelayerNew, 20)
	for i := range channelPairInfos {
		//根据目标地址反查发起地址
		addrs := getChannalPairInfo(channelPairInfos[i])
		if len(addrs) > 0 {
			addrs = utils.DistinctSliceStr(addrs)
			channelPairInfos[i].ChainAAddress = addrs[0]
		}

		key := GenerateDistRelayerId(channelPairInfos[i].ChainA, channelPairInfos[i].ChainAAddress,
			channelPairInfos[i].ChainB, channelPairInfos[i].ChainBAddress)
		pairId := entity.GenerateRelayerPairId(channelPairInfos[i].ChainA, channelPairInfos[i].ChannelA,
			channelPairInfos[i].ChainAAddress, channelPairInfos[i].ChainB, channelPairInfos[i].ChannelB, channelPairInfos[i].ChainBAddress)
		if relayerId, ok := t.distRelayerMap[key]; ok {
			//数据库已存在，更新channel_pair_info
			dbData, cacheOk := dbRelayerMap[key]
			if !cacheOk {
				data, err := relayerRepo.FindOneByRelayerId(relayerId)
				if err != nil {
					logrus.Error(err.Error(), " relayerId:", relayerId)
					continue
				}
				channelpairNoExist := checkNoExist(pairId, *data)
				if channelpairNoExist {
					channelPairInfos[i].PairId = pairId
					data.ChannelPairInfo = append(data.ChannelPairInfo, channelPairInfos[i])
				}
				dbRelayerMap[key] = *data
			} else {
				channelpairNoExist := checkNoExist(pairId, dbData)
				if channelpairNoExist {
					channelPairInfos[i].PairId = pairId
					dbData.ChannelPairInfo = append(dbData.ChannelPairInfo, channelPairInfos[i])
				}
				dbRelayerMap[key] = dbData
			}

		} else {
			//新relayer数据
			if data, exist := newRelayerMap[key]; exist {
				channelpairNoExist := checkNoExist(pairId, data)
				if channelpairNoExist {
					channelPairInfos[i].PairId = pairId
					data.ChannelPairInfo = append(data.ChannelPairInfo, channelPairInfos[i])
					newRelayerMap[key] = data
				}
			} else {
				channelPairInfos[i].PairId = pairId
				newRelayerMap[key] = entity.IBCRelayerNew{
					RelayerId:       primitive.NewObjectID().Hex(),
					ServedChains:    2,
					ChannelPairInfo: []entity.ChannelPairInfo{channelPairInfos[i]},
					CreateAt:        time.Now().Unix(),
					UpdateAt:        time.Now().Unix(),
				}
			}

		}
	}
	return newRelayerMap, dbRelayerMap
}

func (t *RelayerDataTask) saveAndUpdateRelayer(newRelayerMap, dbRelayerMap []entity.IBCRelayerNew) error {
	for _, val := range dbRelayerMap {
		item := getRelayerStatisticData(t.denomPriceMap, &val)
		if err := relayerRepo.Update(item); err != nil {
			return err
		}
	}
	datas := make([]entity.IBCRelayerNew, 0, len(newRelayerMap))
	for _, val := range newRelayerMap {
		item := getRelayerStatisticData(t.denomPriceMap, &val)
		datas = append(datas, *item)
	}
	if len(datas) > 0 {
		if err := relayerRepo.InsertBatch(datas); err != nil {
			return err
		}
	}
	return nil
}

//根据已注册的relayer的地址和链更新channel_pair_info
func matchRegisterRelayerChannelPairInfo(addrPairInfo []entity.ChannelPairInfo) ([]entity.ChannelPairInfo, bool, error) {
	addrs, chains := getRelayerAddrAndChains(addrPairInfo)
	channelAddrs, err := relayerAddrChannelRepo.FindChannels(chains, addrs)
	if err != nil {
		return nil, false, err
	}

	channelMap := make(map[string]*entity.IBCRelayerAddressChannel, len(channelAddrs))
	for _, val := range channelAddrs {
		channelMap[val.Chain+val.RelayerAddress+":"+val.Channel+val.CounterPartyChannel] = val
	}

	chainChannelMap := make(map[string]entity.ChannelPairInfo, len(channelAddrs))
	for _, val := range addrPairInfo {
		chainChannelMap[val.ChainA+val.ChainAAddress] = val
	}
	pairIds := getRelayerPairIds(addrPairInfo)
	addrPairInfoPairIdsLen := len(pairIds)
	channelPairInfos := make([]entity.ChannelPairInfo, 0, len(addrPairInfo))
	for key, val := range channelMap {
		prefixKey := strings.Split(key, ":")[0]
		if value, ok := chainChannelMap[prefixKey]; ok {
			if val.RelayerAddress == value.ChainAAddress {
				value.ChannelB = val.CounterPartyChannel
				value.ChannelA = val.Channel

				//判断是否已存在
				pairId := entity.GenerateRelayerPairId(value.ChainA, value.ChannelA, value.ChainAAddress,
					value.ChainB, value.ChannelB, value.ChainBAddress)
				if !utils.InArray(pairIds, pairId) {
					channelPairInfos = append(channelPairInfos, value)
					pairIds = append(pairIds, pairId)
				}
			}
		}
	}
	addrPairInfoLen := len(addrPairInfo)
	addrPairInfo = removeEmptyChannelData(addrPairInfo)

	if len(channelPairInfos) == 0 {
		return addrPairInfo, addrPairInfoLen != len(addrPairInfo), nil
	}

	setMap := make(map[string]entity.ChannelPairInfo)
	retChannelPair := make([]entity.ChannelPairInfo, 0, len(channelPairInfos))
	for _, val := range channelPairInfos {
		val.PairId = entity.GenerateRelayerPairId(val.ChainA, val.ChannelA, val.ChainAAddress,
			val.ChainB, val.ChannelB, val.ChainBAddress)
		if _, ok := setMap[val.PairId]; ok {
			continue
		}
		setMap[val.PairId] = val
		delete(chainChannelMap, val.ChainA+val.ChainAAddress)
		retChannelPair = append(retChannelPair, val)
	}
	for _, val := range chainChannelMap {
		retChannelPair = append(retChannelPair, val)
	}

	if addrPairInfoPairIdsLen > 0 {
		retChannelPair = append(retChannelPair, addrPairInfo...)
	}
	return retChannelPair, true, nil
}

func getRelayerStatisticData(denomPriceMap map[string]CoinItem, data *entity.IBCRelayerNew) *entity.IBCRelayerNew {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		relayerFeeAmt := AggrRelayerFeeAmt(data)
		feeValue := caculateRelayerTotalValue(denomPriceMap, relayerFeeAmt)
		data.TotalFeeValue = feeValue.CoefficientInt64()
	}()

	go func() {
		defer wg.Done()
		relayerTxsAmt := AggrRelayerTxsAndAmt(data)
		totalTxsValue := caculateRelayerTotalValue(denomPriceMap, relayerTxsAmt)
		data.TotalFeeValue = totalTxsValue.CoefficientInt64()
		for _, val := range relayerTxsAmt {
			data.RelayedTotalTxs += val.Txs
			data.RelayedSuccessTxs += val.TxsSuccess
		}
	}()
	wg.Wait()
	return data
}

func (t *RelayerDataTask) aggrUnknowRelayerChannelPair() {
	skip := int64(0)
	limit := int64(1000)
	distRelayerMap := make(map[string]*entity.IBCRelayerNew, 1000)
	dumpData := make([]string, 0, 1000)
	for {
		relayers, err := relayerRepo.FindAll(skip, limit, repository.RelayerUnknowType)
		if err != nil {
			logrus.Error("find unknow relayer by page fail, ", err.Error())
			return
		}
		for _, relayer := range relayers {
			if len(relayer.ChannelPairInfo) > 0 {
				oneChannelPair := relayer.ChannelPairInfo[0]
				key := GenerateDistRelayerId(oneChannelPair.ChainA, oneChannelPair.ChainAAddress,
					oneChannelPair.ChainB, oneChannelPair.ChainBAddress)
				if data, ok := distRelayerMap[key]; ok {
					dumpData = append(dumpData, relayer.RelayerId)
					for i := range relayer.ChannelPairInfo {
						channelpairNoExist := checkNoExist(relayer.ChannelPairInfo[i].PairId, *data)
						if channelpairNoExist {
							data.ChannelPairInfo = append(data.ChannelPairInfo, relayer.ChannelPairInfo[i])
						}
					}
					data.TotalFeeValue += relayer.TotalFeeValue
					data.RelayedTotalTxsValue += relayer.RelayedTotalTxsValue
					data.RelayedTotalTxs += relayer.RelayedTotalTxs
					data.RelayedSuccessTxs += relayer.RelayedSuccessTxs

					distRelayerMap[key] = data
				} else {
					distRelayerMap[key] = relayer
				}

			}
		}
		if len(relayers) < int(limit) {
			break
		}
		skip += limit
	}
	err := relayerRepo.RemoveDumpData(dumpData)
	if err != nil {
		logrus.Error("remove dump data fail,", err.Error())
		return
	}
	for _, val := range distRelayerMap {
		err := relayerRepo.Update(val)
		if err != nil {
			logrus.Error(err.Error(), " relayerId:", val.RelayerId)
		}
	}

}
