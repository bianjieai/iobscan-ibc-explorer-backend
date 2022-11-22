package task

import (
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
)

// yesterdayUnix 获取昨日第一秒和最后一秒的时间戳
func yesterdayUnix() (int64, int64) {
	date := time.Now().AddDate(0, 0, -1)
	startUnix := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).Unix()
	endUnix := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 59, time.Local).Unix()
	return startUnix, endUnix
}

//==========active accounts======
func caculateActiveAddrsOfChains() {
	startTimeA := time.Now().Unix()
	defer func() {
		logrus.Infof("cronjob execute caculateActiveAddrsOfChains finished, time use %d(s)", time.Now().Unix()-startTimeA)
	}()
	logrus.Infof("cronjob execute caculateActiveAddrsOfChains start...")
	configList, err := chainConfigRepo.FindAllChains()
	if err != nil {
		logrus.Errorf("find chain_config error, %v", err)
		return
	}
	//获取relayer地址
	relayerAddrs, err := relayerAddrCache.FindAll()
	if err != nil {
		logrus.Errorf("find cache relayer error, %v", err)
		return
	}
	relayerMap := make(map[string]struct{}, len(relayerAddrs))
	for _, val := range relayerAddrs {
		relayerMap[val] = struct{}{}
	}

	//获取昨天的时间
	startTime, endTime := yesterdayUnix()
	mapChainAddrs := make(map[string][]string, len(configList))
	for _, chainCfg := range configList {
		res, err := txRepo.GetActiveAccountsOfDay(chainCfg.ChainName, startTime, endTime+1)
		if err != nil {
			logrus.Error("faild get active accounts of day ,err " + err.Error())
			continue
		}
		for _, val := range res {
			//过滤掉relayer地址
			if _, exist := relayerMap[val.Address]; exist {
				continue
			}
			if strings.HasPrefix(val.Address, chainCfg.AddrPrefix) {
				mapChainAddrs[chainCfg.ChainName] = append(mapChainAddrs[chainCfg.ChainName], val.Address)
			}
		}
	}

	dailyDate := time.Now().AddDate(0, 0, -1)
	cache.GetRedisClient().Set(cache.DailyAccountsDate, utils.FmtTime(dailyDate, utils.DateFmtYYYYMMDD), -1)
	if err := statisticsRepo.UpdateOneData(constant.AccountsDailyStatisticName, string(utils.MarshalJsonIgnoreErr(mapChainAddrs))); err != nil && !qmgo.IsDup(err) {
		logrus.Errorf("update statistic data error, %v", err)
	}
}
