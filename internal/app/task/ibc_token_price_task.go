package task

import (
	"encoding/json"
	"fmt"
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

func (t *TokenPriceTask) Run() {
	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return
	}

	var coinIds []string
	for _, v := range baseDenomList {
		if v.CoinId != "" {
			coinIds = append(coinIds, v.CoinId)
		}
	}

	if len(coinIds) == 0 {
		return
	}

	ids := strings.Join(coinIds, ",")
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", global.Config.Spi.CoingeckoPriceUrl, ids)
	bz, err := utils.HttpGet(url)
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return
	}

	var priceResp map[string]map[string]float64
	err = json.Unmarshal(bz, &priceResp)
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return
	}

	priceMap := make(map[string]string)
	for k, v := range priceResp {
		priceMap[k] = strconv.FormatFloat(v["usd"], 'f', 8, 64)
	}

	err = tokenPriceRepo.BatchSet(priceMap)
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return
	}
}
