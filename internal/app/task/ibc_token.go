package task

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type TokenTask struct {
	chainIds []string
}

func (t *TokenTask) Name() string {
	return "ibc_token_task"
}

func (t *TokenTask) Cron() string {
	return ThreeMinute
}

func (t *TokenTask) ExpireTime() time.Duration {
	return 3*time.Minute - 1*time.Second
}

func (t *TokenTask) Run() {
	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s run error, %v", t.Name(), err)
		return
	}

	//baseDenomStrList := make([]string, 0, len(baseDenomList))
	//for _, v := range baseDenomList {
	//	baseDenomStrList = append(baseDenomStrList, v.Denom)
	//}
	existedTokenList, newTokenList, err := t.getAllToken(baseDenomList)
	if err != nil {
		return
	}

	// 部分数据统计出错可以直接忽略error,继续计算后面的指标
	_ = t.setTokenPrice(existedTokenList, newTokenList, baseDenomList)

	_ = t.setDenomSupply(existedTokenList, newTokenList, baseDenomList)

	_ = t.setIbcTransferTxs(existedTokenList, newTokenList)
}

func (t *TokenTask) getAllToken(baseDenomList entity.IBCBaseDenomList) (entity.IBCTokenList, entity.IBCTokenList, error) {
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

	baseDenomMap := baseDenomList.ConvertToMap()
	existedTokenMap := existedTokenList.ConvertToMap()
	var newTokenList entity.IBCTokenList

	for _, v := range tokenList {
		_, ok := existedTokenMap[v.ChainId+v.BaseDenom]
		if ok { // token 已存在
			continue
		}

		var tokenType entity.TokenType
		_, ok = baseDenomMap[v.ChainId+v.BaseDenom]
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
			Supply:         "",
			TransferAmount: "",
			TransferTxs:    0,
			ChainsInvolved: 0,
		})
	}

	return existedTokenList, newTokenList, nil
}

func (t *TokenTask) setTokenPrice(existedTokenList, newTokenList entity.IBCTokenList, baseDenomList entity.IBCBaseDenomList) error {
	tokenPriceMap, err := tokenPriceRepo.GetAll()
	if err != nil {
		logrus.Errorf("task %s `setTokenPrice` error, %v", t.Name(), err)
		return err
	}

	baseDenomMap := baseDenomList.ConvertToMap()
	setPrice := func(tokenList entity.IBCTokenList, tokenPriceMap map[string]float64) {
		for _, v := range tokenList {
			denom, ok := baseDenomMap[v.ChainId+v.BaseDenom]
			if !ok {
				continue
			}

			if denom.CoinId == "" {
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

func (t *TokenTask) setDenomSupply(existedTokenList, newTokenList entity.IBCTokenList, baseDenomStrList entity.IBCBaseDenomList) error {
	// 1、先从链上lcd上获取denom的supply，同时获取ibc denom的supply信息。ibc denom的supply在后面会用，此处一并获取了
	configList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s setSupply error, %v", t.Name(), err)
		return err
	}

	chainIds := make([]string, 0, len(configList))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(configList))
	for _, v := range configList {
		chainIds = append(chainIds, v.ChainId)
		var denoms []string
		for _, b := range baseDenomStrList {
			if b.ChainId == v.ChainId {
				denoms = append(denoms, b.Denom)
			}
		}

		go t.getSupplyFromLcd(v.Lcd, v.ChainId, denoms)
	}
	waitGroup.Wait()

	// 2、给base denom设置supply的值
	for _, v := range existedTokenList {
		if utils.InArray(t.chainIds, v.ChainId) {
			v.Supply, _ = denomSupplyRepo.Get(v.ChainId, v.BaseDenom) // 此处忽略error
		}
	}
	for _, v := range newTokenList {
		if utils.InArray(t.chainIds, v.ChainId) {
			v.Supply, _ = denomSupplyRepo.Get(v.ChainId, v.BaseDenom) // 此处忽略error
		}
	}

	return nil
}

func (t *TokenTask) getSupplyFromLcd(lcd, chainId string, denoms []string) {
	page := 1
	limit := 500
	key := ""
	baseUrl := lcd + constant.LcdSupplyUrl

	for {
		var url string
		if page == 1 {
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
			if strings.HasPrefix(v.Denom, constant.IbcTokenPreFix) || utils.InArray(denoms, v.Denom) {
				_ = denomSupplyRepo.Set(chainId, v.Denom, v.Amount)
			}
		}

		if supplyResp.Pagination.NextKey == "" {
			break
		}
		page++
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
				if tx.BaseDenom == v.BaseDenom && (tx.DcChainId == v.ChainId || tx.ScChainId == v.ChainId) {
					count += tx.Count
				}
			}

			for _, tx := range historyTxsCount {
				if tx.BaseDenom == v.BaseDenom && (tx.DcChainId == v.ChainId || tx.ScChainId == v.ChainId) {
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
