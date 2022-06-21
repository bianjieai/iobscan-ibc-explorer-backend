package task

import (
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/sirupsen/logrus"
)

type TokenTask struct {
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

	existedTokenList, newTokenList, err := t.getAllToken(baseDenomList)
	if err != nil {
		return
	}

	existedTokenList, newTokenList, _ = t.setTokenPrice(existedTokenList, newTokenList, baseDenomList) // 忽略此步的error, 此步出错继续计算后面的指标

	//logrus.Info(utils.MustMarshalJsonToStr(newTokenList))
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

func (t *TokenTask) setTokenPrice(existedTokenList, newTokenList entity.IBCTokenList, baseDenomList entity.IBCBaseDenomList) (entity.IBCTokenList, entity.IBCTokenList, error) {
	tokenPriceMap, err := tokenPriceRepo.GetAll()
	if err != nil {
		logrus.Errorf("task %s `setTokenPrice` error, %v", t.Name(), err)
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
	return existedTokenList, newTokenList, nil
}
