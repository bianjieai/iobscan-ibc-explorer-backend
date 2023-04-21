package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/sirupsen/logrus"
	"time"
)

type IBCCacheTokenTraceStatisticTask struct {
}

func (t *IBCCacheTokenTraceStatisticTask) Name() string {
	return "ibc_cache_token_trace_statistics_task"
}

func (t *IBCCacheTokenTraceStatisticTask) Cron() string {
	if taskConf.IBCCacheTokenTraceStatisticTask != "" {
		return taskConf.IBCCacheTokenTraceStatisticTask
	}
	return "0 */3 * * * ?"
}

func (t *IBCCacheTokenTraceStatisticTask) BeforeHook() error {
	return nil
}

func (t *IBCCacheTokenTraceStatisticTask) Run() {
	res, err := ibcTokenTraceStatisticsRepo.Aggr()
	if err != nil {
		logrus.Errorf("%s Aggr fail, err:%s", t.Name(), err.Error())
		TaskMetricMap.Store(t.Name(), float64(-1))
		return
	}
	tokenReceiveTxs := make(map[string]int64)
	for _, v := range res {
		tokenReceiveTxs[v.Chain+"/"+v.Denom] = v.ReceiveTxs
	}
	_, err = cache.GetRedisClient().MarshalHSet(cache.TokenReceiveTxs, tokenReceiveTxs)
	if err != nil {
		logrus.Errorf("%s cache fail, err:%s", t.Name(), err.Error())
		TaskMetricMap.Store(t.Name(), float64(-1))
		return
	}
	cache.GetRedisClient().Expire(cache.TokenReceiveTxs, time.Minute*6)
	TaskMetricMap.Store(t.Name(), float64(1))
	return
}
