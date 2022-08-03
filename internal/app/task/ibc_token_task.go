package task

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils/bech32"
	v8 "github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenTask struct {
	chainIds         []string                  // 系统支持的chain列表
	chainLcdMap      map[string]string         // chain lcd地址
	chainLcdApiMap   map[string]entity.ApiPath // chain lcd api地址
	escrowAddressMap map[string][]string       // chain ibc跨链托管地址
	baseDenomList    entity.IBCBaseDenomList   // 所有的base denom
	ibcChainDenomMap map[string][]string       // chain id 和其对应的跨链denom的映射关系
	ibcReceiveTxsMap map[string]int64          // ibc hash token denom的recv txs
}

func (t *TokenTask) Name() string {
	return "ibc_token_task"
}

func (t *TokenTask) Cron() int {
	return ThreeMinute
}

func (t *TokenTask) Run() int {
	err := t.analyzeChainConf()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return -1
	}

	if err = t.initDenomData(); err != nil {
		return -1
	}

	_ = t.todayStatistics()

	_ = t.yesterdayStatistics()

	existedTokenList, newTokenList, err := t.getAllToken()
	if err != nil {
		return -1
	}

	// 部分数据统计出错可以直接忽略error,继续计算后面的指标
	_ = t.setTokenPrice(existedTokenList, newTokenList)

	_ = t.setDenomSupply(existedTokenList, newTokenList)

	_ = t.setIbcTransferTxs(existedTokenList, newTokenList)

	_ = t.setIbcTransferAmount(existedTokenList, newTokenList)

	if err = t.ibcReceiveTxs(); err != nil {
		return -1
	}

	t.calculateTokenStatistics(existedTokenList, newTokenList) // 此步计算ibc_token_statistics的数据，同时设置chains involved字段

	// 更新数据到数据库
	if err = tokenRepo.InsertBatch(newTokenList); err != nil {
		logrus.Errorf("task %s insert new token error, %v", t.Name(), err)
	}

	for _, v := range existedTokenList {
		err = tokenRepo.UpdateToken(v)
		if err != nil && err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s insert update token error, %v", t.Name(), err)
		}
	}

	// 更新ibc chain
	t.updateIBCChain()
	return 1
}

func (t *TokenTask) initDenomData() error {
	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s baseDenomRepo.FindAll error, %v", t.Name(), err)
		return err
	}
	t.baseDenomList = baseDenomList

	denomList, err := denomRepo.GetDenomGroupByChainId()
	if err != nil {
		logrus.Errorf("task %s denomRepo.GetDenomGroupByBaseDenom error, %v", t.Name(), err)
		return err
	}

	denomMap := make(map[string][]string, 0)
	for _, v := range denomList {
		denomMap[v.ChainId] = v.Denom
	}
	t.ibcChainDenomMap = denomMap
	return nil
}

func (t *TokenTask) getAllToken() (entity.IBCTokenList, entity.IBCTokenList, error) {
	tokenList, err := denomRepo.FindBaseDenom()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return nil, nil, err
	}

	existedTokenList, err := tokenRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return nil, nil, err
	}

	existedTokenMap := existedTokenList.ConvertToMap()
	var newTokenList entity.IBCTokenList

	for _, v := range tokenList {
		_, ok := existedTokenMap[v.BaseDenom]
		if ok { // token 已存在
			continue
		}

		// 新增的token
		newTokenList = append(newTokenList, &entity.IBCToken{
			BaseDenom:      v.BaseDenom,
			ChainId:        v.ChainId,
			Type:           t.tokenType(v.BaseDenom, v.ChainId),
			Price:          constant.UnknownTokenPrice,
			Currency:       constant.DefaultCurrency,
			Supply:         constant.UnknownDenomAmount,
			TransferAmount: constant.UnknownDenomAmount,
			TransferTxs:    0,
			ChainsInvolved: 1, // 初始值为 1
		})
	}

	return existedTokenList, newTokenList, nil
}

func (t *TokenTask) tokenType(baseDenom, chainId string) entity.TokenType {
	for _, v := range t.baseDenomList {
		if v.ChainId == chainId && v.Denom == baseDenom {
			return entity.TokenTypeAuthed
		}
	}

	return entity.TokenTypeOther
}

func (t *TokenTask) setTokenPrice(existedTokenList, newTokenList entity.IBCTokenList) error {
	tokenPriceMap, err := tokenPriceRepo.GetAll()
	if err != nil {
		logrus.Errorf("task %s `setTokenPrice` error, %v", t.Name(), err)
		return err
	}

	baseDenomMap := t.baseDenomList.ConvertToMap()
	setPrice := func(tokenList entity.IBCTokenList, tokenPriceMap map[string]float64) {
		for _, v := range tokenList {
			denom, ok := baseDenomMap[v.BaseDenom]
			if !ok || denom.CoinId == "" {
				continue
			}

			price, ok := tokenPriceMap[denom.CoinId]
			if ok {
				v.Price = price
			}
		}
	}

	setPrice(existedTokenList, tokenPriceMap)
	setPrice(newTokenList, tokenPriceMap)
	return nil
}

func (t *TokenTask) analyzeChainConf() error {
	configList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s analyzeChainConf error, %v", t.Name(), err)
		return err
	}

	chainIds := make([]string, 0, len(configList))
	chainLcdMap := make(map[string]string)
	chainLcdApiMap := make(map[string]entity.ApiPath)
	escrowAddressMap := make(map[string][]string)
	for _, v := range configList {
		chainIds = append(chainIds, v.ChainId)
		chainLcdMap[v.ChainId] = v.Lcd
		chainLcdApiMap[v.ChainId] = v.LcdApiPath
		address, err := t.analyzeChainEscrowAddress(v.IbcInfo, v.AddrPrefix)
		if err != nil {
			continue
		}
		escrowAddressMap[v.ChainId] = address
	}
	t.chainIds = chainIds
	t.chainLcdMap = chainLcdMap
	t.chainLcdApiMap = chainLcdApiMap
	t.escrowAddressMap = escrowAddressMap
	return nil
}

func (t *TokenTask) analyzeChainEscrowAddress(info []*entity.IbcInfo, addrPrefix string) ([]string, error) {
	var addrList []string
	for _, v := range info {
		for _, p := range v.Paths {
			address, err := t.getEscrowAddress(p.PortId, p.ChannelId, addrPrefix)
			if err != nil {
				continue
			}
			addrList = append(addrList, address)
		}
	}
	return addrList, nil
}

func (t *TokenTask) getEscrowAddress(portID, channelID, addrPrefix string) (string, error) {
	contents := fmt.Sprintf("%s/%s", portID, channelID)
	const version = "ics20-1"
	preImage := []byte(version)
	preImage = append(preImage, 0)
	preImage = append(preImage, contents...)
	hash := sha256.Sum256(preImage)

	addr, err := bech32.ConvertAndEncode(addrPrefix, hash[:20])
	if err != nil {
		logrus.Errorf("task %s getEscrowAddress error, %v", t.Name(), err)
		return "", err
	}

	return addr, nil
}

func (t *TokenTask) setDenomSupply(existedTokenList, newTokenList entity.IBCTokenList) error {
	// 1、先从链上lcd上获取denom的supply，同时获取ibc denom的supply信息。ibc denom的supply在后面会用，此处一并获取了
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(t.chainIds))
	for _, v := range t.chainIds {
		go func(c string) {
			t.getSupplyFromLcd(c)
			waitGroup.Done()
		}(v)
	}
	waitGroup.Wait()

	// 2、给base denom设置supply的值
	setSupply := func(list entity.IBCTokenList) {
		for _, v := range list {
			if utils.InArray(t.chainIds, v.ChainId) {
				supply, err := denomDataRepo.GetSupply(v.ChainId, v.BaseDenom)
				if err == nil {
					v.Supply = supply
				}
			}
		}
	}

	setSupply(existedTokenList)
	setSupply(newTokenList)
	return nil
}

func (t *TokenTask) getSupplyFromLcd(chainId string) {
	lcd := t.chainLcdMap[chainId]
	apiPath := t.chainLcdApiMap[chainId].SupplyPath
	denoms := t.ibcChainDenomMap[chainId]
	baseUrl := fmt.Sprintf("%s%s", lcd, apiPath)
	limit := 500
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
			logrus.Errorf("task %s setSupply error, %v", t.Name(), err)
			return
		}

		var supplyResp vo.SupplyResp
		err = json.Unmarshal(bz, &supplyResp)
		if err != nil {
			logrus.Errorf("task %s setSupply error, %v", t.Name(), err)
			return
		}

		// 第一页查询成功时，清除之前的老数据
		if key == "" {
			_, _ = denomDataRepo.DelSupply(chainId)
		}

		for _, v := range supplyResp.Supply { // ibc denom 和 链原生denom的amount 存下来
			if strings.HasPrefix(v.Denom, constant.IBCTokenPreFix) || utils.InArray(denoms, v.Denom) {
				_ = denomDataRepo.SetSupply(chainId, v.Denom, v.Amount)
			}
		}

		if supplyResp.Pagination.NextKey == nil {
			break
		} else {
			key = *supplyResp.Pagination.NextKey
		}
	}
}

func (t *TokenTask) setIbcTransferTxs(existedTokenList, newTokenList entity.IBCTokenList) error {
	txsCount, err := tokenStatisticsRepo.Aggr()
	if err != nil {
		logrus.Errorf("task %s setIbcTransferTxs error, %v", t.Name(), err)
		return err
	}

	setTxs := func(tokenList entity.IBCTokenList, txsCount []*dto.CountBaseDenomTxsDTO) {
		for _, v := range tokenList {
			var count int64
			for _, tx := range txsCount {
				if tx.BaseDenom == v.BaseDenom {
					count += tx.Count
				}
			}

			v.TransferTxs = count
		}
	}

	setTxs(existedTokenList, txsCount)
	setTxs(newTokenList, txsCount)
	return nil
}

func (t *TokenTask) setIbcTransferAmount(existedTokenList, newTokenList entity.IBCTokenList) error {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(t.chainIds))
	for _, v := range t.chainIds {
		go func(c string, addrs []string) {
			t.getTransAmountFromLcd(c, addrs)
			waitGroup.Done()
		}(v, t.escrowAddressMap[v])
	}
	waitGroup.Wait()

	setTransAmount := func(list entity.IBCTokenList) {
		for _, v := range list {
			if utils.InArray(t.chainIds, v.ChainId) {
				amount, err := denomDataRepo.GetTransferAmount(v.ChainId, v.BaseDenom)
				if err != nil {
					if err == v8.Nil {
						v.TransferAmount = constant.ZeroDenomAmount
					}
					continue
				}
				v.TransferAmount = amount
			}
		}
	}

	setTransAmount(existedTokenList)
	setTransAmount(newTokenList)
	return nil
}

func (t *TokenTask) getTransAmountFromLcd(chainId string, addrList []string) {
	denomTransAmountMap := make(map[string]decimal.Decimal)
	lcd := t.chainLcdMap[chainId]
	apiPath := t.chainLcdApiMap[chainId].BalancesPath
	for _, addr := range addrList { // 一条链上的所有地址都要查询一遍，并按denom分组计数
		limit := 500
		key := ""
		earlyTermination := false
		baseUrl := strings.ReplaceAll(fmt.Sprintf("%s%s", lcd, apiPath), entity.ApiBalancesPathPlaceholder, addr)

		for { // 计算地址上所锁定的denom的数量
			var url string
			if key == "" {
				url = fmt.Sprintf("%s?pagination.limit=%d", baseUrl, limit)
			} else {
				url = fmt.Sprintf("%s?pagination.limit=%d&pagination.key=%s", baseUrl, limit, key)
			}

			bz, err := utils.HttpGet(url)
			if err != nil {
				if isConnectionErr(err) {
					earlyTermination = true
				}
				logrus.Errorf("task %s getTransAmountFromLcd error, %v", t.Name(), err)
				break
			}

			var balancesResp vo.BalancesResp
			err = json.Unmarshal(bz, &balancesResp)
			if err != nil {
				logrus.Errorf("task %s getTransAmountFromLcd error, %v", t.Name(), err)
				break
			}

			for _, v := range balancesResp.Balances {
				amount, err := decimal.NewFromString(v.Amount)
				if err != nil {
					logrus.Errorf("task %s getTransAmountFromLcd error, %v", t.Name(), err)
					continue
				}

				d, ok := denomTransAmountMap[v.Denom]
				if !ok {
					denomTransAmountMap[v.Denom] = amount
				} else {
					denomTransAmountMap[v.Denom] = d.Add(amount)
				}
			}

			if balancesResp.Pagination.NextKey == nil {
				break
			} else {
				key = *balancesResp.Pagination.NextKey
			}
		}

		if earlyTermination {
			break
		}
	}

	if len(denomTransAmountMap) > 0 {
		_, _ = denomDataRepo.DelTransferAmount(chainId) // 清除旧数据
		denomTransAmountStrMap := make(map[string]string)
		for k, v := range denomTransAmountMap {
			denomTransAmountStrMap[k] = v.String()
		}
		if err := denomDataRepo.SetTransferAmount(chainId, denomTransAmountStrMap); err != nil {
			logrus.Errorf("task %s denomDataRepo.SetTransferAmount error, %v", t.Name(), err)
		}
	}
}

func (t *TokenTask) getTokenScale(baseDenom, chainId string) int {
	var scale int
	for _, v := range t.baseDenomList {
		if v.ChainId == chainId && v.Denom == baseDenom {
			scale = v.Scale
			break
		}
	}

	return scale
}

func (t *TokenTask) calculateTokenStatistics(existedTokenList, newTokenList entity.IBCTokenList) {
	for _, v := range existedTokenList {
		chainNum, err := t.ibcTokenStatistics(v)
		if err != nil {
			continue
		}
		v.ChainsInvolved = chainNum
	}

	for _, v := range newTokenList {
		chainNum, err := t.ibcTokenStatistics(v)
		if err != nil {
			continue
		}
		v.ChainsInvolved = chainNum
	}
}

// ==============================================================================================================
// ==============================================================================================================
// ==============================================================================================================
// 以下主要是对于ibc_token_statistics 集合数据的处理与计算

func (t *TokenTask) ibcTokenStatistics(ibcToken *entity.IBCToken) (int64, error) {
	ibcDenomCalculateList, err := denomCalculateRepo.FindByBaseDenom(ibcToken.BaseDenom)
	if err != nil {
		logrus.Errorf("task %s denomCaculateRepo.FindByBaseDenom error, %v", t.Name(), err)
		return 0, nil
	}

	denomList, err := denomRepo.FindByBaseDenom(ibcToken.BaseDenom)
	if err != nil {
		logrus.Errorf("task %s FindByBaseDenom error, %v", t.Name(), err)
		return 0, nil
	}

	ibcDenomCalculateStrList := make([]string, 0, len(ibcDenomCalculateList))
	for _, v := range ibcDenomCalculateList {
		ibcDenomCalculateStrList = append(ibcDenomCalculateStrList, fmt.Sprintf("%s%s", v.ChainId, v.Denom))
	}

	scale := t.getTokenScale(ibcToken.BaseDenom, ibcToken.ChainId)
	allTokenStatisticsList := make([]*entity.IBCTokenTrace, 0, len(denomList))
	chainsSet := utils.NewStringSet()
	for _, v := range denomList {
		denomType := t.ibcTokenStatisticsType(ibcToken.BaseDenom, v.Denom, v.ChainId, ibcDenomCalculateStrList)
		var denomAmount string
		if denomType == entity.TokenTraceTypeGenesis {
			denomAmount = t.ibcDenomAmountGenesis(ibcToken.Supply, ibcToken.TransferAmount)
		} else {
			denomAmount = t.ibcDenomAmount(v.ChainId, v.Denom)
		}

		if denomAmount == constant.ZeroDenomAmount { // 为0说明此链已经没有这个ibc denom
			continue
		}

		allTokenStatisticsList = append(allTokenStatisticsList, &entity.IBCTokenTrace{
			Denom:       v.Denom,
			DenomPath:   v.DenomPath,
			BaseDenom:   ibcToken.BaseDenom,
			ChainId:     v.ChainId,
			OriginalId:  ibcToken.ChainId,
			Type:        denomType,
			IBCHops:     t.ibcHops(v.DenomPath),
			DenomAmount: denomAmount,
			DenomValue:  t.ibcDenomValue(denomAmount, ibcToken.Price, scale).Round(constant.DefaultValuePrecision).String(),
			ReceiveTxs:  t.ibcReceiveTxsMap[fmt.Sprintf("%s%s", v.Denom, v.ChainId)],
		})

		chainsSet.Add(v.ChainId)
	}

	err = tokenTraceRepo.BatchSwap(allTokenStatisticsList, ibcToken.BaseDenom, ibcToken.ChainId) // 删除旧数据，插入新数据
	if err != nil {
		logrus.Errorf("task %s BatchSwap error,base denom:%s, %v", t.Name(), ibcToken.BaseDenom, err)
		return 0, err
	}
	return int64(len(chainsSet)), nil
}

func (t *TokenTask) ibcReceiveTxs() error {
	var txsMap = make(map[string]int64)

	aggr, err := tokenTraceStatisticsRepo.Aggr()
	if err != nil {
		logrus.Errorf("task %s ibcReceiveTxs error, %v", t.Name(), err)
		return err
	}

	for _, v := range aggr {
		txsMap[fmt.Sprintf("%s%s", v.Denom, v.ChainId)] = v.ReceiveTxs
	}

	t.ibcReceiveTxsMap = txsMap
	return nil
}

func (t *TokenTask) ibcDenomValue(amount string, price float64, scale int) decimal.Decimal {
	if amount == constant.UnknownDenomAmount || amount == constant.ZeroDenomAmount || price == 0 || price == constant.UnknownTokenPrice {
		return decimal.Zero
	}

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		logrus.Errorf("task %s ibcDenomValue error, %v", t.Name(), err)
		return decimal.Zero
	}

	value := amountDecimal.Div(decimal.NewFromFloat(math.Pow10(scale))).
		Mul(decimal.NewFromFloat(price))

	return value
}

func (t *TokenTask) ibcDenomAmount(chainId, denom string) string {
	amount, err := denomDataRepo.GetSupply(chainId, denom)
	if err != nil {
		if err == v8.Nil {
			return constant.ZeroDenomAmount
		}
		return constant.UnknownDenomAmount
	}
	return amount
}

func (t *TokenTask) ibcDenomAmountGenesis(supply, transAmount string) string {
	if supply == constant.UnknownDenomAmount || transAmount == constant.UnknownDenomAmount {
		return constant.UnknownDenomAmount
	}

	sd, err := decimal.NewFromString(supply)
	if err != nil {
		return constant.UnknownDenomAmount
	}

	td, err := decimal.NewFromString(transAmount)
	if err != nil {
		return constant.UnknownDenomAmount
	}

	return sd.Sub(td).String()
}

func (t *TokenTask) ibcTokenStatisticsType(baseDenom, denom, chainId string, ibcHash []string) entity.TokenTraceType {
	if baseDenom == denom {
		return entity.TokenTraceTypeGenesis
	}

	if utils.InArray(ibcHash, fmt.Sprintf("%s%s", chainId, denom)) {
		return entity.TokenTraceTypeAuthed
	} else {
		return entity.TokenTraceTypeOther
	}
}

func (t *TokenTask) ibcHops(denomPath string) int {
	return strings.Count(denomPath, constant.IBCHopsIndex)
}

// ==============================================================================================================
// ==============================================================================================================
// ==============================================================================================================
// 以下主要是对于其他被依赖集合数据的更新

func (t *TokenTask) updateIBCChain() {
	ibcChainList, err := tokenTraceRepo.AggregateIBCChain()
	if err != nil {
		logrus.Errorf("task %s updateIBCChain error, %v", t.Name(), err)
		return
	}

	for _, v := range ibcChainList {
		vd := decimal.NewFromFloat(v.DenomValue).Round(constant.DefaultValuePrecision).String()
		if err = chainRepo.UpdateIbcTokenValue(v.ChainId, v.Count, vd); err != nil && err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s updateIBCChain error, %v", t.Name(), err)
		}
	}
}

// ==============================================================================================================
// ==============================================================================================================
// ==============================================================================================================
// 以下主要是增量统计

func (t *TokenTask) todayStatistics() error {
	logrus.Infof("task %s exec today statistics", t.Name())
	startTime, endTime := todayUnix()
	segments := []*segment{
		{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if err := tokenStatisticsTask.deal(segments, opUpdate); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}

	return nil
}

func (t *TokenTask) yesterdayStatistics() error {
	mmdd := time.Now().Format(constant.TimeFormatMMDD)
	incr, _ := statisticsCheckRepo.GetIncr(t.Name(), mmdd)
	if incr > statisticsCheckTimes {
		return nil
	}

	logrus.Infof("task %s check yeaterday statistics, times: %d", t.Name(), incr)
	startTime, endTime := yesterdayUnix()
	segments := []*segment{
		{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}
	if err := tokenStatisticsTask.deal(segments, opUpdate); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}

	_ = statisticsCheckRepo.Incr(t.Name(), mmdd)
	return nil
}
