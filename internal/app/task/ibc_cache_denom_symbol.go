package task

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/sirupsen/logrus"
)

type IBCCacheDenomSymbolTask struct {
}

func (t *IBCCacheDenomSymbolTask) Name() string {
	return "ibc_cache_denom_symbol_task"
}

func (t *IBCCacheDenomSymbolTask) Cron() string {
	if taskConf.IBCCacheDenomSymbolTask != "" {
		return taskConf.IBCCacheDenomSymbolTask
	}
	return "0 */3 * * * ?"
}

func (t *IBCCacheDenomSymbolTask) BeforeHook() error {
	return nil
}

func (t *IBCCacheDenomSymbolTask) Run() {
	res, err := ibcDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("%s find all ibc denom fail, err:%s", t.Name(), err.Error())
		TaskMetricMap.Store(t.Name(), float64(-1))
		return
	}
	tokenSymbol := make(map[string]string)
	for _, v := range res {
		tokenSymbol[v.Chain+"/"+v.Denom] = v.Symbol
	}
	_, err = cache.GetRedisClient().MarshalHSet(cache.DenomSymbol, tokenSymbol)
	if err != nil {
		logrus.Errorf("%s cache fail, err:%s", t.Name(), err.Error())
		TaskMetricMap.Store(t.Name(), float64(-1))
		return
	}
	TaskMetricMap.Store(t.Name(), float64(1))
	return
}
