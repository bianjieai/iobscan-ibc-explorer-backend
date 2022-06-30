package task

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"

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

	existedTokenList, newTokenList, err := t.getAllToken()
	if err != nil {
		return -1
	}

	// 部分数据统计出错可以直接忽略error,继续计算后面的指标
	_ = t.setTokenPrice(existedTokenList, newTokenList)

	_ = t.setDenomSupply(existedTokenList, newTokenList)

	_ = t.setIbcTransferTxs(existedTokenList, newTokenList)

	_ = t.setIbcTransferAmount(existedTokenList, newTokenList)

	t.caculateTokenStatistics(existedTokenList, newTokenList) // 此步计算ibc_token_statistics的数据，同时设置chains involved字段

	// 更新数据到数据库
	err = tokenRepo.InsertBatch(newTokenList)
	if err != nil {
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

	baseDenomMap := t.baseDenomList.ConvertToMap()
	existedTokenMap := existedTokenList.ConvertToMap()
	var newTokenList entity.IBCTokenList

	for _, v := range tokenList {
		_, ok := existedTokenMap[v.BaseDenom]
		if ok { // token 已存在
			continue
		}

		var tokenType entity.TokenType
		_, ok = baseDenomMap[v.BaseDenom]
		if ok {
			tokenType = entity.TokenTypeAuthed
		} else {
			tokenType = entity.TokenTypeOther
		}

		// 新增的token
		newTokenList = append(newTokenList, &entity.IBCToken{
			BaseDenom:      v.BaseDenom,
			ChainId:        v.ChainId,
			Type:           tokenType,
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
	txsCount, err := ibcTxRepo.CountBaseDenomTransferTxs()
	if err != nil {
		logrus.Errorf("task %s setIbcTransferTxs error, %v", t.Name(), err)
		return err
	}

	historyTxsCount, err := ibcTxRepo.CountBaseDenomHistoryTransferTxs()
	if err != nil {
		logrus.Errorf("task %s setIbcTransferTxs error, %v", t.Name(), err)
		return err
	}

	setTxs := func(tokenList entity.IBCTokenList, txsCount, historyTxsCount []*dto.CountBaseDenomTransferAmountDTO) {
		for _, v := range tokenList {
			var count int64
			for _, tx := range txsCount {
				if tx.BaseDenom == v.BaseDenom {
					count += tx.Count
				}
			}

			for _, tx := range historyTxsCount {
				if tx.BaseDenom == v.BaseDenom {
					count += tx.Count
				}
			}
			v.TransferTxs = count
		}
	}

	setTxs(existedTokenList, txsCount, historyTxsCount)
	setTxs(newTokenList, txsCount, historyTxsCount)
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
				if t.isConnectionErr(err) {
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
		denomTransAmountStrMap := make(map[string]string)
		for k, v := range denomTransAmountMap {
			denomTransAmountStrMap[k] = v.String()
		}
		err := denomDataRepo.SetTransferAmount(chainId, denomTransAmountStrMap)
		if err != nil {
			logrus.Errorf("task %s getTransAmountFromLcd error, %v", t.Name(), err)
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

func (t *TokenTask) caculateTokenStatistics(existedTokenList, newTokenList entity.IBCTokenList) {
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
	ibcDenomCaculateList, err := denomCaculateRepo.FindByBaseDenom(ibcToken.BaseDenom)
	if err != nil {
		logrus.Errorf("task %s denomCaculateRepo.FindByBaseDenom error, %v", t.Name(), err)
		return 0, nil
	}

	denomList, err := denomRepo.FindByBaseDenom(ibcToken.BaseDenom)
	if err != nil {
		logrus.Errorf("task %s FindByBaseDenom error, %v", t.Name(), err)
		return 0, nil
	}

	ibcDenomCaculateStrList := make([]string, 0, len(ibcDenomCaculateList))
	for _, v := range ibcDenomCaculateList {
		ibcDenomCaculateStrList = append(ibcDenomCaculateStrList, v.Denom)
	}
	ibcReceiveTxsMap := t.ibcReceiveTxs(ibcToken.BaseDenom)
	scale := t.getTokenScale(ibcToken.BaseDenom, ibcToken.ChainId)

	allTokenStatisticsList := make([]*entity.IBCTokenStatistics, 0, len(denomList))
	for _, v := range denomList {
		denomType := t.ibcTokenStatisticsType(ibcToken.BaseDenom, v.Denom, ibcDenomCaculateStrList)
		var denomAmount string
		if denomType == entity.TokenStatisticsTypeGenesis {
			denomAmount = t.ibcDenomAmountGenesis(ibcToken.Supply, ibcToken.TransferAmount)
		} else {
			denomAmount = t.ibcDenomAmount(v.ChainId, v.Denom)
		}

		allTokenStatisticsList = append(allTokenStatisticsList, &entity.IBCTokenStatistics{
			Denom:       v.Denom,
			DenomPath:   v.DenomPath,
			BaseDenom:   ibcToken.BaseDenom,
			ChainId:     v.ChainId,
			OriginalId:  ibcToken.ChainId,
			Type:        denomType,
			IBCHops:     t.ibcHops(v.DenomPath),
			DenomAmount: denomAmount,
			DenomValue:  t.ibcDenomValue(denomAmount, ibcToken.Price, scale).Round(constant.DefaultValuePrecision).String(),
			ReceiveTxs:  ibcReceiveTxsMap[v.Denom],
		})
	}

	err = tokenStatisticsRepo.BatchSwap(allTokenStatisticsList, ibcToken.BaseDenom, ibcToken.ChainId) // 删除旧数据，插入新数据
	if err != nil {
		logrus.Errorf("task %s BatchSwap error,base denom:%s, %v", t.Name(), ibcToken.BaseDenom, err)
		return 0, err
	}
	return int64(len(allTokenStatisticsList)), nil
}

func (t *TokenTask) ibcReceiveTxs(baseDenom string) map[string]int64 {
	var txsMap = make(map[string]int64)

	txs, err := ibcTxRepo.CountIBCTokenRecvTxs(baseDenom)
	if err != nil {
		logrus.Errorf("task %s ibcReceiveTxs error, %v", t.Name(), err)
	} else {
		for _, v := range txs {
			txsMap[v.Denom] += v.Count
		}
	}

	historyTxs, err := ibcTxRepo.CountIBCTokenHistoryRecvTxs(baseDenom)
	if err != nil {
		logrus.Errorf("task %s ibcReceiveTxs error, %v", t.Name(), err)
	} else {
		for _, v := range historyTxs {
			txsMap[v.Denom] += v.Count
		}
	}

	return txsMap
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

func (t *TokenTask) ibcTokenStatisticsType(baseDenom, denom string, ibcHash []string) entity.TokenStatisticsType {
	if baseDenom == denom {
		return entity.TokenStatisticsTypeGenesis
	}

	if utils.InArray(ibcHash, denom) {
		return entity.TokenStatisticsTypeAuthed
	} else {
		return entity.TokenStatisticsTypeOther
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
	ibcChainList, err := tokenStatisticsRepo.AggregateIBCChain()
	if err != nil {
		logrus.Errorf("task %s updateIBCChain error, %v", t.Name(), err)
		return
	}

	for _, v := range ibcChainList {
		vd := decimal.NewFromFloat(v.DenomValue).String()
		if err = chainRepo.UpdateIbcTokenValue(v.ChainId, v.Count, vd); err != nil && err != mongo.ErrNoDocuments {
			logrus.Errorf("task %s updateIBCChain error, %v", t.Name(), err)
		}
	}
}

func (t *TokenTask) isConnectionErr(err error) bool {
	return strings.Contains(err.Error(), "connection refused")
}
