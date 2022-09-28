package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor/metrics"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

var (
	cronTaskStatusMetric  metrics.Guage
	lcdConnectStatsMetric metrics.Guage
	redisStatusMetric     metrics.Guage
	TagName               = "taskname"
	ChainTag              = "chain_id"

	chainConfigRepo   repository.IChainConfigRepo   = new(repository.ChainConfigRepo)
	chainRegistryRepo repository.IChainRegistryRepo = new(repository.ChainRegistryRepo)
)

const (
	v1beta1         = "v1beta1"
	v1              = "v1"
	v1beta1Channels = "/ibc/core/channel/v1beta1/channels"
	v1Channels      = "/ibc/core/channel/v1/channels"
	v1beta1Params   = "/cosmos/staking/v1beta1/params"
	v1Params        = "/cosmos/staking/v1/params"

	apiChannels    = "/ibc/core/channel/%s/channels?pagination.offset=OFFSET&pagination.limit=LIMIT&pagination.count_total=true"
	apiClientState = "/ibc/core/channel/%s/channels/CHANNEL/ports/PORT/client_state"
	//apiBalances    = "/cosmos/bank/%s/balances/{address}"
	//apiParams      = "/cosmos/{module}/%s/params"
	//apiSupply      = "/cosmos/bank/%s/supply"
)

func NewMetricCronWorkStatus() metrics.Guage {
	syncWorkStatusMetric := metrics.NewGuage(
		"ibc_explorer_backend",
		"",
		"cron_task_status",
		"ibc_explorer_backend cron task working status (1:Normal  -1:UNormal)",
		[]string{TagName},
	)
	syncWorkStatus, _ := metrics.CovertGuage(syncWorkStatusMetric)
	return syncWorkStatus
}

func NewMetricRedisStatus() metrics.Guage {
	redisNodeStatusMetric := metrics.NewGuage(
		"ibc_explorer_backend",
		"redis",
		"connection_status",
		"ibc_explorer_backend  node connection status of redis service (1:Normal  -1:UNormal)",
		nil,
	)
	redisStatus, _ := metrics.CovertGuage(redisNodeStatusMetric)
	return redisStatus
}

func NewMetricLcdStatus() metrics.Guage {
	lcdConnectionStatusMetric := metrics.NewGuage(
		"ibc_explorer_backend",
		"lcd",
		"connection_status",
		"ibc_explorer_backend  lcd connection status of blockchain (1:Normal  -1:UNormal)",
		[]string{ChainTag},
	)
	connectionStatus, _ := metrics.CovertGuage(lcdConnectionStatusMetric)
	return connectionStatus
}

func SetCronTaskStatusMetricValue(taskName string, value float64) {
	if cronTaskStatusMetric != nil {
		cronTaskStatusMetric.With(TagName, taskName).Set(value)
	}
}

func lcdConnectionStatus(quit chan bool) {
	for {
		t := time.NewTimer(time.Duration(15) * time.Second)
		select {
		case <-t.C:
			chainCfgs, err := chainConfigRepo.FindAllChainInfs()
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			for _, val := range chainCfgs {
				if checkLcd(val.Lcd) {
					lcdConnectStatsMetric.With(ChainTag, val.ChainId).Set(float64(1))
				} else {
					if switchLcd(val) {
						lcdConnectStatsMetric.With(ChainTag, val.ChainId).Set(float64(1))
					} else {
						lcdConnectStatsMetric.With(ChainTag, val.ChainId).Set(float64(-1))
						logrus.Errorf("monitor chain %s lcd is unavailable", val.ChainId)
					}
				}
			}

		case <-quit:
			logrus.Debug("quit signal recv  lcdConnectionStatus")
			return

		}
	}
}

// checkLcd If lcd is ok, return true. Else return false
func checkLcd(lcd string) bool {
	_, err := utils.HttpGet(fmt.Sprintf("%s/node_info", lcd))
	if err != nil {
		_, err = utils.HttpGet(fmt.Sprintf("%s/blocks/latest", lcd))
		if err != nil {
			return false
		}
	}

	return true
}

// switchLcd If Switch lcd succeeded, return true. Else return false
func switchLcd(chainConf *entity.ChainConfig) bool {
	chainRegistry, err := chainRegistryRepo.FindOne(chainConf.ChainId)
	if err != nil {
		logrus.Errorf("lcd monitor error: %v", err)
		return false
	}

	bz, err := utils.HttpGet(chainRegistry.ChainJsonUrl)
	if err != nil {
		logrus.Errorf("lcd monitor get chain json error: %v", err)
		return false
	}

	var chainRegisterResp vo.ChainRegisterResp
	_ = json.Unmarshal(bz, &chainRegisterResp)
	for _, v := range chainRegisterResp.Apis.Rest {
		if !checkLcd(v.Address) {
			continue
		}

		chainConf.Lcd = v.Address
		if _, err := utils.HttpGet(fmt.Sprintf("%s%s", v.Address, v1beta1Channels)); err == nil ||
			!strings.Contains(err.Error(), "501 Not Implemented") {
			chainConf.LcdApiPath.ChannelsPath = fmt.Sprintf(apiChannels, v1beta1)
			chainConf.LcdApiPath.ClientStatePath = fmt.Sprintf(apiClientState, v1beta1)
		} else {
			chainConf.LcdApiPath.ChannelsPath = fmt.Sprintf(apiChannels, v1)
			chainConf.LcdApiPath.ClientStatePath = fmt.Sprintf(apiClientState, v1)
		}

		//if _, err := utils.HttpGet(fmt.Sprintf("%s%s", v.Address, v1beta1Params)); err == nil || !strings.Contains(err.Error(), "501 Not Implemented") {
		//	chainConf.LcdApiPath.BalancesPath = fmt.Sprintf(apiBalances, v1beta1)
		//	chainConf.LcdApiPath.SupplyPath = fmt.Sprintf(apiSupply, v1beta1)
		//	chainConf.LcdApiPath.ParamsPath = fmt.Sprintf(apiParams, v1beta1)
		//} else {
		//	chainConf.LcdApiPath.BalancesPath = fmt.Sprintf(apiBalances, v1)
		//	chainConf.LcdApiPath.SupplyPath = fmt.Sprintf(apiSupply, v1)
		//	chainConf.LcdApiPath.ParamsPath = fmt.Sprintf(apiParams, v1)
		//}

		if err = chainConfigRepo.UpdateLcdApi(chainConf); err != nil {
			logrus.Error("switch lcd error: %v", err)
			return false
		} else {
			return true
		}
	}

	return false
}

func redisClientStatus(quit chan bool) {
	for {
		t := time.NewTimer(time.Duration(10) * time.Second)
		select {
		case <-t.C:
			if cache.RedisStatus() {
				redisStatusMetric.Set(float64(1))
			} else {
				redisStatusMetric.Set(float64(-1))
			}
		case <-quit:
			logrus.Debug("quit signal recv redisClientStatus")
			return
		}
	}
}

func Start(port string) {
	quit := make(chan bool)
	defer func() {
		close(quit)
		if err := recover(); err != nil {
			logrus.Error("monitor server occur error ", err)
			os.Exit(1)
		}
	}()
	logrus.Info("monitor server start")
	// start monitor
	server := metrics.NewMonitor(port)
	cronTaskStatusMetric = NewMetricCronWorkStatus()
	redisStatusMetric = NewMetricRedisStatus()
	lcdConnectStatsMetric = NewMetricLcdStatus()
	server.Report(func() {
		go redisClientStatus(quit)
		go lcdConnectionStatus(quit)
	})
}
