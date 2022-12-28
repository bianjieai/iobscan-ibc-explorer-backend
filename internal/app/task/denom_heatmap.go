package task

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/pkg/distributiontask"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type DenomHeatmapTask struct {
	authDenomList  entity.AuthDenomList
	chainConfigMap map[string]*entity.ChainConfig
}

var denomHeatmapTask distributiontask.CronTask = new(DenomHeatmapTask)

func (t *DenomHeatmapTask) Name() string {
	return "denom_heatmap_task"
}
func (t *DenomHeatmapTask) Cron() string {
	if taskConf.CronDenomHeatmapTask != "" {
		return taskConf.CronDenomHeatmapTask
	}

	return "0 * * * * ?"
}

func (t *DenomHeatmapTask) BeforeHook() error {
	return nil
}

func (t *DenomHeatmapTask) Run() {
	if err := t.init(); err != nil {
		return
	}

	nowTime := time.Now()
	statisticsTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), nowTime.Minute(), 0, 0, time.Local)
	wg := sync.WaitGroup{}
	wg.Add(3)
	var txVolumeMap map[string]*dto.Aggr24hDenomVolumeDTO
	var coinPriceMap map[string]float64
	var coinPriceErr error

	go func() {
		defer wg.Done()
		t.supplyHandler()
	}()

	go func() {
		defer wg.Done()
		txVolumeMap = t.volume24h()
	}()

	go func() {
		defer wg.Done()
		coinPriceMap, coinPriceErr = t.coinPriceHandler()
		if coinPriceErr != nil {
			logrus.Errorf("task %s coinPriceHandler err, %v", t.Name(), coinPriceErr)
			coinPriceMap, _ = tokenPriceRepo.GetAll()
		}
	}()

	wg.Wait()

	t.aggrData(statisticsTime, coinPriceMap, txVolumeMap)
}

func (t *DenomHeatmapTask) init() error {
	authDenomList, err := authDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s authDenomRepo.FindAll error, %v", t.Name(), err)
		return err
	}

	t.authDenomList = authDenomList

	chainConfig, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", t.Name(), err)
		return err
	}

	t.chainConfigMap = chainConfig

	return nil
}

// coinPriceHandler Get price of auth dneom, then save price info into cache.
func (t *DenomHeatmapTask) coinPriceHandler() (map[string]float64, error) {
	var coinIds []string
	for _, v := range t.authDenomList {
		if v.CoinId != "" {
			coinIds = append(coinIds, v.CoinId)
		}
	}

	if len(coinIds) == 0 {
		return nil, fmt.Errorf("no coin id")
	}

	ids := strings.Join(coinIds, ",")
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", global.Config.Spi.CoingeckoPriceUrl, ids)
	bz, err := utils.HttpGet(url)
	if err != nil {
		logrus.Errorf("task %s get coin price error, %v", t.Name(), err)
		return nil, err
	}

	var priceResp map[string]map[string]float64
	err = json.Unmarshal(bz, &priceResp)
	if err != nil {
		logrus.Errorf("task %s get coin price error, %v", t.Name(), err)
		return nil, err
	}

	priceMap := make(map[string]string, len(coinIds))
	priceFloatMap := make(map[string]float64, len(coinIds))
	for k, v := range priceResp {
		priceFloatMap[k] = v["usd"]

		result := strconv.FormatFloat(v["usd"], 'f', 12, 64)
		for strings.HasSuffix(result, "0") {
			result = strings.TrimSuffix(result, "0")
		}
		if strings.HasSuffix(result, ".") {
			result = strings.TrimSuffix(result, ".")
		}
		priceMap[k] = result
	}

	if err = tokenPriceRepo.BatchSet(priceMap); err != nil {
		logrus.Errorf("task %s set coin price cache error, %v", t.Name(), err)
	}
	return priceFloatMap, nil
}

// supplyHandler Get supply of denoms, then save supply info to cache
func (t *DenomHeatmapTask) supplyHandler() {
	wg := sync.WaitGroup{}
	wg.Add(len(t.chainConfigMap))
	for _, v := range t.chainConfigMap {
		cf := v
		go func() {
			defer wg.Done()
			t.getSupplyFromLcd(cf)
		}()
	}
	wg.Wait()
}

func (t *DenomHeatmapTask) getSupplyFromLcd(chainCfg *entity.ChainConfig) {
	chain := chainCfg.ChainName
	lcd := chainCfg.GrpcRestGateway
	apiPath := chainCfg.LcdApiPath.SupplyPath
	baseUrl := fmt.Sprintf("%s%s", lcd, apiPath)
	limit := 1000
	key := ""

	for {
		var url string
		if key == "" {
			url = fmt.Sprintf("%s?pagination.limit=%d", baseUrl, limit)
		} else {
			url = fmt.Sprintf("%s?pagination.limit=%d&pagination.key=%s", baseUrl, limit, key)
		}

		bz, err := utils.HttpGet(url)
		if err != nil {
			logrus.Errorf("task %s chain: %s setSupply error, %v", t.Name(), chain, err)
			return
		}

		var supplyResp vo.SupplyResp
		err = json.Unmarshal(bz, &supplyResp)
		if err != nil {
			logrus.Errorf("task %s chain: %s setSupply error, %v", t.Name(), chain, err)
			return
		}

		supplyMap := make(map[string]string, len(supplyResp.Supply))
		for _, v := range supplyResp.Supply {
			supplyMap[v.Denom] = v.Amount
		}

		// 写数据之前，清除之前的老数据
		if key == "" {
			_, _ = denomDataRepo.DelSupply(chain)
		}
		if err = denomDataRepo.BatchSetSupply(chain, supplyMap); err != nil {
			logrus.Errorf("task %s BatchSetSupply err, %v", t.Name(), err)
		}

		if supplyResp.Pagination.NextKey == nil {
			break
		} else {
			key = *supplyResp.Pagination.NextKey
		}
	}
}

func (t *DenomHeatmapTask) aggrData(statisticsTime time.Time, coinPriceMap map[string]float64, txVolumeMap map[string]*dto.Aggr24hDenomVolumeDTO) {
	valueFunc := func(denomAmount decimal.Decimal, price float64, scale int) decimal.Decimal {
		priceDecimal := decimal.NewFromFloat(price)
		value := denomAmount.Div(decimal.NewFromFloat(math.Pow10(scale))).
			Mul(priceDecimal).Round(4)

		return value
	}

	unix := time.Now().Unix()
	denomHeatmapList := make([]*entity.DenomHeatmap, 0, len(t.authDenomList))
	for _, v := range t.authDenomList {
		supplyDecimal := decimal.Zero
		marketCapDecimal := decimal.Zero
		volume24hDecimal := decimal.Zero

		supply, err := denomDataRepo.GetSupply(v.Chain, v.Denom)
		if err != nil {
			logrus.Errorf("task %s GetSupply(%s-%s) err, %v", t.Name(), v.Chain, v.Denom, err)
		} else {
			supplyDecimal, _ = decimal.NewFromString(supply)
		}

		price, ok := coinPriceMap[v.CoinId]
		if !ok {
			denomHeatmapList = append(denomHeatmapList, &entity.DenomHeatmap{
				Denom:             v.Denom,
				Chain:             v.Chain,
				StatisticsTime:    statisticsTime,
				Price:             0,
				Supply:            supplyDecimal.String(),
				MarketCap:         marketCapDecimal.String(),
				TransferVolume24h: volume24hDecimal.String(),
				CreateAt:          unix,
				UpdateAt:          unix,
			})

			continue
		}

		marketCapDecimal = valueFunc(supplyDecimal, price, v.Scale)
		txVolume, ok := txVolumeMap[fmt.Sprintf("%s_%s", v.Denom, v.Chain)]
		if ok {
			denomAmount := decimal.NewFromFloat(txVolume.DenomAmount)
			volume24hDecimal = valueFunc(denomAmount, price, v.Scale)
		}

		denomHeatmapList = append(denomHeatmapList, &entity.DenomHeatmap{
			Denom:             v.Denom,
			Chain:             v.Chain,
			StatisticsTime:    statisticsTime,
			Price:             price,
			Supply:            supplyDecimal.String(),
			MarketCap:         marketCapDecimal.String(),
			TransferVolume24h: volume24hDecimal.String(),
			CreateAt:          unix,
			UpdateAt:          unix,
		})
	}

	if err := denomHeatmap.InsertMany(denomHeatmapList); err != nil {
		logrus.Errorf("task %s denomHeatmap.InsertMany err, %v", t.Name(), err)
	}
}

func (t *DenomHeatmapTask) marketCapHandler(coinPriceMap map[string]string) map[string]decimal.Decimal {
	capMap := make(map[string]decimal.Decimal)

	for _, v := range t.authDenomList {
		key := fmt.Sprintf("%s%s", v.Denom, v.Chain)
		supply, err := denomDataRepo.GetSupply(v.Chain, v.Denom)
		if err != nil {
			logrus.Errorf("task %s GetSupply(%s-%s) err, %v", t.Name(), v.Chain, v.Denom, err)
			capMap[key] = decimal.Zero
			continue
		}

		price, ok := coinPriceMap[key]
		if !ok {
			capMap[key] = decimal.Zero
			continue
		}

		supplyDecimal, _ := decimal.NewFromString(supply)
		priceDecimal, _ := decimal.NewFromString(price)
		value := supplyDecimal.Div(decimal.NewFromFloat(math.Pow10(v.Scale))).
			Mul(priceDecimal).Round(4)
		capMap[key] = value
	}

	return capMap
}

func (t *DenomHeatmapTask) volume24h() map[string]*dto.Aggr24hDenomVolumeDTO {
	volumeMap := make(map[string]*dto.Aggr24hDenomVolumeDTO)
	startTime, _ := last24hTimeUnix()
	txVolumeList, err := ibcTxRepo.Aggr24hDenomVolume(startTime)
	if err != nil {
		logrus.Errorf("task %s Aggr24hDenomVolume err, %v", t.Name(), err)
		return volumeMap
	}

	txVolumeMap := make(map[string]*dto.Aggr24hDenomVolumeDTO, len(txVolumeList))
	for _, v := range txVolumeList {
		txVolumeMap[fmt.Sprintf("%s_%s", v.BaseDenom, v.BaseDenomChain)] = v
	}
	return txVolumeMap
}

// ===================================================================
// ===================================================================
// ===================================================================

type IBCDenomHopsTask struct {
}

func (t *IBCDenomHopsTask) Name() string {
	return "ibc_denom_hops_task"
}

func (t *IBCDenomHopsTask) Switch() bool {
	return true
}

func (t *IBCDenomHopsTask) Run() int {
	denomList, err := denomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s denomRepo.FindAll err, %v", t.Name(), err)
		return -1
	}

	for _, v := range denomList {
		hops := ibcHops(v.DenomPath)
		if hops == 0 {
			continue
		}

		if err = denomRepo.UpdateHops(v.Chain, v.Denom, hops); err != nil {
			logrus.Errorf("task %s UpdateHops %s-%s err, %v", t.Name(), v.Chain, v.Denom, err)
		}
	}

	return 1
}
