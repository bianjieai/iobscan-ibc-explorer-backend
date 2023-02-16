package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/shopspring/decimal"
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
	chainFeeStatistics, err := ibcChainFeeStatisticsRepo.ChainFeeStatistics(chain, startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	relayerFeeStatistics, err := ibcChainFeeStatisticsRepo.RelayerFeeStatistics(chain, startTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	denomPriceMap := cache.TokenPriceMap()

	for _, chainFee := range chainFeeStatistics {
		var chainDenom vo.ChainDenomFeeStatistics
		chainDenom.Denom = chainFee.FeeDenom
		chainDenom.TotalAmount = decimal.NewFromFloat(chainFee.FeeAmount).String()
		if coin, ok := denomPriceMap[fmt.Sprintf("%s%s", chainFee.FeeDenom, chain)]; ok {
			if coin.Scale > 0 {
				chainDenom.DenomUSDPrice = decimal.NewFromFloat(coin.Price).String()
				chainDenom.TotalUSDValue = decimal.NewFromFloat(chainFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price)).String()
			}
		}
		denomFeeMap[chainFee.FeeDenom] = chainDenom
	}
	for _, relayerFee := range relayerFeeStatistics {
		if chainDenom, ok := denomFeeMap[relayerFee.FeeDenom]; ok {
			chainDenom.RelayerAmount = decimal.NewFromFloat(relayerFee.FeeAmount).String()
			if coin, exists := denomPriceMap[fmt.Sprintf("%s%s", relayerFee.FeeDenom, chain)]; exists {
				if coin.Scale > 0 {
					chainDenom.RelayerUSDValue = decimal.NewFromFloat(relayerFee.FeeAmount).Div(decimal.NewFromFloat(math.Pow10(coin.Scale))).Mul(decimal.NewFromFloat(coin.Price)).String()
				}
			}
			denomFeeMap[relayerFee.FeeDenom] = chainDenom
		} else {
			return nil, errors.WrapDetail(errors.ErrSystemError, "data error")
		}
	}

	var resp vo.ChainFeeStatisticsResp
	for _, v := range denomFeeMap {
		resp.Fee = append(resp.Fee, v)
	}
	return &resp, nil
}
