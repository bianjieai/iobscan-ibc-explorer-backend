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
	//key: relayerId
	relayerTxsDataMap map[string]TxsItem
	//key: relayerId+Chain+Channel
	relayerValueMap map[string]decimal.Decimal
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

	CoinItem struct {
		Price float64
		Scale int
	}
)

type RelayerHandle func(relayer *entity.IBCRelayer)

func relayerAmtsMapKey(chainId, baseDenom, dcChainAddr, dcChannel string) string {
	return fmt.Sprintf("%s:%s:%s:%s", chainId, baseDenom, dcChainAddr, dcChannel)
}

func relayerAmtValueMapKey(relayerId, chainId, channel string) string {
	return fmt.Sprintf("%s:%s:%s", relayerId, chainId, channel)
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
	t.getTokenPriceMap()
	_ = t.todayStatistics()
	_ = t.yesterdayStatistics()
	t.cacheChainUnbondTimeFromLcd()
	t.updateIbcChainsRelayer()
	t.cacheIbcChannelRelayer()

	t.caculateRelayerTotalValue()
	t.CheckAndChangeRelayer(func(relayer *entity.IBCRelayer) {
		t.updateRelayerStatus(relayer)
		t.saveOrUpdateRelayerTxsAndValue(relayer)
	})

	return 1
}

//use cache map check relayer if exist
func filterDbExist(relayers []entity.IBCRelayer, historyData bool) []entity.IBCRelayer {
	dbRelayers, err := relayerCache.FindAll()
	if err != nil {
		return relayers
	}
	relayerMap := make(map[string]string, len(dbRelayers))
	for _, val := range dbRelayers {
		key := fmt.Sprintf("%s:%s:%s", val.ChainA, val.ChainAAddress, val.ChannelA)
		key1 := fmt.Sprintf("%s:%s:%s", val.ChainB, val.ChainBAddress, val.ChannelB)
		relayerMap[key] = ""
		relayerMap[key1] = ""
	}
	var distinctArr []entity.IBCRelayer
	for _, val := range relayers {
		if val.ChainAAddress == "" {
			val.ChainAAllAddress = getSrcChainAddress(&dto.GetRelayerInfoDTO{
				ScChainId:      val.ChainA,
				ScChannel:      val.ChannelA,
				DcChainId:      val.ChainB,
				DcChannel:      val.ChannelB,
				DcChainAddress: val.ChainBAddress,
			}, historyData)
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
func getSrcChainAddress(info *dto.GetRelayerInfoDTO, historyData bool) []string {
	//查询relayer在原链所有地址
	var (
		chainAAddress []string
		msgPacketId   string
	)

	if historyData {
		ibcTx, err := ibcTxRepo.GetHistoryOneRelayerScTxPacketId(info)
		if err == nil {
			msgPacketId = ibcTx.ScTxInfo.Msg.Msg.PacketId
		}
	} else {
		ibcTx, err := ibcTxRepo.GetOneRelayerScTxPacketId(info)
		if err == nil {
			msgPacketId = ibcTx.ScTxInfo.Msg.Msg.PacketId
		}
	}
	if msgPacketId != "" {
		scAddrs, err := txRepo.GetRelayerScChainAddr(msgPacketId, info.ScChainId)
		if err != nil {
			logrus.Errorf("get srAddr relayer packetId fail, %s", err.Error())
		}
		for _, val := range scAddrs {
			if val.ScChainAddress != "" {
				chainAAddress = append(chainAAddress, val.ScChainAddress)
			}
		}
	}
	return chainAAddress
}
func distinctRelayer(relayers []entity.IBCRelayer) []entity.IBCRelayer {
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

//dependence: cacheChainUnbondTimeFromLcd
func (t *IbcRelayerCronTask) updateRelayerStatus(relayer *entity.IBCRelayer) {
	timePeriod, updateTime, err := t.getTimePeriodAndupdateTime(relayer)
	if err != nil {
		logrus.Error("get relayer timePeriod and updateTime fail, ", err.Error())
		return
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
	t.handleOneRelayerStatusAndTime(relayer, updateTime, timePeriod)
	t.updateIbcChannelRelayerInfo(relayer, updateTime)
}
func (t *IbcRelayerCronTask) CheckAndChangeRelayer(handle func(relayer *entity.IBCRelayer)) {
	skip := int64(0)
	limit := int64(50)
	for {
		relayers, err := relayerRepo.FindAll(skip, limit)
		if err != nil {
			logrus.Error("find relayer by page fail, ", err.Error())
			return
		}
		for _, relayer := range relayers {
			handle(relayer)
		}
		if len(relayers) < int(limit) {
			break
		}
		skip += limit
	}
}

func (t *IbcRelayerCronTask) getTokenPriceMap() {
	coinIdPriceMap, _ := tokenPriceRepo.GetAll()
	baseDenoms, err := baseDenomCache.FindAll()
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
	mmdd := time.Now().Format(constant.TimeFormatMMDD)
	incr, _ := statisticsCheckRepo.GetIncr("getUnbondtimeFromLcd", mmdd)
	if incr > 1 {
		return
	}
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
		go func(baseUrl, chainId string) {
			getStakeParams(baseUrl, chainId)
			group.Done()
		}(baseUrl, val.ChainId)
	}
	group.Wait()
	_ = statisticsCheckRepo.Incr("getUnbondtimeFromLcd", mmdd)
}

func getStakeParams(baseUrl, chainId string) {
	bz, err := utils.HttpGet(baseUrl)
	if err != nil {
		logrus.Errorf(" staking %s params error, %v", baseUrl, err)
		return
	}

	var stakeparams vo.StakeParams
	err = json.Unmarshal(bz, &stakeparams)
	if err != nil {
		logrus.Errorf("unmarshal staking params error, %v", err)
		return
	}
	_ = unbondTimeCache.SetUnbondTime(chainId, stakeparams.Params.UnbondingTime)
}
func (t *IbcRelayerCronTask) handleOneRelayerStatusAndTime(relayer *entity.IBCRelayer, updateTime, timePeriod int64) {
	//Running=>Close: update_client 时间大于relayer基准周期
	status := entity.RelayerStatus(0)
	if relayer.TimePeriod > 0 && relayer.UpdateTime > 0 && relayer.TimePeriod < updateTime-relayer.UpdateTime {
		if relayer.Status == entity.RelayerRunning {
			status = entity.RelayerStop
		}
	}
	paths := t.getChannelsStatus(relayer.ChainA, relayer.ChainB)
	//Running=>Close: relayer中继通道只要出现状态不是STATE_OPEN
	if relayer.Status == entity.RelayerRunning {
		for _, path := range paths {
			if path.ChannelId == relayer.ChannelA {
				if path.State != constant.ChannelStateOpen {
					status = entity.RelayerStop
					break
				}
			}
			if path.Counterparty.ChannelId == relayer.ChannelB {
				if path.Counterparty.State != constant.ChannelStateOpen {
					status = entity.RelayerStop
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
					status = entity.RelayerRunning
				}
			}
		}
	}
	if status != entity.RelayerRunning && status != entity.RelayerStop {
		status = entity.RelayerStop
	}
	if err := relayerRepo.UpdateStatusAndTime(relayer.RelayerId, int(status), updateTime, timePeriod); err != nil {
		logrus.Error("update relayer about time_period and update_time fail, ", err.Error())
	}

}

// dependence: cacheIbcChannelRelayer
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
func collectTxs(data []*dto.CountRelayerPacketTxsCntDTO, startTime, endTime int64, hookTxs func(startTime, endTime int64) ([]*dto.CountRelayerPacketTxsCntDTO,
	error)) []*dto.CountRelayerPacketTxsCntDTO {

	relayerTxsCntDto, err := hookTxs(startTime, endTime)
	if err != nil {
		logrus.Error("collectTx hookTxs have fail, ", err.Error())
		return data
	}

	data = append(data, relayerTxsCntDto...)

	return data
}

func (t *IbcRelayerCronTask) AggrRelayerPacketTxs() {
	res, err := relayerStatisticsRepo.AggregateRelayerTxs()
	if err != nil {
		logrus.Error("aggregate relayer txs have fail, ", err.Error())
		return
	}
	t.relayerTxsDataMap = make(map[string]TxsItem, 20)
	for _, item := range res {
		key := item.RelayerId
		value, exist := t.relayerTxsDataMap[key]
		if exist {
			value.Txs += item.TotalTxs
			value.TxsSuccess += item.SuccessTotalTxs
			t.relayerTxsDataMap[key] = value
		} else {
			t.relayerTxsDataMap[key] = TxsItem{Txs: item.TotalTxs, TxsSuccess: item.SuccessTotalTxs}
		}
	}
}

func createAmounts(relayerAmounts []*dto.CountRelayerPacketAmountDTO) map[string]decimal.Decimal {
	relayerAmtsMap := make(map[string]decimal.Decimal, 20)
	for _, amt := range relayerAmounts {
		if amt.DcChainAddress != "" && amt.DcChainId != "" && amt.DcChannel != "" {
			key := relayerAmtsMapKey(amt.DcChainId, amt.BaseDenom, amt.DcChainAddress, amt.DcChannel)
			decAmt := decimal.NewFromFloat(amt.Amount)
			value, exist := relayerAmtsMap[key]
			if exist {
				value = value.Add(decAmt)
				relayerAmtsMap[key] = value
			} else {
				relayerAmtsMap[key] = decAmt
			}
		}
	}
	return relayerAmtsMap
}

func createIBCRelayerStatistics(channel, chainId, relayerId, baseDenom string, amount decimal.Decimal, successTxs, txs,
	startTime, endTime int64) entity.IBCRelayerStatistics {
	return entity.IBCRelayerStatistics{
		RelayerId:         relayerId,
		ChainId:           chainId,
		Channel:           channel,
		TransferBaseDenom: baseDenom,
		TransferAmount:    amount.String(),
		SuccessTotalTxs:   successTxs,
		TotalTxs:          txs,
		SegmentStartTime:  startTime,
		SegmentEndTime:    endTime,
		CreateAt:          time.Now().Unix(),
		UpdateAt:          time.Now().Unix(),
	}
}

func (t *IbcRelayerCronTask) caculateRelayerTotalValue() {
	baseDenomAmtDtos, err := relayerStatisticsRepo.CountRelayerBaseDenomAmt()
	if err != nil {
		logrus.Error("count relayer basedenom amount failed, ", err.Error())
		return
	}
	createAmtValue := func(baseDenomAmtDtos []*dto.CountRelayerBaseDenomAmtDTO) map[string]decimal.Decimal {
		relayerAmtValueMap := make(map[string]decimal.Decimal, 20)
		for _, amt := range baseDenomAmtDtos {
			key := relayerAmtValueMapKey(amt.RelayerId, amt.ChainId, amt.Channel)
			decAmt := decimal.NewFromFloat(amt.Amount)
			baseDenomValue := decimal.NewFromFloat(0)
			if coin, ok := t.denomPriceMap[amt.BaseDenom]; ok {
				if coin.Scale > 0 {
					baseDenomValue = decAmt.Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price))
				}
			}
			value, exist := relayerAmtValueMap[key]
			if exist {
				value = value.Add(baseDenomValue)
				relayerAmtValueMap[key] = value
			} else {
				relayerAmtValueMap[key] = baseDenomValue
			}
		}
		return relayerAmtValueMap
	}
	t.relayerValueMap = createAmtValue(baseDenomAmtDtos)
}

//dependence: caculateRelayerTotalValue, AggregateRelayerPacketTxs
func (t *IbcRelayerCronTask) saveOrUpdateRelayerTxsAndValue(val *entity.IBCRelayer) {
	t.AggrRelayerPacketTxs()
	getRelayerValue := func(data *entity.IBCRelayer) string {
		keyA := relayerAmtValueMapKey(data.RelayerId, data.ChainA, data.ChannelA)
		keyB := relayerAmtValueMapKey(data.RelayerId, data.ChainB, data.ChannelB)
		totalAValue, _ := t.relayerValueMap[keyA]
		totalBValue, _ := t.relayerValueMap[keyB]
		totalValue := decimal.NewFromFloat(0).Add(totalAValue).Add(totalBValue).
			Round(constant.DefaultValuePrecision).String()
		return totalValue
	}

	getRelayerTxs := func(data *entity.IBCRelayer) (int64, int64) {
		totalTxsAValue, _ := t.relayerTxsDataMap[data.RelayerId]
		totalTxsBValue, _ := t.relayerTxsDataMap[data.RelayerId]
		txsSuccess := totalTxsAValue.TxsSuccess + totalTxsBValue.TxsSuccess
		txs := totalTxsAValue.Txs + totalTxsBValue.Txs
		return txs, txsSuccess
	}

	totalValue := getRelayerValue(val)
	txs, txsSuccess := getRelayerTxs(val)
	if err := relayerRepo.UpdateTxsInfo(val.RelayerId, txs, txsSuccess, totalValue); err != nil && err != qmgo.ErrNoSuchDocuments {
		logrus.Error("update txs,txsSuccess,totalValue failed, ", err.Error())
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

func handleIbcTxLatest(startTime, endTime int64) []entity.IBCRelayer {
	relayerDtos, err := ibcTxRepo.GetRelayerInfo(startTime, endTime)
	if err != nil {
		logrus.Errorf("get relayer info fail, %s", err.Error())
		return nil
	}
	var relayers []entity.IBCRelayer
	for _, val := range relayerDtos {
		relayers = append(relayers, createRelayerData(val))
	}
	return relayers
}

func createRelayerData(dto *dto.GetRelayerInfoDTO) entity.IBCRelayer {
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

func (t *IbcRelayerCronTask) todayStatistics() error {
	logrus.Infof("task %s exec today statistics", t.Name())
	startTime, endTime := todayUnix()
	segments := []*segment{
		{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if err := relayerStatisticsTask.deal(segments, opUpdate); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}
	handleNewRelayerOnce(segments, false)

	return nil
}

func (t *IbcRelayerCronTask) yesterdayStatistics() error {
	mmdd := time.Now().Format(constant.TimeFormatMMDD)
	incr, _ := statisticsCheckRepo.GetIncr(t.Name(), mmdd)
	if incr > statisticsCheckTimes {
		return nil
	}

	logrus.Infof("task %s check yeaterday statistics, time: %d", t.Name(), incr)
	startTime, endTime := yesterdayUnix()
	segments := []*segment{
		{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if err := relayerStatisticsTask.deal(segments, opUpdate); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}
	handleNewRelayerOnce(segments, false)

	_ = statisticsCheckRepo.Incr(t.Name(), mmdd)
	return nil
}