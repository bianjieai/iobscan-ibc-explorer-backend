package task

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IbcRelayerCronTask struct {
	relayerTxsMap  map[string]TxsItem
	relayerAmtsMap map[string]AmtItem
	//key: ChainA+ChainB+ChannelA+ChannelB
	channelRelayerCnt map[string]int64
	//key: BaseDenom
	denomPriceMap map[string]CoinItem
}
type (
	TxsItem struct {
		Txs        int64
		TxsSuccess int64
	}
	AmtItem struct {
		Amount decimal.Decimal
		Value  decimal.Decimal
	}
	CoinItem struct {
		Price float64
		Scale int
	}
)

func relayerTxsMapKey(chainId, dcChainAddr, dcChannel string) string {
	return fmt.Sprintf("%s:%s:%s", chainId, dcChainAddr, dcChannel)
}

func relayerAmtsMapKey(chainId, baseDenom, dcChainAddr, dcChannel string) string {
	return fmt.Sprintf("%s:%s:%s:%s", chainId, baseDenom, dcChainAddr, dcChannel)
}

func (t *IbcRelayerCronTask) Name() string {
	return "ibc_relayer_task"
}
func (t *IbcRelayerCronTask) Cron() int {
	if taskConf.CronTimeRelayerTask > 0 {
		return taskConf.CronTimeRelayerTask
	}
	return ThreeMinute
}

func (t *IbcRelayerCronTask) Run() int {
	t.handleNewRelayer()
	t.CheckAndChangeStatus()
	t.saveOrUpdateRelayerTxs()
	return 1
}

func (t *IbcRelayerCronTask) handleNewRelayer() {
	relayer, err := relayerRepo.FindLatestOne()
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		logrus.Errorf("findLatestone relayer fail, %s", err.Error())
		return
	}
	latestTxTime := int64(0)
	if relayer != nil {
		latestTxTime = relayer.LatestTxTime
	}
	currentLatestTxTime, _ := ibcTxRepo.GetLatestTxTime()
	relayersData := t.handleIbcTxLatest(latestTxTime)
	if len(relayersData) > 0 && currentLatestTxTime > latestTxTime {
		relayersData[len(relayersData)-1].LatestTxTime = currentLatestTxTime
	}
	relayersHistoryData := t.handleIbcTxHistory(latestTxTime)
	relayersData = append(relayersData, relayersHistoryData...)
	if len(relayersData) > 0 {
		relayersData = t.distinctRelayer(relayersData)
		relayersData = t.filterDbExist(relayersData)
		if err := relayerRepo.Insert(relayersData); err != nil && !qmgo.IsDup(err) {
			logrus.Error("insert  relayer data fail, ", err.Error())
		}
		t.updateIbcChainsRelayer()
		t.cacheIbcChannelRelayer()
		t.cacheChainUnbondTimeFromLcd()
	}
	t.CountRelayerPacketTxs()
	t.CountRelayerPacketTxsAmount()
}

//use cache map check relayer if exist
func (t *IbcRelayerCronTask) filterDbExist(relayers []entity.IBCRelayer) []entity.IBCRelayer {
	relayerMap, err := relayerCache.FindAll()
	if err != nil {
		return relayers
	}
	var distinctArr []entity.IBCRelayer
	for _, val := range relayers {
		if val.ChainAAddress == "" {
			val.ChainAAllAddress = t.getSrcChainAddress(&dto.GetRelayerInfoDTO{
				ScChainId:      val.ChainA,
				ScChannel:      val.ChannelA,
				DcChainId:      val.ChainB,
				DcChannel:      val.ChannelB,
				DcChainAddress: val.ChainBAddress,
			})
			if len(val.ChainAAllAddress) > 0 {
				val.ChainAAddress = val.ChainAAllAddress[0]
			}
		}
		key := fmt.Sprintf("%s:%s:%s", val.ChainA, val.ChainAAddress, val.ChannelA)
		key1 := fmt.Sprintf("%s:%s:%s", val.ChainB, val.ChainBAddress, val.ChannelB)
		_, exist := relayerMap[key]
		_, exist1 := relayerMap[key1]
		if !exist && !exist1 {
			val.RelayerId = utils.Md5(val.ChannelA + val.ChannelB + val.ChainA + val.ChainB + val.ChainAAddress + val.ChainBAddress)
			distinctArr = append(distinctArr, val)
		}
	}
	return distinctArr
}
func (t *IbcRelayerCronTask) getSrcChainAddress(val *dto.GetRelayerInfoDTO) []string {
	//查询relayer在原链所有地址
	var scAddrs []*dto.GetRelayerScChainAddreeDTO
	if ibcTx, err := ibcTxRepo.GetHistoryOneRelayerScTxPacketId(val); err == nil {
		scAddrs, err = txRepo.GetRelayerScChainAddr(ibcTx.ScTxInfo.Msg.Msg.PacketId, val.ScChainId)
		if err != nil {
			logrus.Errorf("get srAddr relayer packetId fail, %s", err.Error())
		}
	}
	chainAAddress := make([]string, 0, len(scAddrs))
	for _, val := range scAddrs {
		if val.ScChainAddress != "" {
			chainAAddress = append(chainAAddress, val.ScChainAddress)
		}
	}
	return chainAAddress
}
func (t *IbcRelayerCronTask) distinctRelayer(relayers []entity.IBCRelayer) []entity.IBCRelayer {
	distRelayerMap := make(map[string]bool, len(relayers))
	var distinctArr []entity.IBCRelayer
	checkSameMap := make(map[string]string, 20)
	//收集relayer双向链的地址
	for _, val := range relayers {
		if val.ChainBAddress != "" {
			keyB := fmt.Sprintf("%s%s", val.ChainB, val.ChannelB)
			checkSameMap[keyB] = val.ChainBAddress
		} else if val.ChainAAddress != "" {
			keyA := fmt.Sprintf("%s%s", val.ChainA, val.ChannelA)
			checkSameMap[keyA] = val.ChainAAddress
		}
	}
	for _, val := range relayers {
		//获取对方链地址信息
		if val.ChainBAddress == "" {
			key := fmt.Sprintf("%s%s", val.ChainB, val.ChannelB)
			val.ChainBAddress, _ = checkSameMap[key]
		}
		//获取对方链地址信息
		if val.ChainAAddress == "" {
			key := fmt.Sprintf("%s%s", val.ChainA, val.ChannelA)
			val.ChainAAddress, _ = checkSameMap[key]
		}
		key := fmt.Sprintf("%s:%s:%s", val.ChainA, val.ChannelA, val.ChainAAddress)
		key1 := fmt.Sprintf("%s:%s:%s", val.ChainB, val.ChannelB, val.ChainBAddress)
		_, exist := distRelayerMap[key]
		_, exist1 := distRelayerMap[key1]
		if !exist && !exist1 {
			distRelayerMap[key] = true
			distRelayerMap[key1] = true
			distinctArr = append(distinctArr, val)
		}
	}
	return distinctArr
}

func (t *IbcRelayerCronTask) CheckAndChangeStatus() {
	skip := int64(0)
	limit := int64(50)
	for {
		relayers, err := relayerRepo.FindAll(skip, limit)
		if err != nil {
			logrus.Error("find relayer by page fail, ", err.Error())
			return
		}
		for _, relayer := range relayers {
			timePeriod, updateTime, err := t.getTimePeriodAndupdateTime(relayer)
			if err != nil {
				logrus.Error("get relayer timePeriod and updateTime fail, ", err.Error())
				continue
			}
			if timePeriod == -1 {
				//get unbonding time from cache
				var chainAUnbondT, chainBUnbondT int64
				chainAUnbondTime, _ := unbondTimeCache.GetUnbondTime(relayer.ChainA)
				if chainAUnbondTime != "" {
					chainAUnbondT, _ = strconv.ParseInt(strings.ReplaceAll(chainAUnbondTime, "s", ""), 10, 64)
				}
				chainBUnbondTime, _ := unbondTimeCache.GetUnbondTime(relayer.ChainB)
				if chainBUnbondTime != "" {
					chainBUnbondT, _ = strconv.ParseInt(strings.ReplaceAll(chainBUnbondTime, "s", ""), 10, 64)
				}
				if chainAUnbondT > 0 && chainBUnbondT > 0 {
					if chainAUnbondT >= chainBUnbondT {
						timePeriod = chainBUnbondT
					} else {
						timePeriod = chainAUnbondT
					}
				}
			}
			t.handleOneRelayerInfo(relayer, updateTime, timePeriod)
			t.updateIbcChannelRelayerInfo(relayer, updateTime)
		}
		if len(relayers) < int(limit) {
			break
		}
		skip += limit
	}
}

func (t *IbcRelayerCronTask) getTokenPriceMap() {
	coinIdPriceMap, _ := tokenPriceRepo.GetAll()
	baseDenoms, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Error("find base_denom fail, ", err.Error())
		return
	}
	if len(coinIdPriceMap) == 0 {
		return
	}
	t.denomPriceMap = make(map[string]CoinItem, len(baseDenoms))
	for _, val := range baseDenoms {
		if price, ok := coinIdPriceMap[val.CoinId]; ok {
			t.denomPriceMap[val.Denom] = CoinItem{Price: price, Scale: val.Scale}
		}
	}
}

func (t *IbcRelayerCronTask) cacheChainUnbondTimeFromLcd() {
	configList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s cacheChainUnbondTimeFromLcd error, %v", t.Name(), err)
		return
	}
	if len(configList) == 0 {
		return
	}
	group := sync.WaitGroup{}
	group.Add(len(configList))
	for _, val := range configList {
		baseUrl := strings.ReplaceAll(fmt.Sprintf("%s%s", val.Lcd, val.LcdApiPath.ParamsPath), entity.ParamsModulePathPlaceholder, entity.StakeModule)
		value, err := unbondTimeCache.GetUnbondTime(val.ChainId)
		if err == nil && len(value) > 0 {
			continue
		}
		go func(baseUrl, chainId string) {
			defer group.Done()
			bz, err := utils.HttpGet(baseUrl)
			if err != nil {
				logrus.Errorf("task %s staking %s params error, %v", t.Name(), baseUrl, err)
				return
			}

			var stakeparams vo.StakeParams
			err = json.Unmarshal(bz, &stakeparams)
			if err != nil {
				logrus.Errorf("%s unmarshal staking params error, %v", t.Name(), err)
				return
			}
			_ = unbondTimeCache.SetUnbondTime(chainId, stakeparams.Params.UnbondingTime)
		}(baseUrl, val.ChainId)
	}
	group.Wait()
}

func (t *IbcRelayerCronTask) handleOneRelayerInfo(relayer *entity.IBCRelayer, updateTime, timePeriod int64) {
	//Running=>Close: update_client 时间大于relayer基准周期
	if relayer.TimePeriod > 0 && relayer.UpdateTime > 0 && relayer.TimePeriod < updateTime-relayer.UpdateTime {
		if relayer.Status == entity.RelayerRunning {
			relayer.Status = entity.RelayerStop
		}
	}
	paths := t.getChannelsStatus(relayer.ChainA, relayer.ChainB)
	//Running=>Close: relayer中继通道只要出现状态不是STATE_OPEN
	if relayer.Status == entity.RelayerRunning {
		for _, path := range paths {
			if path.ChannelId == relayer.ChannelA {
				if path.State != constant.ChannelStateOpen {
					relayer.Status = entity.RelayerStop
					break
				}
			}
			if path.Counterparty.ChannelId == relayer.ChannelB {
				if path.Counterparty.State != constant.ChannelStateOpen {
					relayer.Status = entity.RelayerStop
					break
				}
			}
		}
	} else {
		// Close=>Running: relayer的双向通道状态均为STATE_OPEN且update_client 时间小于relayer基准周期
		if relayer.TimePeriod > updateTime-relayer.UpdateTime {
			var channelStatus []string
			for _, path := range paths {
				if path.ChannelId == relayer.ChannelA {
					channelStatus = append(channelStatus, path.State)
				}
				if path.Counterparty.ChannelId == relayer.ChannelB {
					channelStatus = append(channelStatus, path.Counterparty.State)
				}
			}
			if len(channelStatus) == 2 {
				if channelStatus[0] == channelStatus[1] && channelStatus[0] == constant.ChannelStateOpen {
					relayer.Status = entity.RelayerRunning
				}
			}
		}
	}
	if relayer.Status != entity.RelayerRunning && relayer.Status != entity.RelayerStop {
		relayer.Status = entity.RelayerStop
	}
	if err := relayerRepo.UpdateStatusAndTime(relayer.RelayerId, int(relayer.Status), updateTime, timePeriod); err != nil {
		logrus.Error("update relayer about time_period and update_time fail, ", err.Error())
	}

}

func (t *IbcRelayerCronTask) updateIbcChannelRelayerInfo(relayer *entity.IBCRelayer, updateTime int64) {
	if len(t.channelRelayerCnt) > 0 || updateTime > 0 {
		var relayerCnt int64
		if len(t.channelRelayerCnt) > 0 {
			if relayercnt, ok := t.channelRelayerCnt[relayer.ChainA+relayer.ChainB+relayer.ChannelA+relayer.ChannelB]; ok {
				relayerCnt += relayercnt
			}
			if relayercnt, ok := t.channelRelayerCnt[relayer.ChainB+relayer.ChainA+relayer.ChannelB+relayer.ChannelA]; ok {
				relayerCnt += relayercnt
			}
		}

		ChannelId := generateChannelId(relayer.ChainA, relayer.ChannelA, relayer.ChainB, relayer.ChannelB)
		if err := channelRepo.UpdateOne(ChannelId, updateTime, relayerCnt); err != nil && err != mongo.ErrNoDocuments {
			logrus.Error("update ibc_channel about relayer fail, ", err.Error())
		}
	}
}

//set cache value redis key: ibc_channel_relayer_cnt
func (t *IbcRelayerCronTask) cacheIbcChannelRelayer() {
	channelRelayersDto, err := relayerRepo.CountChannelRelayers()
	if err != nil {
		logrus.Error("count channel relayers fail, ", err.Error())
		return
	}
	t.channelRelayerCnt = make(map[string]int64, len(channelRelayersDto))
	for _, one := range channelRelayersDto {
		t.channelRelayerCnt[one.ChainA+one.ChainB+one.ChannelA+one.ChannelB] = one.Count
	}
}
func collectTxs(data []*dto.CountRelayerPacketTxsCntDTO, hookTxs func() ([]*dto.CountRelayerPacketTxsCntDTO,
	error)) []*dto.CountRelayerPacketTxsCntDTO {

	relayerTxsCntDto, err := hookTxs()
	if err != nil {
		logrus.Error("collectTx hookTxs have fail, ", err.Error())
		return data
	}

	data = append(data, relayerTxsCntDto...)

	return data
}

//cache transferTxs,successTransferTx
func (t *IbcRelayerCronTask) CountRelayerPacketTxs() {
	relayerTxs := make([]*dto.CountRelayerPacketTxsCntDTO, 0, 20)
	relayerSuccessTxs := make([]*dto.CountRelayerPacketTxsCntDTO, 0, 20)
	//relayer txs count
	relayerTxs = collectTxs(relayerTxs, ibcTxRepo.CountRelayerPacketTxs)
	relayerTxs = collectTxs(relayerTxs, ibcTxRepo.CountHistoryRelayerPacketTxs)

	//relayer success txs count
	relayerSuccessTxs = collectTxs(relayerSuccessTxs, ibcTxRepo.CountRelayerSuccessPacketTxs)
	relayerSuccessTxs = collectTxs(relayerSuccessTxs, ibcTxRepo.CountHistoryRelayerSuccessPacketTxs)

	t.relayerTxsMap = make(map[string]TxsItem, 20)
	for _, tx := range relayerTxs {
		if tx.DcChainAddress != "" && tx.DcChainId != "" && tx.DcChannel != "" {
			key := relayerTxsMapKey(tx.DcChainId, tx.DcChainAddress, tx.DcChannel)
			value, exist := t.relayerTxsMap[key]
			if exist {
				value.Txs += tx.Count
				t.relayerTxsMap[key] = value
			} else {
				t.relayerTxsMap[key] = TxsItem{Txs: tx.Count}
			}
		}
	}

	for _, tx := range relayerSuccessTxs {
		if tx.DcChainAddress != "" && tx.DcChainId != "" && tx.DcChannel != "" {
			key := relayerTxsMapKey(tx.DcChainId, tx.DcChainAddress, tx.DcChannel)
			value, exist := t.relayerTxsMap[key]
			if exist {
				value.TxsSuccess += tx.Count
				t.relayerTxsMap[key] = value
			} else {
				t.relayerTxsMap[key] = TxsItem{TxsSuccess: tx.Count}
			}
		}
	}
}

func (t *IbcRelayerCronTask) CountRelayerPacketTxsAmount() {
	t.getTokenPriceMap()
	relayerAmounts := make([]*dto.CountRelayerPacketAmountDTO, 0, 20)
	if amounts, err := ibcTxRepo.CountRelayerPacketAmount(); err != nil {
		logrus.Error(err.Error())
	} else {
		relayerAmounts = append(relayerAmounts, amounts...)
	}
	if amounts, err := ibcTxRepo.CountHistoryRelayerPacketAmount(); err != nil {
		logrus.Error(err.Error())
	} else {
		relayerAmounts = append(relayerAmounts, amounts...)
	}

	createAmounts := func(relayerAmounts []*dto.CountRelayerPacketAmountDTO) map[string]AmtItem {
		relayerAmtsMap := make(map[string]AmtItem, 20)
		for _, amt := range relayerAmounts {
			if amt.DcChainAddress != "" && amt.DcChainId != "" && amt.DcChannel != "" {
				key := relayerAmtsMapKey(amt.DcChainId, amt.BaseDenom, amt.DcChainAddress, amt.DcChannel)
				decAmt := decimal.NewFromFloat(amt.Amount)
				baseDenomValue := decimal.NewFromFloat(0)
				if coin, ok := t.denomPriceMap[amt.BaseDenom]; ok {
					if coin.Scale > 0 {
						baseDenomValue = decAmt.DivRound(decimal.NewFromFloat(math.Pow10(coin.Scale)), constant.DefaultValuePrecision).Mul(decimal.NewFromFloat(coin.Price))
					}
				}
				value, exist := relayerAmtsMap[key]
				if exist {
					value.Amount = value.Amount.Add(decAmt)
					value.Value = value.Value.Add(baseDenomValue)
					relayerAmtsMap[key] = value
				} else {
					relayerAmtsMap[key] = AmtItem{Amount: decAmt, Value: baseDenomValue}
				}
			}
		}
		return relayerAmtsMap
	}

	t.relayerAmtsMap = createAmounts(relayerAmounts)
	t.caculateRelayerTotalValue()
}

//this function use data returned by CountRelayerPacketTxsAmount
func (t *IbcRelayerCronTask) caculateRelayerTotalValue() {
	var relayerStatics []entity.IBCRelayerStatistics
	for key, value := range t.relayerAmtsMap {
		if arrs := strings.Split(key, ":"); len(arrs) == 4 {
			chainId, baseDenom, relayerAddr, channel := arrs[0], arrs[1], arrs[2], arrs[3]
			relayerData, err := relayerRepo.FindRelayerId(chainId, relayerAddr, channel)
			if err != nil {
				if err != qmgo.ErrNoSuchDocuments {
					logrus.Warn(chainId, relayerAddr, channel, "find relayer id fail, ", err.Error())
				}
				continue
			}
			item := createIBCRelayerStatistics(channel, chainId, relayerData.RelayerId, baseDenom, value.Amount, value.Value)
			relayerStatics = append(relayerStatics, item)
		}
	}
	for _, val := range relayerStatics {
		if err := relayerStatisticsRepo.InserOrUpdate(val); err != nil {
			logrus.Error("insert or update relayer statistic fail ", err.Error())
		}
	}
	return
}

func createIBCRelayerStatistics(channel, chainId, relayerId, baseDenom string, amount, value decimal.Decimal) entity.IBCRelayerStatistics {
	return entity.IBCRelayerStatistics{
		RelayerId:          relayerId,
		ChainId:            chainId,
		Channel:            channel,
		TransferBaseDenom:  baseDenom,
		TransferAmount:     amount.String(),
		TransferTotalValue: value.Round(constant.DefaultValuePrecision).String(),
		CreateAt:           time.Now().Unix(),
		UpdateAt:           time.Now().Unix(),
	}
}
func (t *IbcRelayerCronTask) saveOrUpdateRelayerTxs() {
	if len(t.relayerTxsMap) > 0 {
		totalValueDtos, err := relayerStatisticsRepo.CountRelayerTotalValue()
		if err != nil {
			logrus.Error("count relayer transfer_total_txs_value failed, ", err.Error())
		}
		totalValueMap := make(map[string]float64, len(totalValueDtos))
		for _, val := range totalValueDtos {
			totalValueMap[val.ChainId+val.RelayerId+val.Channel] = val.Amount
		}

		for key, val := range t.relayerTxsMap {
			if arrs := strings.Split(key, ":"); len(arrs) == 3 {
				chainId, relayerAddr, channel := arrs[0], arrs[1], arrs[2]
				relayerData, err := relayerRepo.FindRelayerId(chainId, relayerAddr, channel)
				if err != nil {
					if err != qmgo.ErrNoSuchDocuments {
						logrus.Warn(chainId, relayerAddr, channel, "find relayer id fail, ", err.Error())
					}
					continue
				}
				totalValue, _ := totalValueMap[chainId+relayerData.RelayerId+channel]
				if err := relayerRepo.UpdateTxsInfo(relayerData.RelayerId, val.Txs, val.TxsSuccess, totalValue); err != nil {
					logrus.Error(err.Error())
				}
			}
		}
	}
}

func (t *IbcRelayerCronTask) getChannelsStatus(chainId, dcChainId string) []*entity.ChannelPath {
	// use cache find channels
	var ibcPaths []*entity.ChannelPath
	if paths, _ := ibcInfoCache.Get(chainId, dcChainId); paths != nil {
		data := paths.(string)
		utils.UnmarshalJsonIgnoreErr([]byte(data), &ibcPaths)
	}
	return ibcPaths
}

func (t *IbcRelayerCronTask) updateIbcChainsRelayer() {
	res, err := chainCache.FindAll()
	if err != nil {
		logrus.Error("find ibc_chains data fail, ", err.Error())
		return
	}

	for _, val := range res {
		relayerCnt, err := relayerRepo.FindRelayersCnt(val.ChainId)
		if err != nil {
			logrus.Error("count relayers of chain fail, ", err.Error())
			continue
		}
		if relayerCnt > 0 {
			if err := chainRepo.UpdateRelayers(val.ChainId, relayerCnt); err != nil {
				logrus.Error("update ibc_chain relayers fail, ", err.Error())
			}
		}
	}
	return
}

func (t *IbcRelayerCronTask) handleIbcTxLatest(latestTxTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetRelayerInfo(latestTxTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, val := range relayerDtos {
		relayers = append(relayers, t.createRelayerData(val))
	}
	return relayers
}

func (t *IbcRelayerCronTask) handleIbcTxHistory(latestTxTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetHistoryRelayerInfo(latestTxTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, val := range relayerDtos {
		relayers = append(relayers, t.createRelayerData(val))
	}
	return relayers
}

func (t *IbcRelayerCronTask) createRelayerData(dto *dto.GetRelayerInfoDTO) entity.IBCRelayer {
	return entity.IBCRelayer{
		ChainA:        dto.ScChainId,
		ChainB:        dto.DcChainId,
		ChannelA:      dto.ScChannel,
		ChannelB:      dto.DcChannel,
		ChainBAddress: dto.DcChainAddress,
		CreateAt:      time.Now().Unix(),
		UpdateAt:      time.Now().Unix(),
	}
}

//1: timePeriod
//2: updateTime
//3: error
func (t *IbcRelayerCronTask) getTimePeriodAndupdateTime(relayer *entity.IBCRelayer) (int64, int64, error) {
	updateTimeA, timePeriodA, err := txRepo.GetTimePeriodByUpdateClient(relayer.ChainA, relayer.ChainAAddress, relayer.UpdateTime)
	if err != nil {
		return 0, 0, err
	}
	updateTimeB, timePeriodB, err := txRepo.GetTimePeriodByUpdateClient(relayer.ChainB, relayer.ChainBAddress, relayer.UpdateTime)
	if err != nil {
		return 0, 0, err
	}
	timePeriod := timePeriodB
	updateTime := updateTimeB
	if timePeriodA >= timePeriodB && timePeriodB > 0 {
		// 两条链对应timePeriodB均不为-1，表示均正常取最大基准周期
		timePeriod = timePeriodA
		if updateTimeA >= updateTimeB {
			updateTime = updateTimeA
		}
	} else if timePeriodA == timePeriodB && timePeriodB == -1 {
		// 两条链对应timePeriodB均为-1，表示均正常取最大基准周期
		timePeriod = -1
		if updateTimeA >= updateTimeB {
			updateTime = updateTimeA
		}
	} else if timePeriodA == -1 || timePeriodB == -1 {
		//如果有一条链update_client没有查到，就不更新updateTime
		updateTime = relayer.UpdateTime
		timePeriod = relayer.TimePeriod
	}
	//判断更新时间如果小于历史更新时间，就不更新
	if updateTime < relayer.UpdateTime {
		updateTime = relayer.UpdateTime
	}
	return timePeriod, updateTime, nil
}
