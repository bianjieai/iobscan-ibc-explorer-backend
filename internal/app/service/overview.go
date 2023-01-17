package service

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/shopspring/decimal"
)

type IOverviewService interface {
	MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error)
	TokenDistribution(req *vo.TokenDistributionReq) (*vo.TokenDistributionResp, errors.Error)
	ChainVolumeTrend(req *vo.ChainVolumeTrendReq) (*vo.ChainVolumeTrendResp, errors.Error)
	ChainVolume(req *vo.ChainVolumeReq) (*vo.ChainVolumeResp, errors.Error)
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
	TotalTransferVolume := decimal.Zero
	var atomPrice float64
	atomMarketCap := decimal.Zero

	for _, v := range denomHeatmapList {
		marketCapDecimal, _ := decimal.NewFromString(v.MarketCap)
		totalMarketCap = totalMarketCap.Add(marketCapDecimal)

		transferVolume24hDecimal, _ := decimal.NewFromString(v.TransferVolume24h)
		TotalTransferVolume = TotalTransferVolume.Add(transferVolume24hDecimal)

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
		TotalDenomNumber:     len(heatmapItemList),
		MarketCapGrowthRate:  0,
		MarketCapTrend:       constant.IncreaseSymbol,
		TotalTransferVolume:  TotalTransferVolume.String(),
		AtomPrice:            atomPrice,
		AtomDominance:        atomDominance,
	}

	return &vo.MarketHeatmapResp{
		Items:     heatmapItemList,
		TotalInfo: heatmapTotalInfo,
	}
}

func (svc *OverviewService) TokenDistribution(req *vo.TokenDistributionReq) (*vo.TokenDistributionResp, errors.Error) {
	ibcDenoms, err := denomRepo.FindByBaseDenom(req.BaseDenom, req.BaseDenomChain)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(ibcDenoms) == 0 {
		return nil, errors.WrapNoDataErr()
	}

	ibcDenomMap := ibcDenoms.ConvertToMap()
	ibcTokens, err := tokenStatisticsRepo.FindByBaseDenom(req.BaseDenom, req.BaseDenomChain)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(ibcTokens) == 0 {
		return nil, errors.WrapNoDataErr()
	}

	var genesisDenomElem *vo.GraphElem
	mapHopsData := make(map[int][]*vo.GraphElem, len(ibcTokens))
	usedDenomFlagMap := make(map[string]bool, len(ibcTokens))
	for _, val := range ibcTokens {
		if val.Denom == req.BaseDenom && val.Chain == req.BaseDenomChain {
			genesisDenomElem = &vo.GraphElem{
				Supply: val.DenomSupply,
				Amount: val.DenomAmount,
				Denom:  val.Denom,
				Chain:  val.Chain,
				Hops:   val.IBCHops,
			}
		}

		hop := val.IBCHops
		if hop == 0 {
			continue
		}

		elem := &vo.GraphElem{
			Supply: val.DenomSupply,
			Amount: val.DenomAmount,
			Denom:  val.Denom,
			Chain:  val.Chain,
			Hops:   val.IBCHops,
		}
		uk := fmt.Sprintf("%s%s", val.Chain, val.Denom)
		usedDenomFlagMap[uk] = false
		if td, ok := ibcDenomMap[uk]; ok {
			elem.PrevChain = td.PrevChain
			elem.PrevDenom = td.PrevDenom
		}

		hopDatas, exist := mapHopsData[hop]
		if exist {
			mapHopsData[hop] = append(hopDatas, elem)
		} else {
			mapHopsData[hop] = []*vo.GraphElem{elem}
		}
	}

	if genesisDenomElem == nil || genesisDenomElem.Supply == constant.UnknownDenomAmount || genesisDenomElem.Supply == constant.ZeroDenomAmount {
		return nil, errors.WrapNoDataErr()
	}

	graphData := &vo.GraphData{
		Children:  []*vo.GraphData{},
		GraphElem: genesisDenomElem,
	}
	hop := 1
	oneHopDenoms, ok := mapHopsData[hop]
	if !ok {
		graphData = svc.fillJumpHopsElem(ibcTokens, usedDenomFlagMap, graphData)
		return &vo.TokenDistributionResp{GraphData: graphData}, nil
	}

	for _, val := range oneHopDenoms {
		children := vo.GraphData{
			Children:  nil,
			GraphElem: val,
		}
		children = svc.findChildrens(mapHopsData, usedDenomFlagMap, children)
		graphData.Children = append(graphData.Children, &children)
		usedDenomFlagMap[fmt.Sprintf("%s%s", val.Chain, val.Denom)] = true
	}

	graphData = svc.fillJumpHopsElem(ibcTokens, usedDenomFlagMap, graphData)
	resp := vo.TokenDistributionResp{GraphData: graphData}
	return &resp, nil
}

func (svc *OverviewService) findChildrens(mapHopsData map[int][]*vo.GraphElem, usedDenomFlagMap map[string]bool, ret vo.GraphData) vo.GraphData {
	hopDenoms, ok := mapHopsData[ret.Hops+1]
	if !ok {
		ret.Children = []*vo.GraphData{}
		return ret
	}

	if ret.Supply == constant.ZeroDenomAmount || ret.Supply == constant.UnknownDenomAmount {
		ret.Children = []*vo.GraphData{}
		return ret
	}

	ret.Children = make([]*vo.GraphData, 0)
	for _, val := range hopDenoms {
		if val.PrevDenom == ret.Denom && val.PrevChain == ret.Chain {
			children := vo.GraphData{
				Children:  nil,
				GraphElem: val,
			}
			children = svc.findChildrens(mapHopsData, usedDenomFlagMap, children)
			ret.Children = append(ret.Children, &children)
			usedDenomFlagMap[fmt.Sprintf("%s%s", val.Chain, val.Denom)] = true
		}
	}

	return ret
}

func (svc *OverviewService) fillJumpHopsElem(ibcTokens []*entity.IBCTokenTrace, usedDenomFlagMap map[string]bool, ret *vo.GraphData) *vo.GraphData {
	tokenMap := make(map[string]*entity.IBCTokenTrace)
	for _, v := range ibcTokens {
		tokenMap[fmt.Sprintf("%s%s", v.Chain, v.Denom)] = v
	}

	for k, flag := range usedDenomFlagMap {
		if flag {
			continue
		}
		token, ok := tokenMap[k]
		if !ok || token.IBCHops <= 1 || token.DenomSupply == constant.ZeroDenomAmount || token.DenomSupply == constant.UnknownDenomAmount {
			continue
		}

		var gdp = &vo.GraphData{
			Children: nil,
			GraphElem: &vo.GraphElem{
				Supply: token.DenomSupply,
				Amount: "0",
				Hops:   1,
			},
		}
		rootGdp := gdp
		for i := 2; i <= token.IBCHops; i++ {
			elem := &vo.GraphElem{
				Supply: token.DenomSupply,
				Amount: "0",
				Hops:   i,
			}
			if i == token.IBCHops {
				elem.Chain = token.Chain
				elem.Denom = token.Denom
				elem.Amount = token.DenomAmount
			}

			gdp.Children = append(gdp.Children, &vo.GraphData{
				Children:  []*vo.GraphData{},
				GraphElem: elem,
			})
			gdp = gdp.Children[0]
		}

		ret.Children = append(ret.Children, rootGdp)
	}

	return ret
}

func (svc *OverviewService) ChainVolumeTrend(req *vo.ChainVolumeTrendReq) (*vo.ChainVolumeTrendResp, errors.Error) {
	fillVolumeItems := func(items []vo.VolumeItem) []vo.VolumeItem { // 若items 不足365个，则补足至365个
		volumeMap := make(map[string]string, len(items))
		for _, v := range items {
			volumeMap[v.Datetime] = v.Value
		}
		date := time.Now().AddDate(0, 0, -constant.ChainFlowTrendDays+1)
		startUnix := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).Unix()
		newItems := make([]vo.VolumeItem, 0, constant.ChainFlowTrendDays)
		for i := 0; i < constant.ChainFlowTrendDays; i++ {
			dt := time.Unix(startUnix+int64(i*86400), 0).Format(constant.DateFormat)
			value := "0"
			if _, ok := volumeMap[dt]; ok {
				value = volumeMap[dt]
			}
			newItems = append(newItems, vo.VolumeItem{
				Datetime: dt,
				Value:    value,
			})
		}
		return newItems
	}

	if req.Chain != "" {
		//check chain if exists
		_, err := chainCfgRepo.FindOne(req.Chain)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		inVolumes, err := chainFlowCacheRepo.GetInflowTrend(constant.ChainFlowTrendDays, req.Chain)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		outVolumes, err := chainFlowCacheRepo.GetOutflowTrend(constant.ChainFlowTrendDays, req.Chain)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		return &vo.ChainVolumeTrendResp{
			VolumeIn:  fillVolumeItems(inVolumes),
			VolumeOut: fillVolumeItems(outVolumes),
			Chain:     req.Chain,
		}, nil
	}

	chainsCfg, err := chainCfgRepo.FindAllChainInfos()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	inVolumeMap := make(map[string]decimal.Decimal, 1)
	outVolumeMap := make(map[string]decimal.Decimal, 1)
	for _, val := range chainsCfg {
		inVolumes, err := chainFlowCacheRepo.GetInflowTrend(constant.ChainFlowTrendDays, val.ChainName)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		for _, volu := range inVolumes {
			value, _ := decimal.NewFromString(volu.Value)
			if data, ok := inVolumeMap[volu.Datetime]; ok {
				data = data.Add(value)
				inVolumeMap[volu.Datetime] = data
			} else {
				inVolumeMap[volu.Datetime] = value
			}
		}

		outVolumes, err := chainFlowCacheRepo.GetOutflowTrend(constant.ChainFlowTrendDays, val.ChainName)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		for _, volu := range outVolumes {
			value, _ := decimal.NewFromString(volu.Value)
			if data, ok := outVolumeMap[volu.Datetime]; ok {
				data = data.Add(value)
				outVolumeMap[volu.Datetime] = data
			} else {
				outVolumeMap[volu.Datetime] = value
			}
		}
	}

	date := time.Now().AddDate(0, 0, -constant.ChainFlowTrendDays+1)
	startUnix := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).Unix()
	inVolumes := make([]vo.VolumeItem, 0, constant.ChainFlowTrendDays)
	outVolumes := make([]vo.VolumeItem, 0, constant.ChainFlowTrendDays)
	for i := 0; i < constant.ChainFlowTrendDays; i++ {
		dt := time.Unix(startUnix+int64(i*86400), 0).Format(constant.DateFormat)
		inValue := "0"
		outValue := "0"
		if _, ok := inVolumeMap[dt]; ok {
			inValue = inVolumeMap[dt].String()
		}
		if _, ok := outVolumeMap[dt]; ok {
			outValue = outVolumeMap[dt].String()
		}
		inVolumes = append(inVolumes, vo.VolumeItem{
			Datetime: dt,
			Value:    inValue,
		})
		outVolumes = append(outVolumes, vo.VolumeItem{
			Datetime: dt,
			Value:    outValue,
		})
	}

	return &vo.ChainVolumeTrendResp{
		VolumeIn:  inVolumes,
		VolumeOut: outVolumes,
		Chain:     req.Chain,
	}, nil
}

func (svc *OverviewService) ChainVolume(req *vo.ChainVolumeReq) (*vo.ChainVolumeResp, errors.Error) {
	chainInVolumesMap, err := chainFlowCacheRepo.GetAllInflowVolume(constant.ChainFlowTrendDays)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	allInVolumes := float64(0)
	for _, val := range chainInVolumesMap {
		allInVolumes += val
	}

	chainOutVolumesMap, err := chainFlowCacheRepo.GetAllOutflowVolume(constant.ChainFlowTrendDays)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	allOutVolumes := float64(0)
	for _, val := range chainOutVolumesMap {
		allOutVolumes += val
	}

	chainsCfg, err := chainCfgRepo.FindAllChainInfos()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	resp := make(vo.ChainVolumeResp, 0, len(chainsCfg))
	resp = append(resp, vo.ChainVolumeItem{
		Chain:               "all_chain",
		TransferVolumeIn:    strconv.FormatFloat(allInVolumes, 'f', 4, 64),
		TransferVolumeOut:   strconv.FormatFloat(allOutVolumes, 'f', 4, 64),
		TransferVolumeTotal: strconv.FormatFloat(allInVolumes+allOutVolumes, 'f', 4, 64),
	})
	for _, val := range chainsCfg {
		inVolume := chainInVolumesMap[val.ChainName]
		outVolume := chainOutVolumesMap[val.ChainName]
		totalVolume := inVolume + outVolume
		item := vo.ChainVolumeItem{
			Chain:               val.ChainName,
			TransferVolumeIn:    strconv.FormatFloat(inVolume, 'f', 4, 64),
			TransferVolumeOut:   strconv.FormatFloat(outVolume, 'f', 4, 64),
			TransferVolumeTotal: strconv.FormatFloat(totalVolume, 'f', 4, 64),
		}
		resp = append(resp, item)
	}
	return &resp, nil
}
