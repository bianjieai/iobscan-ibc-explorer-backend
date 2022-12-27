package service

import (
	"fmt"
	"math"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/shopspring/decimal"
)

type IOverviewService interface {
	MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error)
}

var _ IOverviewService = new(OverviewService)

type OverviewService struct {
}

func (svc *OverviewService) MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error) {
	nowTime := time.Now()
	before24hTime := nowTime.AddDate(0, 0, -1)
	statisticsTime, err := denomHeatmapRepo.FindLastStatisticsTime(nowTime)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	denomHeatmapList, err := denomHeatmapRepo.FindByStatisticsTime(statisticsTime)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	resp := svc.buildMarketHeatmapResp(denomHeatmapList)

	before24hStatisticsTime, err := denomHeatmapRepo.FindLastStatisticsTime(before24hTime)
	if err != nil {
		return resp, nil
	}

	before24hDenomHeatmapList, err := denomHeatmapRepo.FindByStatisticsTime(before24hStatisticsTime)
	if err != nil {
		return resp, nil
	}

	if len(before24hDenomHeatmapList) > 0 {
		resp = svc.completeMarketHeatmapResp(resp, before24hDenomHeatmapList)
	}

	return resp, nil
}

func (svc *OverviewService) completeMarketHeatmapResp(resp *vo.MarketHeatmapResp, denomHeatmapList []*entity.DenomHeatmap) *vo.MarketHeatmapResp {
	oldPriceMap := make(map[string]float64, len(denomHeatmapList))
	oldTotalMarketMap := decimal.Zero
	for _, v := range denomHeatmapList {
		oldPriceMap[fmt.Sprintf("%s_%s", v.Chain, v.Denom)] = v.Price
		marketCapDecimal, _ := decimal.NewFromString(v.MarketCap)
		oldTotalMarketMap = oldTotalMarketMap.Add(marketCapDecimal)
	}

	growthRateFunc := func(d1, d2 decimal.Decimal) (rate float64, trend string) {
		rate = d1.DivRound(d2, 4).
			Sub(decimal.NewFromInt(1)).InexactFloat64()
		if rate >= 0 {
			trend = constant.IncreaseSymbol
		} else {
			trend = constant.DecreaseSymbol
		}

		rate = math.Abs(rate)
		return rate, trend
	}

	for i, v := range resp.Items {
		if price, ok := oldPriceMap[fmt.Sprintf("%s_%s", v.Chain, v.Denom)]; ok {
			if price == 0 {
				continue
			}

			resp.Items[i].PriceGrowthRate, resp.Items[i].PriceTrend = growthRateFunc(decimal.NewFromFloat(v.Price), decimal.NewFromFloat(price))
		}
	}

	if !oldTotalMarketMap.Equal(decimal.Zero) {
		totalMarketCap, _ := decimal.NewFromString(resp.TotalInfo.TotalMarketCap)
		resp.TotalInfo.MarketCapGrowthRate, resp.TotalInfo.MarketCapTrend = growthRateFunc(totalMarketCap, oldTotalMarketMap)
	}

	return resp
}

func (svc *OverviewService) buildMarketHeatmapResp(denomHeatmapList []*entity.DenomHeatmap) *vo.MarketHeatmapResp {
	var stableCoinMap entity.IBCBaseDenomMap
	stableCoins, err := authDenomRepo.FindStableCoins()
	if err == nil {
		stableCoinMap = stableCoins.ConvertToMap()
	}

	heatmapItemList := make([]vo.HeatmapItem, 0, len(denomHeatmapList))
	totalMarketCap := decimal.Zero
	stablecoinsMarketCap := decimal.Zero
	transferVolumeTotal := decimal.Zero
	var atomPrice float64
	atomMarketCap := decimal.Zero

	for _, v := range denomHeatmapList {
		marketCapDecimal, _ := decimal.NewFromString(v.MarketCap)
		totalMarketCap = totalMarketCap.Add(marketCapDecimal)

		transferVolume24hDecimal, _ := decimal.NewFromString(v.TransferVolume24h)
		transferVolumeTotal = transferVolumeTotal.Add(transferVolume24hDecimal)

		if _, ok := stableCoinMap[fmt.Sprintf("%s%s", v.Chain, v.Denom)]; ok {
			stablecoinsMarketCap = stablecoinsMarketCap.Add(marketCapDecimal)
		}

		if v.Chain == constant.ChainNameCosmosHub && v.Denom == constant.DenomAtom {
			atomPrice = v.Price
			atomMarketCap, _ = decimal.NewFromString(v.MarketCap)
		}

		heatmapItemList = append(heatmapItemList, vo.HeatmapItem{
			Price:               v.Price,
			PriceGrowthRate:     0,
			PriceTrend:          constant.IncreaseSymbol,
			Denom:               v.Denom,
			Chain:               v.Chain,
			MarketCapValue:      v.MarketCap,
			TransferVolumeValue: v.TransferVolume24h,
		})

	}

	var atomDominance float64
	if !totalMarketCap.Equal(decimal.Zero) {
		atomDominance = atomMarketCap.DivRound(totalMarketCap, 4).InexactFloat64()
	}

	heatmapTotalInfo := vo.HeatmapTotalInfo{
		StablecoinsMarketCap: stablecoinsMarketCap.String(),
		TotalMarketCap:       totalMarketCap.String(),
		MarketCapGrowthRate:  0,
		MarketCapTrend:       constant.IncreaseSymbol,
		TransferVolumeTotal:  transferVolumeTotal.String(),
		AtomPrice:            atomPrice,
		AtomDominance:        atomDominance,
	}

	return &vo.MarketHeatmapResp{
		Items:     heatmapItemList,
		TotalInfo: heatmapTotalInfo,
	}
}
