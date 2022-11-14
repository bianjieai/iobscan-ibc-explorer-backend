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
	v1beta1        = "v1beta1"
	v1             = "v1"
	nodeInfo       = "/node_info"
	v1Channels     = "/ibc/core/channel/v1/channels?pagination.limit=1"
	apiChannels    = "/ibc/core/channel/%s/channels?pagination.offset=OFFSET&pagination.limit=LIMIT&pagination.count_total=true"
	apiClientState = "/ibc/core/channel/%s/channels/CHANNEL/ports/PORT/client_state"
)

// unbelievableLcd 不可信的lcd
var unbelievableLcd = map[string][]string{
	"sifchain_1": {"https://api.sifchain.chaintools.tech/"},
}

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
		t := time.NewTimer(time.Duration(120) * time.Second)
		select {
		case <-t.C:
			chainCfgs, err := chainConfigRepo.FindAllOpenChainInfos()
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			for _, val := range chainCfgs {
				if checkAndUpdateLcd(val.Lcd, val) {
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

// checkAndUpdateLcd If lcd is ok, update db and return true. Else return false
func checkAndUpdateLcd(lcd string, cf *entity.ChainConfig) bool {
	unLcds, ex := unbelievableLcd[cf.ChainId]
	if ex && utils.InArray(unLcds, lcd) {
		return false
	}
	if resp, err := utils.HttpGet(fmt.Sprintf("%s%s", lcd, nodeInfo)); err == nil {
		var data struct {
			NodeInfo struct {
				Network string `json:"network"`
			} `json:"node_info"`
		}
		if err := json.Unmarshal(resp, &data); err != nil {
			return false
		}
		network := strings.ReplaceAll(data.NodeInfo.Network, "-", "_")
		if network != cf.ChainId {
			//return false, if lcd node_info network no match chain_id
			return false
		}

	} else {
		// return false,if lcd node_ifo api is not reach
		return false
	}

	var ok bool
	var version string
	if _, err := utils.HttpGet(fmt.Sprintf("%s%s", lcd, v1Channels)); err == nil {
		ok = true
		version = v1
	} else if strings.Contains(err.Error(), "501 Not Implemented") {
		ok = true
		version = v1beta1
	} else {
		ok = false
	}

	if ok {
		if cf.Lcd == lcd && cf.LcdApiPath.ChannelsPath == fmt.Sprintf(apiChannels, version) && cf.LcdApiPath.ClientStatePath == fmt.Sprintf(apiClientState, version) {
			return true
		}

		cf.Lcd = lcd
		cf.LcdApiPath.ChannelsPath = fmt.Sprintf(apiChannels, version)
		cf.LcdApiPath.ClientStatePath = fmt.Sprintf(apiClientState, version)
		if err := chainConfigRepo.UpdateLcdApi(cf); err != nil {
			logrus.Errorf("lcd monitor update api error: %v", err)
			return false
		} else {
			return true
		}
	}

	return false
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
		if ok := checkAndUpdateLcd(v.Address, chainConf); ok {
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
