package cache

import (
	"fmt"
	"strconv"
)

type StatisticsCheckCacheRepo struct {
}

func (repo *StatisticsCheckCacheRepo) GetIncr(task, date string) (int, error) {
	key := fmt.Sprintf(statisticsCheck, task, date)
	incr, err := rc.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(incr)
}

func (repo *StatisticsCheckCacheRepo) Incr(task, date string) (int64, error) {
	key := fmt.Sprintf(statisticsCheck, task, date)
	incr, err := rc.Incr(key)
	if err != nil {
		return 0, err
	}

	if incr == 1 { // 第一次设置过期时间
		rc.Expire(key, oneDay)
	}
	return incr, nil
}
