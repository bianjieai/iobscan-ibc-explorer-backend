package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"time"
)

type IbcStatisticCronTask struct {
}

func (t *IbcStatisticCronTask) Name() string {
	return "ibc_statistic_task"
}
func (t *IbcStatisticCronTask) Cron() int {
	if taskConf.CronTimeStatisticTask > 0 {
		return taskConf.CronTimeStatisticTask
	}
	return EveryMinute
}

func (t *IbcStatisticCronTask) Run() int {
	if err := t.updateChannelAndChains24h(); err != nil {
		return -1
	}
	if err := t.updateChannelInfo(); err != nil {
		return -1
	}

	if err := t.handleDenomIncre(constant.BaseDenomAllStatisticName, denomRepo.BasedDenomCount); err != nil {
		return -1
	}

	if err := t.handleDenomIncre(constant.DenomAllStatisticName, denomRepo.Count); err != nil {
		return -1
	}

	chainsAll, err := chainConfigRepo.Count()
	if err != nil {
		return -1
	}
	if err := statisticsRepo.UpdateOne(constant.ChainsAllStatisticName, chainsAll); err != nil {
		return -1
	}
	return 1
}

func (t *IbcStatisticCronTask) updateChannelInfo() error {
	channelOpen, err := channelRepo.CountStatus(entity.ChannelStatusOpened)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return err
	}
	channelClose, err := channelRepo.CountStatus(entity.ChannelStatusClosed)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return err
	}
	channelAll := channelOpen + channelClose
	if err := statisticsRepo.UpdateOne(constant.ChannelOpenStatisticName, channelOpen); err != nil {
		return err
	}
	if err := statisticsRepo.UpdateOne(constant.ChannelCloseStatisticName, channelClose); err != nil {
		return err
	}
	if err := statisticsRepo.UpdateOne(constant.ChannelAllStatisticName, channelAll); err != nil {
		return err
	}
	return nil
}

func (t *IbcStatisticCronTask) updateChannelAndChains24h() error {
	//获取最近24小时前的时间戳
	startTime := time.Now().Unix() - 24*3600
	channelDtos, err := ibcTxRepo.Aggr24hActiveChannelTxs(startTime)
	if err != nil {
		return err
	}
	setChannelMap := make(map[string]struct{}, len(channelDtos))
	chainIdMap := make(map[string]struct{}, len(channelDtos))
	count := int64(0)
	for _, val := range channelDtos {
		chainIdMap[val.ScChainId] = struct{}{}
		chainIdMap[val.DcChainId] = struct{}{}
		channelIdPrefix := fmt.Sprintf("%s|%s", val.ScChainId, val.ScChannel)
		channelIdEndwith := fmt.Sprintf("%s|%s", val.DcChainId, val.DcChannel)
		_, existPrefix := setChannelMap[channelIdPrefix]
		_, existEndWith := setChannelMap[channelIdEndwith]
		if !existEndWith && !existPrefix {
			setChannelMap[channelIdPrefix] = struct{}{}
			setChannelMap[channelIdEndwith] = struct{}{}
			count++
		}
	}
	if err := statisticsRepo.UpdateOne(constant.Channel24hStatisticName, count); err != nil {
		return err
	}

	chains24hAll := len(chainIdMap)
	if err := statisticsRepo.UpdateOne(constant.Chains24hStatisticName, int64(chains24hAll)); err != nil {
		return err
	}
	return nil
}

func (t *IbcStatisticCronTask) handleDenomIncre(statisticName string, call func(createAt int64, record bool) (int64, error)) error {
	statisticData, err := statisticsRepo.FindOne(statisticName)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		return err
	}
	if statisticData.StatisticsName == "" {
		statisticData = entity.IbcStatistic{
			StatisticsName: statisticName,
			StatisticsInfo: "",
			Count:          0,
		}
		denomAllCnt, err := call(0, false)
		if err != nil {
			return err
		}
		statisticData.Count = denomAllCnt

		latestCreateAt, err := denomRepo.LatestCreateAt()
		if err != nil {
			return err
		}
		currentDenomCnt, err := call(latestCreateAt, true)
		if err != nil {
			return err
		}
		statisticData.StatisticsInfo = string(utils.MarshalJsonIgnoreErr(IncreInfo{Count: currentDenomCnt, CreateAt: latestCreateAt}))
		statisticData.CreateAt = time.Now().Unix()
		if err := statisticsRepo.UpdateOneIncre(statisticData); err != nil {
			return err
		}
		return nil
	} else if statisticData.StatisticsInfo == "" {
		denomAllCnt, err := call(0, false)
		if err != nil {
			return err
		}
		statisticData.Count = denomAllCnt
		latestCreateAt, err := denomRepo.LatestCreateAt()
		if err != nil {
			return err
		}
		currentDenomCnt, err := call(latestCreateAt, true)
		if err != nil {
			return err
		}
		statisticData.StatisticsInfo = string(utils.MarshalJsonIgnoreErr(IncreInfo{Count: currentDenomCnt, CreateAt: latestCreateAt}))
		statisticData.CreateAt = time.Now().Unix()
		if err := statisticsRepo.UpdateOneIncre(statisticData); err != nil {
			return err
		}
		return nil
	}

	var increData IncreInfo
	utils.UnmarshalJsonIgnoreErr([]byte(statisticData.StatisticsInfo), &increData)

	denomAllCnt, err := call(increData.CreateAt, false)
	if err != nil {
		return err
	}
	if denomAllCnt > increData.Count {
		increValue := denomAllCnt - increData.Count
		statisticData.Count = statisticData.Count + increValue
		latestCreateAt, err := denomRepo.LatestCreateAt()
		if err != nil {
			return err
		}
		currentDenomCnt, err := call(latestCreateAt, true)
		if err != nil {
			return err
		}
		statisticData.StatisticsInfo = string(utils.MarshalJsonIgnoreErr(IncreInfo{Count: currentDenomCnt, CreateAt: latestCreateAt}))
		statisticData.CreateAt = time.Now().Unix()
		if err := statisticsRepo.UpdateOneIncre(statisticData); err != nil {
			return err
		}
	}
	return nil
}

type IncreInfo struct {
	Count    int64 `json:"count" bson:"count"`
	CreateAt int64 `json:"create_at" bson:"create_at"`
}
