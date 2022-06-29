package task

import (
	"encoding/json"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/monitor"
	"strconv"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type TokenPriceTask struct {
}

func (t *TokenPriceTask) Name() string {
	return "ibc_token_price_task"
}

func (t *TokenPriceTask) Cron() string {
	return ThreeMinute
}

func (t *TokenPriceTask) ExpireTime() time.Duration {
	return 3*time.Minute - 1*time.Second
}

func (t *TokenPriceTask) Run() int {
	monitor.SetCronTaskStatusMetricValue(t.Name(), -1)
	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}

	var coinIds []string
	for _, v := range baseDenomList {
		if v.CoinId != "" {
			coinIds = append(coinIds, v.CoinId)
		}
	}

	if len(coinIds) == 0 {
		return -1
	}

	ids := strings.Join(coinIds, ",")
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", global.Config.Spi.CoingeckoPriceUrl, ids)
	bz, err := utils.HttpGet(url)
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}

	var priceResp map[string]map[string]float64
	err = json.Unmarshal(bz, &priceResp)
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}

	priceMap := make(map[string]string)
	for k, v := range priceResp {
		result := strconv.FormatFloat(v["usd"], 'f', 12, 64)
		for strings.HasSuffix(result, "0") {
			result = strings.TrimSuffix(result, "0")
		}
		if strings.HasSuffix(result, ".") {
			result = strings.TrimSuffix(result, ".")
		}
		priceMap[k] = result
	}

	err = tokenPriceRepo.BatchSet(priceMap)
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}
	return 1
}
