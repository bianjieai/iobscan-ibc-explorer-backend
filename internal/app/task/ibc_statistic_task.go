package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
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
	if err := t.updateChannel24h(); err != nil {
		return -1
	}
	if err := t.updateChannelInfo(); err != nil {
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

func (t *IbcStatisticCronTask) updateChannel24h() error {
	//获取最近24小时前的时间戳
	startTime := time.Now().Unix() - 24*3600
	channelDtos, err := ibcTxRepo.Aggr24hActiveChannelTxs(startTime)
	if err != nil {
		return err
	}
	setMap := make(map[string]struct{}, len(channelDtos))
	for _, val := range channelDtos {
		channelIdPrefix := fmt.Sprintf("%s|%s", val.ScChainId, val.ScChannel)
		channelIdEndwith := fmt.Sprintf("%s|%s", val.DcChainId, val.DcChannel)
		_, existPrefix := setMap[channelIdPrefix]
		_, existEndWith := setMap[channelIdEndwith]
		if !existEndWith && !existPrefix {
			setMap[channelIdPrefix] = struct{}{}
			setMap[channelIdEndwith] = struct{}{}
		}
	}
	count := int64(len(setMap))
	if err := statisticsRepo.UpdateOne(constant.Channel24hStatisticName, count); err != nil {
		return err
	}
	return nil
}
