package task

import (
	"fmt"
	"math"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type ChannelStatisticsTask struct {
	segmentMinTime              int64
	segmentChannelStatisticsMap map[string][]*dto.ChannelStatisticsDTO
}

var channelStatisticsTask ChannelStatisticsTask

func (t *ChannelStatisticsTask) Name() string {
	return "ibc_channel_statistics_task"
}

func (t *ChannelStatisticsTask) Switch() bool {
	return global.Config.Task.SwitchIbcChannelStatisticsTask
}

func ChannelIncrementStatistics(segments []*segment) error {
	return channelStatisticsTask.deal(segments, false)
}

func (t *ChannelStatisticsTask) Run() int {
	t.segmentMinTime = math.MaxInt64
	t.segmentChannelStatisticsMap = make(map[string][]*dto.ChannelStatisticsDTO)

	if err := channelStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s CreateNew err, %v", t.Name(), err)
		return -1
	}

	segments, err := getTxTimeSegment(false, segmentStepLatest)
	if err != nil {
		logrus.Errorf("task %s getSegment err, %v", t.Name(), err)
		return -1
	}

	historySegments, err := getTxTimeSegment(true, segmentStepHistory)
	if err != nil {
		logrus.Errorf("task %s getHistorySegment err, %v", t.Name(), err)
		return -1
	}

	t.segmentMinTime = segments[0].StartTime
	// 优先处理历史分段
	logrus.Infof("task %s deal history segment total: %d", t.Name(), len(historySegments))
	if err = t.dealHistory(historySegments); err != nil {
		logrus.Errorf("task %s dealHistory err, %v", t.Name(), err)
		return -1
	}

	logrus.Infof("task %s deal segment total: %d", t.Name(), len(segments))
	if err = t.deal(segments, true); err != nil {
		logrus.Errorf("task %s deal err, %v", t.Name(), err)
		return -1
	}

	if err = channelStatisticsRepo.SwitchColl(); err != nil {
		logrus.Errorf("task %s SwitchColl err, %v", t.Name(), err)
		return -1
	}

	return 1
}

// dealHistory 处理历史记录，针对ex_ibc_tx
func (t *ChannelStatisticsTask) dealHistory(segments []*segment) error {
	for _, v := range segments {
		logrus.Infof("task %s dealHistory segment [%d, %d]", t.Name(), v.StartTime, v.EndTime)
		txs, err := ibcTxRepo.AggrIBCChannelHistoryTxs(v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s AggrIBCChannelHistoryTxs err, %v", t.Name(), err)
			return err
		}

		if len(txs) == 0 {
			continue
		}

		channelStatisticsAggr := t.aggr(txs)
		if v.StartTime >= t.segmentMinTime {
			// 将新老表重叠的分段数据记录到map
			t.segmentChannelStatisticsMap[fmt.Sprintf("%d-%d", v.StartTime, v.EndTime)] = channelStatisticsAggr
		}
		if err = t.saveData(channelStatisticsAggr, v, true, true); err != nil {
			logrus.Errorf("task %s dealHistory saveData err, %v", t.Name(), err)
		}
	}
	return nil
}

// deal 处理最新的记录，针对ex_ibc_tx_latest
func (t *ChannelStatisticsTask) deal(segments []*segment, fullStatistics bool) error {
	for _, v := range segments {
		logrus.Infof("task %s deal segment [%d, %d]", t.Name(), v.StartTime, v.EndTime)
		txs, err := ibcTxRepo.AggrIBCChannelTxs(v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s AggrIBCChannelTxs err, %v", t.Name(), err)
			return err
		}

		if len(txs) == 0 {
			continue
		}

		channelStatisticsAggr := t.aggr(txs)
		if fullStatistics {
			channelStatisticsAggr = t.integrationStatisticsData(channelStatisticsAggr, v)
		}
		if err = t.saveData(channelStatisticsAggr, v, false, fullStatistics); err != nil {
			logrus.Errorf("task %s deal saveData err, %v", t.Name(), err)
		}
	}
	return nil
}

func (t *ChannelStatisticsTask) integrationStatisticsData(aggrRes []*dto.ChannelStatisticsDTO, seg *segment) []*dto.ChannelStatisticsDTO {
	// 新表中的段与历史表重和，需要整合数据
	hirtoryAggrRes, ok := t.segmentChannelStatisticsMap[fmt.Sprintf("%d-%d", seg.StartTime, seg.EndTime)]
	if !ok {
		return aggrRes
	}

	integrationDataMap := make(map[string]*dto.ChannelStatisticsDTO, len(aggrRes))
	for _, v := range aggrRes {
		key := fmt.Sprintf("%s%s%s%d", v.ChannelId, v.BaseDenom, v.BaseDenomChain, v.Status)
		if data, ok := integrationDataMap[key]; !ok {
			integrationDataMap[key] = v
		} else {
			data.TxsCount += v.TxsCount
			data.TxsAmount = data.TxsAmount.Add(v.TxsAmount)
		}
	}

	for _, v := range hirtoryAggrRes {
		key := fmt.Sprintf("%s%s%s%d", v.ChannelId, v.BaseDenom, v.BaseDenomChain, v.Status)
		if data, ok := integrationDataMap[key]; !ok {
			integrationDataMap[key] = v
		} else {
			data.TxsCount += v.TxsCount
			data.TxsAmount = data.TxsAmount.Add(v.TxsAmount)
		}
	}

	integrationDataList := make([]*dto.ChannelStatisticsDTO, 0, len(integrationDataMap))
	for _, v := range integrationDataMap {
		integrationDataList = append(integrationDataList, v)
	}
	return integrationDataList
}

func (t *ChannelStatisticsTask) aggr(txs []*dto.AggrIBCChannelTxsDTO) []*dto.ChannelStatisticsDTO {
	var cl []*dto.ChannelStatisticsDTO
	for _, v := range txs {
		isExisted := false
		ChannelId := generateChannelId(v.ScChain, v.ScChannel, v.DcChain, v.DcChannel)
		for _, c := range cl {
			if c.ChannelId == ChannelId && v.BaseDenom == c.BaseDenom && v.BaseDenomChain == c.BaseDenomChain &&
				v.Status == c.Status { // 同一个channel
				c.TxsCount += v.Count
				c.TxsAmount = c.TxsAmount.Add(decimal.NewFromFloat(v.Amount))
				isExisted = true
				break
			}
		}

		if !isExisted {
			cl = append(cl, &dto.ChannelStatisticsDTO{
				ChannelId:      ChannelId,
				BaseDenom:      v.BaseDenom,
				BaseDenomChain: v.BaseDenomChain,
				TxsCount:       v.Count,
				TxsAmount:      decimal.NewFromFloat(v.Amount),
				Status:         v.Status,
			})
		}
	}

	return cl
}

func (t *ChannelStatisticsTask) saveData(dtoList []*dto.ChannelStatisticsDTO, seg *segment, targetHistory, fullStatistics bool) error {
	var statistics = make([]*entity.IBCChannelStatistics, 0, len(dtoList))
	for _, v := range dtoList {
		statistics = append(statistics, &entity.IBCChannelStatistics{
			ChannelId:        v.ChannelId,
			BaseDenom:        v.BaseDenom,
			BaseDenomChain:   v.BaseDenomChain,
			TransferTxs:      v.TxsCount,
			TransferAmount:   v.TxsAmount.String(),
			Status:           entity.IbcTxStatus(v.Status),
			SegmentStartTime: seg.StartTime,
			SegmentEndTime:   seg.EndTime,
			CreateAt:         time.Now().Unix(),
			UpdateAt:         time.Now().Unix(),
		})
	}

	var err error
	if fullStatistics {
		if targetHistory {
			err = channelStatisticsRepo.BatchInsertToNew(statistics)
		} else {
			err = channelStatisticsRepo.BatchSwapNew(seg.StartTime, seg.EndTime, statistics)
		}
	} else {
		err = channelStatisticsRepo.BatchSwap(seg.StartTime, seg.EndTime, statistics)
	}

	return err
}
