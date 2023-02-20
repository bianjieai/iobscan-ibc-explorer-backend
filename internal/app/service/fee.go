package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"math"
)

type IFeeService interface {
	ChainFeeStatistics(chain string, startTime, endTime int64) (*vo.ChainFeeStatisticsResp, errors.Error)
}

var _ IFeeService = new(FeeService)

type FeeService struct {
}

func (svc *FeeService) ChainFeeStatistics(chain string, startTime, endTime int64) (*vo.ChainFeeStatisticsResp, errors.Error) {
	denomFeeMap := make(map[string]vo.ChainDenomFeeStatistics)
	var chainStatistics []*dto.DenomAmountStatisticsDTO
	var relayerStatistics []*dto.DenomAmountStatisticsDTO
	responseChannel := make(chan map[string][]*dto.DenomAmountStatisticsDTO, 1)
	g := new(errgroup.Group)
	g.Go(func() error {
		chainFeeStatistics, err := ibcChainFeeStatisticsRepo.ChainFeeStatistics(chain, startTime, endTime)
		if err != nil {
			return err
		} else {
			select {
			case tagMap := <-responseChannel:
				tagMap["chain"] = chainFeeStatistics
				responseChannel <- tagMap
			default:
				tagMap := make(map[string][]*dto.DenomAmountStatisticsDTO)
				tagMap["chain"] = chainFeeStatistics
				responseChannel <- tagMap
			}
		}
		return nil
	})

	g.Go(func() error {
		relayerFeeStatistics, err := ibcChainFeeStatisticsRepo.RelayerFeeStatistics(chain, startTime, endTime)
		if err != nil {
			return err
		} else {
			select {
			case tagMap := <-responseChannel:
				tagMap["relayer"] = relayerFeeStatistics
				responseChannel <- tagMap
			default:
				tagMap := make(map[string][]*dto.DenomAmountStatisticsDTO)
				tagMap["relayer"] = relayerFeeStatistics
				responseChannel <- tagMap
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err)
	} else {
		close(responseChannel)
	}

	for response := range responseChannel {
		for k, v := range response {
			if k == "chain" {
				chainStatistics = v
			}
			if k == "relayer" {
				relayerStatistics = v
			}
		}
	}

	denomPriceMap := cache.TokenPriceMap()
	authDenoms, err := authDenomRepo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	denomSymbolMap := authDenoms.ConvertToMap()

	for _, chainFee := range chainStatistics {
		var chainDenom vo.ChainDenomFeeStatistics
		if coin, ok := denomPriceMap[fmt.Sprintf("%s%s", chainFee.FeeDenom, chain)]; ok {
			chainDenom.Denom = coin.Symbol
			chainDenom.TotalAmount = decimal.NewFromFloat(chainFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).String()
			chainDenom.DenomUSDPrice = decimal.NewFromFloat(coin.Price).String()
			chainDenom.TotalUSDValue = decimal.NewFromFloat(chainFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price)).String()
			denomFeeMap[chainFee.FeeDenom] = chainDenom
		} else {
			if denomSymbol, exists := denomSymbolMap[fmt.Sprintf("%s%s", chain, chainFee.FeeDenom)]; exists {
				chainDenom.Denom = denomSymbol.Symbol
				chainDenom.TotalAmount = decimal.NewFromFloat(chainFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(denomSymbol.Scale))).String()
			} else {
				chainDenom.Denom = chainFee.FeeDenom
				chainDenom.TotalAmount = decimal.NewFromFloat(chainFee.FeeAmount).String()
			}
			denomFeeMap[chainFee.FeeDenom] = chainDenom
		}
	}
	for _, relayerFee := range relayerStatistics {
		if chainDenom, ok := denomFeeMap[relayerFee.FeeDenom]; ok {
			if coin, exists := denomPriceMap[fmt.Sprintf("%s%s", relayerFee.FeeDenom, chain)]; exists {
				chainDenom.RelayerAmount = decimal.NewFromFloat(relayerFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).String()
				chainDenom.RelayerUSDValue = decimal.NewFromFloat(relayerFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price)).String()
				denomFeeMap[relayerFee.FeeDenom] = chainDenom
			} else {
				if denomSymbol, has := denomSymbolMap[fmt.Sprintf("%s%s", chain, relayerFee.FeeDenom)]; has {
					chainDenom.RelayerAmount = decimal.NewFromFloat(relayerFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(denomSymbol.Scale))).String()
				} else {
					chainDenom.RelayerAmount = decimal.NewFromFloat(relayerFee.FeeAmount).String()
				}
				denomFeeMap[relayerFee.FeeDenom] = chainDenom
			}
		}
	}

	var resp vo.ChainFeeStatisticsResp
	for _, v := range denomFeeMap {
		resp.Fee = append(resp.Fee, v)
	}
	return &resp, nil
}
