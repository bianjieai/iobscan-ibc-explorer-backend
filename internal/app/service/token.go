package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"time"
)

type ITokenService interface {
	PopularSymbols(minHops int, minReceiveTxs int64) (*vo.PopularSymbolsResp, errors.Error)
}

var _ ITokenService = new(TokenService)

type TokenService struct {
}

func (svc *TokenService) PopularSymbols(minHops int, minReceiveTxs int64) (*vo.PopularSymbolsResp, errors.Error) {
	var resp vo.PopularSymbolsResp
	hopsTokens, err := ibcDenomRepo.FindSymbolDenomsByHops(minHops)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var tokenReceiveTxsMap map[string]int64
	err = cache.GetRedisClient().UnmarshalHGetAll(cache.TokenReceiveTxs, &tokenReceiveTxsMap)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	ttl, err := cache.GetRedisClient().TTL(cache.TokenReceiveTxs)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	timestamp := time.Now().Unix() - (360 - int64(ttl.Seconds()))
	resp.TimeStamp = timestamp
	symbols := utils.NewStringSet()
	for _, token := range hopsTokens {
		if count, ok := tokenReceiveTxsMap[token.Chain+"/"+token.Denom]; ok {
			if count >= minReceiveTxs {
				symbols.Add(token.Symbol)
			}
		} else {
			if minReceiveTxs == 0 {
				symbols.Add(token.Symbol)
			}
		}
	}
	resp.Symbols = symbols.ToSlice()
	return &resp, nil
}
