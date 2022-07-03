package task

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelTask struct {
	allChannelIds    []string
	channelStatusMap map[string]entity.ChannelStatus
	baseDenomMap     entity.IBCBaseDenomMap // 所有的base denom
	chainTxsMap      map[string]int64
	chainTxsValueMap map[string]decimal.Decimal
}

func (t *ChannelTask) Name() string {
	return "ibc_channel_task"
}

func (t *ChannelTask) Cron() int {
	if taskConf.CronTimeChannelTask > 0 {
		return taskConf.CronTimeChannelTask
	}
	return ThreeMinute
}

func (t *ChannelTask) Run() int {
	t.clear()
	if err := t.analyzeChainConfig(); err != nil {
		return -1
	}

	existedChannelList, newChannelList, err := t.getAllChannel()
	if err != nil {
		return -1
	}

	// 部分数据统计出错可以直接忽略error,继续计算后面的指标
	_ = t.setLatestSettlementTime(existedChannelList, newChannelList)

	t.setStatusAndOperatingPeriod(existedChannelList, newChannelList)

	_ = t.todayStatistics()

	_ = t.yesterdayStatistics()

	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}
	t.baseDenomMap = baseDenomList.ConvertToMap()

	if err = t.setTransferTxs(existedChannelList, newChannelList); err != nil { // 计算txs和交易价值，同时更新ibc_channel_statistics
		logrus.Errorf("task %s setTransferTxs error, %v", t.Name(), err)
		return -1
	}

	if err = channelRepo.InsertBatch(newChannelList); err != nil {
		logrus.Errorf("task %s InsertBatch error, %v", t.Name(), err)
	}

	for _, v := range existedChannelList {
		if err = channelRepo.UpdateChannel(v); err != nil && err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s UpdateChannel error, %v", t.Name(), err)
		}
	}

	// 更新ibc_chain
	for chainId, txs := range t.chainTxsMap {
		txsValue := t.chainTxsValueMap[chainId].Round(constant.DefaultValuePrecision).String()
		if err = chainRepo.UpdateTransferTxs(chainId, txs, txsValue); err != nil && err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s update chain %s error, %v", t.Name(), chainId, err)
		}
	}
	return 1
}

func (t *ChannelTask) clear() {
	t.chainTxsMap = make(map[string]int64)
	t.chainTxsValueMap = make(map[string]decimal.Decimal)
}

func (t *ChannelTask) analyzeChainConfig() error {
	confList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s analyzeChainConfig error, %v", t.Name(), err)
		return err
	}

	var channelIds []string
	channelStatusMap := make(map[string]entity.ChannelStatus)

	var chainA, channelA, chainB, channelB string
	for _, v := range confList {
		chainA = v.ChainId
		for _, info := range v.IbcInfo {
			chainB = info.ChainId
			for _, p := range info.Paths {
				channelA = p.ChannelId
				channelB = p.Counterparty.ChannelId
				channelId := generateChannelId(chainA, channelA, chainB, channelB)

				if utils.InArray(channelIds, channelId) { // 已经存在
					continue
				}

				channelIds = append(channelIds, channelId)
				if p.State == constant.ChannelStateOpen || p.Counterparty.State == constant.ChannelStateOpen {
					channelStatusMap[channelId] = entity.ChannelStatusOpened
				} else {
					channelStatusMap[channelId] = entity.ChannelStatusClosed
				}
			}
		}
	}

	t.allChannelIds = channelIds
	t.channelStatusMap = channelStatusMap
	return nil
}

func generateChannelId(chainA, channelA, chainB, channelB string) string {
	if strings.Contains(strings.ToLower(chainA), constant.Cosmos) {
		return fmt.Sprintf("%s|%s|%s|%s", chainA, channelA, chainB, channelB)
	}

	if strings.Contains(strings.ToLower(chainB), constant.Cosmos) {
		return fmt.Sprintf("%s|%s|%s|%s", chainB, channelB, chainA, channelA)
	}

	if strings.Contains(strings.ToLower(chainA), constant.Iris) {
		return fmt.Sprintf("%s|%s|%s|%s", chainA, channelA, chainB, channelB)
	}

	if strings.Contains(strings.ToLower(chainB), constant.Iris) {
		return fmt.Sprintf("%s|%s|%s|%s", chainB, channelB, chainA, channelA)
	}

	compare := strings.Compare(strings.ToLower(chainA), strings.ToLower(chainB))
	if compare < 0 {
		return fmt.Sprintf("%s|%s|%s|%s", chainA, channelA, chainB, channelB)
	} else {
		return fmt.Sprintf("%s|%s|%s|%s", chainB, channelB, chainA, channelA)
	}
}

func (t *ChannelTask) parseChannelId(channelId string) (chainA, channelA, chainB, channelB string, err error) {
	split := strings.Split(channelId, "|")
	if len(split) != 4 {
		logrus.Errorf("task %s parseChannelId error, %v", t.Name(), err)
		return "", "", "", "", fmt.Errorf("channel id format error")
	}
	return split[0], split[1], split[2], split[3], nil
}

func (t *ChannelTask) getAllChannel() (entity.IBCChannelList, entity.IBCChannelList, error) {
	existedChannelList, err := channelRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s getAllChannel error, %v", t.Name(), err)
		return nil, nil, err
	}

	existedIds := existedChannelList.GetChannelIds()
	var newChannelList entity.IBCChannelList
	for _, v := range t.allChannelIds {
		isExist := false
		for _, e := range existedIds {
			if v == e {
				isExist = true
				break
			}
		}

		if isExist {
			continue
		}

		newChannelList = append(newChannelList)
		chainA, channelA, chainB, channelB, err := t.parseChannelId(v)
		if err != nil {
			return nil, nil, err
		}

		newChannelList = append(newChannelList, &entity.IBCChannel{
			ChannelId:        v,
			ChainA:           chainA,
			ChainB:           chainB,
			ChannelA:         channelA,
			ChannelB:         channelB,
			Status:           entity.ChannelStatusOpened, // 默认开启状态
			OperatingPeriod:  0,
			LatestOpenTime:   0,
			Relayers:         0,
			TransferTxs:      0,
			TransferTxsValue: "",
			CreateAt:         time.Now().Unix(),
			UpdateAt:         time.Now().Unix(),
		})
	}

	return existedChannelList, newChannelList, nil
}

func (t *ChannelTask) setLatestSettlementTime(existedChannelList entity.IBCChannelList, newChannelList entity.IBCChannelList) error {
	// todo
	for _, v := range newChannelList {
		// 查询,初始的LatestSettlementTime 为channel的 open confirm 时间
		// channel open confirm 时间的获取当前从配置读取
		chanConf, err := channelConfigRepo.Find(v.ChainA, v.ChannelA, v.ChainB, v.ChannelB)
		if err != nil {
			continue
		}
		v.LatestOpenTime = chanConf.ChannelOpenAt
	}

	for _, v := range existedChannelList {
		// 之前没有设置open 时间且是open状态的
		if v.LatestOpenTime == 0 && v.Status == entity.ChannelStatusOpened {
			if chanConf, err := channelConfigRepo.Find(v.ChainA, v.ChannelA, v.ChainB, v.ChannelB); err == nil {
				v.LatestOpenTime = chanConf.ChannelOpenAt
			}
		}

		// 之前关闭了,现在重新打开channel
		if v.Status == entity.ChannelStatusClosed && t.channelStatusMap[v.ChannelId] == entity.ChannelStatusOpened {
			// 查询

		}
	}
	return nil
}

func (t *ChannelTask) setStatusAndOperatingPeriod(existedChannelList entity.IBCChannelList, newChannelList entity.IBCChannelList) {
	set := func(list entity.IBCChannelList) {
		for _, v := range list {
			currentStatus, ok := t.channelStatusMap[v.ChannelId]
			if !ok {
				currentStatus = entity.ChannelStatusOpened
			}

			if v.LatestOpenTime == 0 { // channel open 时间不确定，设置状态，处理下一个
				v.Status = currentStatus
				continue
			}

			// 1、channel 一直是close的, 持续工作时间不变
			// 2、channel 从open->close, close->open, open->open 状态变化时，持续工作时间更新
			if v.Status == entity.ChannelStatusClosed && currentStatus == entity.ChannelStatusClosed {
				continue
			}

			now := time.Now().Unix()
			v.OperatingPeriod = now - v.LatestOpenTime
			v.Status = currentStatus
		}
	}

	set(existedChannelList)
	set(newChannelList)
}

func (t *ChannelTask) setTransferTxs(existedChannelList entity.IBCChannelList, newChannelList entity.IBCChannelList) error {
	statistics, err := channelStatisticsRepo.Aggr()
	if err != nil {
		logrus.Errorf("task %s channelStatisticsRepo.Aggr error, %v", t.Name(), err)
		return err
	}

	for _, v := range existedChannelList {
		count, value := t.calculateChannelStatistics(v.ChannelId, statistics)
		v.TransferTxs = count
		v.TransferTxsValue = value.Round(constant.DefaultValuePrecision).String()
	}

	for _, v := range newChannelList {
		count, value := t.calculateChannelStatistics(v.ChannelId, statistics)
		v.TransferTxs = count
		v.TransferTxsValue = value.Round(constant.DefaultValuePrecision).String()
	}

	return nil
}

func (t *ChannelTask) calculateChannelStatistics(channelId string, statistics []*dto.ChannelStatisticsAggrDTO) (int64, decimal.Decimal) {
	var txsCount int64 = 0
	var txsValue = decimal.Zero

	for _, v := range statistics {
		if channelId == v.ChannelId {
			valueDecimal := t.calculateValue(v.TxsAmount, v.BaseDenom)
			txsCount += v.TxsCount
			txsValue = txsValue.Add(valueDecimal)

			chainA, _, chainB, _, _ := t.parseChannelId(channelId)
			t.chainTxsMap[chainA] += v.TxsCount
			t.chainTxsMap[chainB] += v.TxsCount
			d, ok := t.chainTxsValueMap[chainA]
			if ok {
				t.chainTxsValueMap[chainA] = d.Add(valueDecimal)
			} else {
				t.chainTxsValueMap[chainA] = valueDecimal
			}

			d, ok = t.chainTxsValueMap[chainB]
			if ok {
				t.chainTxsValueMap[chainB] = d.Add(valueDecimal)
			} else {
				t.chainTxsValueMap[chainB] = valueDecimal
			}
		}
	}

	return txsCount, txsValue
}

func (t *ChannelTask) calculateValue(amount float64, baseDenom string) decimal.Decimal {
	denom, ok := t.baseDenomMap[baseDenom]
	if !ok || denom.CoinId == "" {
		return decimal.Zero
	}

	price, err := tokenPriceRepo.Get(denom.CoinId)
	if err != nil {
		logrus.Errorf("task %s calculateValue error, %v", t.Name(), err)
		return decimal.Zero
	}

	value := decimal.NewFromFloat(amount).Div(decimal.NewFromFloat(math.Pow10(denom.Scale))).
		Mul(decimal.NewFromFloat(price))

	return value
}

func (t *ChannelTask) todayStatistics() error {
	logrus.Infof("task %s exec today statistics", t.Name())
	startTime, endTime := todayUnix()
	segments := []*segment{
		{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if err := channelStatisticsTask.deal(segments, opUpdate); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}

	return nil
}

func (t *ChannelTask) yesterdayStatistics() error {
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
	if err := channelStatisticsTask.deal(segments, opUpdate); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}

	_ = statisticsCheckRepo.Incr(t.Name(), mmdd)
	return nil
}
