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
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	denomPriceMap  map[string]dto.CoinItem
}

func (t *RelayerDataTask) Name() string {
	return "ibc_relayer_data_task"
}

func (t *RelayerDataTask) Switch() bool {
	return global.Config.Task.SwitchOnlyInitRelayerData
}

func (t *RelayerDataTask) Run() int {
	t.denomPriceMap = cache.TokenPriceMap()
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
				key := entity.GenerateDistRelayerId(pairInfo.ChainA, pairInfo.ChainAAddress, pairInfo.ChainB, pairInfo.ChainBAddress)
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
		if val.PairId != "" {
			pairIds = append(pairIds, val.PairId)
		}
	}
	return pairIds
}

func doRegisterRelayer(denomPriceMap map[string]dto.CoinItem) {

	skip := int64(0)
	limit := int64(1000)
	for {
		relayers, err := relayerRepo.FindAll(skip, limit, repository.RelayerRegisterType)
		if err != nil {
			logrus.Error("find unknow relayer by page fail, ", err.Error())
			return
		}
		for _, relayer := range relayers {
			handleRelayerChannelPair(relayer)
			handleRelayerStatistic(denomPriceMap, relayer)
		}
		if len(relayers) < int(limit) {
			break
		}
		skip += limit
	}
}

func handleRelayerChannelPair(relayer *entity.IBCRelayerNew) {

	channelPairs, change, err := matchRelayerChannelPairInfo(relayer.ChannelPairInfo)
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

func handleRelayerStatistic(denomPriceMap map[string]dto.CoinItem, relayer *entity.IBCRelayerNew) {
	item := getRelayerStatisticData(denomPriceMap, relayer)
	if err := relayerRepo.UpdateTxsInfo(item.RelayerId, item.RelayedTotalTxs, item.RelayedSuccessTxs,
		item.RelayedTotalTxsValue, item.TotalFeeValue); err != nil {
		logrus.Error("update register relayer statistic fail, ", err.Error())
	}
}

func (t *RelayerDataTask) handleNewRelayerOnce(segments []*segment, historyData bool) {
	t.initdistRelayerMap()
	doHandleSegments(t.Name(), 5, segments, historyData, t.doOneSegment)
}

func (t *RelayerDataTask) doOneSegment(v *segment, historyData bool) {
	var channelPairInfos []*entity.ChannelPairInfo
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
	emptyChannelData := make([]entity.ChannelPairInfo, 0, len(addrPairInfo))
	for i := range addrPairInfo {
		if addrPairInfo[i].ChannelA == "" || addrPairInfo[i].ChannelB == "" {
			emptyChannelData = append(emptyChannelData, addrPairInfo[i])
		}
	}
	if len(emptyChannelData) > 0 {
		//需要清除chain,address在列表里已经存在但channel为空的数据
		dataMap := make(map[string]struct{}, len(addrPairInfo))
		channelPairArrs := make([]entity.ChannelPairInfo, 0, len(addrPairInfo))
		for _, val := range addrPairInfo {
			if val.ChannelB != "" && val.ChannelA != "" {
				key := entity.GenerateDistRelayerId(val.ChainA, val.ChainAAddress, val.ChainB, val.ChainBAddress)
				dataMap[key] = struct{}{}
				channelPairArrs = append(channelPairArrs, val)
			}
		}

		for _, val := range emptyChannelData {
			key := entity.GenerateDistRelayerId(val.ChainA, val.ChainAAddress, val.ChainB, val.ChainBAddress)
			if _, ok := dataMap[key]; !ok {
				channelPairArrs = append(channelPairArrs, val)
			}
		}
		return channelPairArrs
	}
	return addrPairInfo
}

func (t *RelayerDataTask) handleIbcTxLatest(startTime, endTime int64) []*entity.ChannelPairInfo {
	relayerDtos, err := ibcTxRepo.GetRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []*entity.ChannelPairInfo
	for _, val := range relayerDtos {
		item := t.createChannelPairInfoData(val)
		if item.Valid() && item.ChainA != item.ChainB {
			relayers = append(relayers, &item)
		}
	}
	return relayers
}

func (t *RelayerDataTask) handleIbcTxHistory(startTime, endTime int64) []*entity.ChannelPairInfo {
	relayerDtos, err := ibcTxRepo.GetHistoryRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []*entity.ChannelPairInfo
	for _, val := range relayerDtos {
		item := t.createChannelPairInfoData(val)
		if item.Valid() && item.ChainA != item.ChainB {
			relayers = append(relayers, &item)
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
		ChainAAddress: dto.ScChainAddress,
		ChainBAddress: dto.DcChainAddress,
	}
	chainA, _ := entity.ConfirmRelayerPair(channelPairInfo.ChainA, channelPairInfo.ChainB)
	if chainA != channelPairInfo.ChainA {
		channelPairInfo = entity.ChannelPairInfo{
			ChainA:        dto.DcChainId,
			ChainB:        dto.ScChainId,
			ChannelA:      dto.DcChannel,
			ChannelB:      dto.ScChannel,
			ChainAAddress: dto.DcChainAddress,
			ChainBAddress: dto.ScChainAddress,
		}
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

func (t *RelayerDataTask) handleChannelPairInfo(channelPairInfos []*entity.ChannelPairInfo) (map[string]entity.IBCRelayerNew, map[string]entity.IBCRelayerNew) {

	newRelayerMap := make(map[string]entity.IBCRelayerNew, 20)
	dbRelayerMap := make(map[string]entity.IBCRelayerNew, 20)
	for i := range channelPairInfos {
		//忽略地址为空的channel_pair
		if channelPairInfos[i].ChainAAddress == "" || channelPairInfos[i].ChainBAddress == "" {
			continue
		}
		key := entity.GenerateDistRelayerId(channelPairInfos[i].ChainA, channelPairInfos[i].ChainAAddress,
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
					data.ChannelPairInfo = append(data.ChannelPairInfo, *channelPairInfos[i])
				}
				dbRelayerMap[key] = *data
			} else {
				channelpairNoExist := checkNoExist(pairId, dbData)
				if channelpairNoExist {
					channelPairInfos[i].PairId = pairId
					dbData.ChannelPairInfo = append(dbData.ChannelPairInfo, *channelPairInfos[i])
				}
				dbRelayerMap[key] = dbData
			}

		} else {
			//新relayer数据
			if data, exist := newRelayerMap[key]; exist {
				channelpairNoExist := checkNoExist(pairId, data)
				if channelpairNoExist {
					channelPairInfos[i].PairId = pairId
					data.ChannelPairInfo = append(data.ChannelPairInfo, *channelPairInfos[i])
					newRelayerMap[key] = data
				}
			} else {
				channelPairInfos[i].PairId = pairId
				newRelayerMap[key] = entity.IBCRelayerNew{
					RelayerId:       primitive.NewObjectID().Hex(),
					ServedChains:    2,
					ChannelPairInfo: []entity.ChannelPairInfo{*channelPairInfos[i]},
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
		if err := relayerRepo.Update(&val); err != nil {
			return err
		}
	}
	datas := make([]entity.IBCRelayerNew, 0, len(newRelayerMap))
	for _, val := range newRelayerMap {
		datas = append(datas, val)
	}
	if len(datas) > 0 {
		if err := relayerRepo.InsertBatch(datas); err != nil {
			return err
		}
	}
	return nil
}

//根据relayer的地址和链更新channel_pair_info
func matchRelayerChannelPairInfo(addrPairInfo []entity.ChannelPairInfo) ([]entity.ChannelPairInfo, bool, error) {
	pairIds := getRelayerPairIds(addrPairInfo)
	addrChannelPairInfos := make([]entity.ChannelPairInfo, 0, len(addrPairInfo))
	for _, val := range addrPairInfo {
		pairInfos, err := repository.GetChannelPairInfoByAddressPair(val.ChainA, val.ChainAAddress, val.ChainB, val.ChainBAddress)
		if err != nil {
			logrus.Error("GetChannelPairInfoByAddressPair fail, "+err.Error(),
				" chainA:", val.ChainA, " chainB:", val.ChainB, " chainAAddr:", val.ChainAAddress, " chainBAddr:", val.ChainBAddress)
			continue
		}
		addrChannelPairInfos = append(addrChannelPairInfos, pairInfos...)
	}

	//存放新增的channel_pair
	channelPairInfos := make([]entity.ChannelPairInfo, 0, len(addrPairInfo))
	for _, val := range addrChannelPairInfos {
		if !utils.InArray(pairIds, val.PairId) {
			channelPairInfos = append(channelPairInfos, val)
			pairIds = append(pairIds, val.PairId)
		}
	}
	addrPairInfoLen := len(addrPairInfo)
	addrPairInfo = removeEmptyChannelData(addrPairInfo)

	//没有新增的channel_pair
	if len(channelPairInfos) == 0 {
		return addrPairInfo, addrPairInfoLen != len(addrPairInfo), nil
	}

	addrPairInfo = append(addrPairInfo, channelPairInfos...)

	return addrPairInfo, true, nil
}

func getRelayerStatisticData(denomPriceMap map[string]dto.CoinItem, data *entity.IBCRelayerNew) *entity.IBCRelayerNew {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var (
		totalFeeValue                      decimal.Decimal
		totalTxsValue                      decimal.Decimal
		relayedTotalTxs, relayedSuccessTxs int64
	)
	go func() {
		defer wg.Done()
		relayedFeeTotalTxs := int64(0)
		relayerFeeAmt := AggrRelayerFeeAmt(data)
		totalFeeValue = caculateRelayerTotalValue(denomPriceMap, relayerFeeAmt)
		txsItem := make([]vo.DenomFeeItem, 0, len(relayerFeeAmt))
		for _, val := range relayerFeeAmt {
			relayedFeeTotalTxs += val.Txs
			txsItem = append(txsItem, vo.DenomFeeItem{
				Denom:      val.Denom,
				DenomChain: val.ChainId,
				Txs:        val.Txs,
				FeeValue:   val.AmtValue.String(),
			})
		}
		res := vo.TotalFeeCostResp{
			TotalTxs:        relayedFeeTotalTxs,
			TotalFeeValue:   totalFeeValue.String(),
			TotalDenomCount: int64(len(txsItem)),
			DenomList:       txsItem,
		}
		_ = relayerDataCache.SetTotalFeeCost(data.RelayerId, &res)
	}()

	go func() {
		defer wg.Done()
		relayerTxsAmt := AggrRelayerTxsAndAmt(data)
		txsItem := make([]vo.DenomTxsItem, 0, len(relayerTxsAmt))
		totalTxsValue = caculateRelayerTotalValue(denomPriceMap, relayerTxsAmt)
		for _, val := range relayerTxsAmt {
			txsItem = append(txsItem, vo.DenomTxsItem{
				BaseDenom:      val.Denom,
				BaseDenomChain: val.ChainId,
				Txs:            val.Txs,
				TxsValue:       val.AmtValue.String(),
			})
			relayedTotalTxs += val.Txs
			relayedSuccessTxs += val.TxsSuccess
		}

		res := vo.TotalRelayedValueResp{
			TotalTxs:        relayedTotalTxs,
			TotalTxsValue:   totalTxsValue.String(),
			TotalDenomCount: int64(len(txsItem)),
			DenomList:       txsItem,
		}
		_ = relayerDataCache.SetTotalRelayedValue(data.RelayerId, &res)

	}()
	wg.Wait()
	data.RelayedTotalTxsValue = totalTxsValue.String()
	data.TotalFeeValue = totalFeeValue.String()
	data.RelayedTotalTxs = relayedTotalTxs
	data.RelayedSuccessTxs = relayedSuccessTxs
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
				key := entity.GenerateDistRelayerId(oneChannelPair.ChainA, oneChannelPair.ChainAAddress,
					oneChannelPair.ChainB, oneChannelPair.ChainBAddress)
				if data, ok := distRelayerMap[key]; ok {
					dumpData = append(dumpData, relayer.RelayerId)
					for i := range relayer.ChannelPairInfo {
						channelpairNoExist := checkNoExist(relayer.ChannelPairInfo[i].PairId, *data)
						if channelpairNoExist {
							data.ChannelPairInfo = append(data.ChannelPairInfo, relayer.ChannelPairInfo[i])
						}
					}

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
		//caculate relayer statistic
		item := getRelayerStatisticData(t.denomPriceMap, val)
		err := relayerRepo.Update(item)
		if err != nil {
			logrus.Error(err.Error(), " relayerId:", val.RelayerId)
		}
	}

}
