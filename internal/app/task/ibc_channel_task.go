package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor"
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

func (t *ChannelTask) Cron() string {
	return ThreeMinute
}

func (t *ChannelTask) ExpireTime() time.Duration {
	return 3*time.Minute - 1*time.Second
}

func (t *ChannelTask) Run() int {
	monitor.SetCronTaskStatusMetricValue(t.Name(), -1)
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

	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}
	t.baseDenomMap = baseDenomList.ConvertToMap()

	statistics, err := t.channelStatistics()
	if err != nil {
		return -1
	}

	t.setTransferTxs(existedChannelList, newChannelList, statistics) // 计算txs和交易价值，同时更新ibc_channel_statistics

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
		if err = chainRepo.UpdateTransferTxs(chainId, txs, txsValue); err != nil {
			logrus.Errorf("task %s update chain %s error, %v", t.Name(), chainId, err)
		}
	}
	return 1
}

func (t *ChannelTask) analyzeChainConfig() error {
	confList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s analyzeChainConfig error, %v", t.Name(), err)
		return err
	}

	var channelIds, channelMirrorIds []string
	channelStatusMap := make(map[string]entity.ChannelStatus)

	var chainA, channelA, chainB, channelB string
	for _, v := range confList {
		chainA = v.ChainId
		for _, info := range v.IbcInfo {
			chainB = info.ChainId
			for _, p := range info.Paths {
				channelA = p.ChannelId
				channelB = p.Counterparty.ChannelId
				channelId, mirrorChannelId := generateChannelId(chainA, channelA, chainB, channelB)

				if utils.InArray(channelMirrorIds, channelId) || utils.InArray(channelMirrorIds, mirrorChannelId) { // 已经存在
					continue
				}

				channelMirrorIds = append(channelMirrorIds, channelId, mirrorChannelId)
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

func generateChannelId(chainA, channelA, chainB, channelB string) (string, string) {
	return fmt.Sprintf("%s|%s|%s|%s", chainA, channelA, chainB, channelB), fmt.Sprintf("%s|%s|%s|%s", chainB, channelB, chainA, channelA)
}

func (t *ChannelTask) parseChannelId(channelId string) (chainA, channelA, chainB, channelB string, err error) {
	split := strings.Split(channelId, "|")
	if len(split) != 4 {
		logrus.Errorf("task %s parseChannelId error, %v", t.Name(), err)
		return "", "", "", "", fmt.Errorf("channel id format error")
	}
	return split[0], split[1], split[2], split[3], nil
}

func (t *ChannelTask) getChannelMirrorId(channelId string) string {
	chainA, channelA, chainB, channelB, err := t.parseChannelId(channelId)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s|%s|%s|%s", chainB, channelB, chainA, channelA)
}

func (t *ChannelTask) channelEqual(channelId1, channelId2 string) bool {
	if channelId1 == channelId2 {
		return true
	}

	chainA1, channelA1, chainB1, channelB1, err := t.parseChannelId(channelId1)
	if err != nil {
		return false
	}

	chainA2, channelA2, chainB2, channelB2, err := t.parseChannelId(channelId1)
	if err != nil {
		return false
	}

	if chainA1 == chainA2 && channelA1 == channelA2 && chainB1 == chainB2 && channelB1 == channelB2 {
		return true
	}

	if chainA1 == chainB2 && channelA1 == channelB2 && chainB1 == chainA2 && channelB1 == channelA2 {
		return true
	}

	return false
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
			if t.channelEqual(e, v) {
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
			ChannelId:            v,
			ChainA:               chainA,
			ChainB:               chainB,
			ChannelA:             channelA,
			ChannelB:             channelB,
			Status:               entity.ChannelStatusOpened,
			OperatingPeriod:      0,
			LatestSettlementTime: 0,
			Relayers:             0,
			TransferTxs:          0,
			TransferTxsValue:     "",
			CreateAt:             time.Now().Unix(),
			UpdateAt:             time.Now().Unix(),
		})
	}

	return existedChannelList, newChannelList, nil
}

func (t *ChannelTask) setLatestSettlementTime(existedChannelList entity.IBCChannelList, newChannelList entity.IBCChannelList) error {
	// todo
	for _, v := range newChannelList {
		// 查询,初始的LatestSettlementTime 为channel的 open confirm 时间
		// channel open confirm 时间的获取当前没有实现，先设为0
		v.LatestSettlementTime = 0
	}

	for _, v := range existedChannelList {
		// 之前没有设置open 时间或者之前关闭了，现在重新打开channel
		if (v.LatestSettlementTime == 0) || (v.Status == entity.ChannelStatusClosed && t.channelStatusMap[v.ChannelId] == entity.ChannelStatusOpened) {
			// 查询

		}
	}
	return nil
}

func (t *ChannelTask) setStatusAndOperatingPeriod(existedChannelList entity.IBCChannelList, newChannelList entity.IBCChannelList) {
	set := func(list entity.IBCChannelList) {
		for _, v := range list {
			status, ok := t.channelStatusMap[v.ChannelId]
			if ok {
				v.Status = status
			}

			if v.LatestSettlementTime != 0 {
				now := time.Now().Unix()
				v.OperatingPeriod += now - v.LatestSettlementTime
				v.LatestSettlementTime = now
			}
		}
	}

	set(existedChannelList)
	set(newChannelList)
}

func (t *ChannelTask) channelStatistics() ([]*dto.ChannelStatisticsDTO, error) {
	channelTxs, err := ibcTxRepo.AggrIBCChannelTxs()
	if err != nil {
		logrus.Errorf("task %s channelStatistics error, %v", t.Name(), err)
		return nil, err
	}

	historyChannelTxs, err := ibcTxRepo.AggrIBCChannelHistoryTxs()
	if err != nil {
		logrus.Errorf("task %s channelStatistics error, %v", t.Name(), err)
		return nil, err
	}

	// channel 不区分交易方向，将属于一个channel的 A->B, B->A 的交易信息整合一下
	integration := func(cl []*dto.ChannelStatisticsDTO, aggr []*dto.AggrIBCChannelTxsDTO) []*dto.ChannelStatisticsDTO {
		for _, v := range aggr {
			isExisted := false
			ChannelId, MirrorChannelId := generateChannelId(v.ScChainId, v.ScChannel, v.DcChainId, v.DcChannel)
			for _, c := range cl {
				if (t.channelEqual(ChannelId, c.ChannelId) || t.channelEqual(MirrorChannelId, c.ChannelId)) && v.BaseDenom == c.BaseDenom { // 同一个channel
					c.TxsCount += v.Count
					c.TxsAmount = c.TxsAmount.Add(decimal.NewFromFloat(v.Amount))
					isExisted = true
					break
				}
			}

			if !isExisted {
				cl = append(cl, &dto.ChannelStatisticsDTO{
					ChannelId:       ChannelId,
					MirrorChannelId: MirrorChannelId,
					BaseDenom:       v.BaseDenom,
					TxsCount:        v.Count,
					TxsAmount:       decimal.NewFromFloat(v.Amount),
				})
			}
		}
		return cl
	}

	var cslist []*dto.ChannelStatisticsDTO
	cslist = integration(cslist, channelTxs)
	cslist = integration(cslist, historyChannelTxs)

	return cslist, nil
}

func (t *ChannelTask) setTransferTxs(existedChannelList entity.IBCChannelList, newChannelList entity.IBCChannelList, statistics []*dto.ChannelStatisticsDTO) {
	for _, v := range existedChannelList {
		count, value, err := t.calculateChannelStatistics(v.ChannelId, statistics)
		if err != nil {
			continue
		}

		v.TransferTxs = count
		v.TransferTxsValue = value.Round(constant.DefaultValuePrecision).String()
	}

	for _, v := range newChannelList {
		count, value, err := t.calculateChannelStatistics(v.ChannelId, statistics)
		if err != nil {
			continue
		}

		v.TransferTxs = count
		v.TransferTxsValue = value.Round(constant.DefaultValuePrecision).String()
	}
}

func (t *ChannelTask) calculateChannelStatistics(channelId string, statistics []*dto.ChannelStatisticsDTO) (int64, decimal.Decimal, error) {
	var ibcList []*entity.IBCChannelStatistics
	var txsCount int64 = 0
	var txsValue = decimal.Zero
	if t.chainTxsMap == nil {
		t.chainTxsMap = make(map[string]int64)
	}
	if t.chainTxsValueMap == nil {
		t.chainTxsValueMap = make(map[string]decimal.Decimal)
	}

	for _, v := range statistics {
		if channelId == v.ChannelId || channelId == v.MirrorChannelId {
			valueDecimal := t.calculateValue(v.TxsAmount, v.BaseDenom)
			ibcList = append(ibcList, &entity.IBCChannelStatistics{
				ChannelId:          channelId,
				TransferBaseDenom:  v.BaseDenom,
				TransferAmount:     v.TxsAmount.String(),
				TransferTotalValue: valueDecimal.Round(constant.DefaultValuePrecision).String(),
				CreateAt:           time.Now().Unix(),
				UpdateAt:           time.Now().Unix(),
			})

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

	err := channelStatisticsRepo.BatchSwap(channelId, ibcList)
	if err != nil {
		logrus.Errorf("task %s calculateChannelStatistics error, %v", t.Name(), err)
		return 0, decimal.Zero, err
	}

	return txsCount, txsValue, nil
}

func (t *ChannelTask) calculateValue(amount decimal.Decimal, baseDenom string) decimal.Decimal {
	var coinId string
	var scale int
	for _, v := range t.baseDenomMap {
		if v.Denom == baseDenom {
			coinId = v.CoinId
			scale = v.Scale
			break
		}
	}

	if coinId == "" {
		return decimal.Zero
	}

	price, err := tokenPriceRepo.Get(coinId)
	if err != nil {
		logrus.Errorf("task %s calculateValue error, %v", t.Name(), err)
		return decimal.Zero
	}

	value := amount.Div(decimal.NewFromFloat(math.Pow10(scale))).
		Mul(decimal.NewFromFloat(price))

	return value
}
