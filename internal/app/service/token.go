package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"sort"
)

type ITokenService interface {
	PopularSymbols(minHops int, minReceiveTxs int64) (*vo.PopularSymbolsResp, errors.Error)
}

var _ ITokenService = new(TokenService)

type TokenService struct {
}

func (svc *TokenService) PopularSymbols(minHops int, minReceiveTxs int64) (*vo.PopularSymbolsResp, errors.Error) {
	var resp vo.PopularSymbolsResp
	tokens, err := ibcTokenTraceRepo.FindByHopsAndReceiveTcs(minHops, minReceiveTxs)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var tokenSymbolMap map[string]string
	err = cache.GetRedisClient().UnmarshalHGetAll(cache.DenomSymbol, &tokenSymbolMap)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	maxUpdateAt, err := ibcTokenTraceRepo.FindMaxUpdateAt()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	resp.TimeStamp = maxUpdateAt
	symbols := utils.NewStringSet()
	symbolReceiveTxsAmountMap := make(map[string]int64)
	for _, token := range tokens {
		if symbol, ok := tokenSymbolMap[token.Chain+"/"+token.Denom]; ok {
			if symbol == "" {
				symbols.Add(token.BaseDenom)
				symbolReceiveTxsAmountMap[token.BaseDenom] += token.ReceiveTxs
			} else {
				symbols.Add(symbol)
				symbolReceiveTxsAmountMap[symbol] += token.ReceiveTxs
			}
		}
	}
	symbolsList := symbols.ToSlice()
	var symbolDetails []vo.SymbolDetail
	for _, symbol := range symbolsList {
		var symbolDetail vo.SymbolDetail
		symbolDetail.Symbol = symbol
		symbolDetail.TotalReceiveTxs = symbolReceiveTxsAmountMap[symbol]
		symbolDetails = append(symbolDetails, symbolDetail)
	}

	sort.Slice(symbolDetails, func(i, j int) bool {
		return symbolDetails[i].TotalReceiveTxs > symbolDetails[j].TotalReceiveTxs
	})
	resp.Symbols = symbolDetails
	return &resp, nil
}
