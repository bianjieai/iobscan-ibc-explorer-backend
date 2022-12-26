package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	v8 "github.com/go-redis/redis/v8"
	"strings"
)

type IOverviewService interface {
	MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error)
	TokenDistribution(req *vo.TokenDistributionReq) (*vo.TokenDistributionResp, errors.Error)
}

var _ IOverviewService = new(OverviewService)

type OverviewService struct {
}

func (t *OverviewService) MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error) {
	return nil, nil
}

func (t *OverviewService) TokenDistribution(req *vo.TokenDistributionReq) (*vo.TokenDistributionResp, errors.Error) {
	ibcDenoms, err := denomRepo.FindByBaseDenom(req.BaseDenom, req.BaseDenomChain)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	getHops := func(denomPath string) int {
		return strings.Count(denomPath, "/channel")
	}
	mapHopsData := make(map[int]entity.IBCDenomList, 1)
	mapChainData := make(map[string]string, 1)
	for _, val := range ibcDenoms {

		amount, err := supportCache.GetSupply(val.Chain, val.Denom)
		if err != nil {
			if err == v8.Nil {
				amount = constant.ZeroDenomAmount
			}
			amount = constant.UnknownDenomAmount
		}
		_, ok := mapChainData[val.Chain+val.Denom]
		if !ok {
			mapChainData[val.Chain+val.Denom] = amount
		}

		//todo replace with val.hops
		hop := getHops(val.DenomPath)

		if val.DenomPath == "" || hop == 0 {
			continue
		}

		hopDatas, exist := mapHopsData[hop]
		if exist {
			hopDatas = append(hopDatas, val)
			mapHopsData[hop] = hopDatas
		} else {
			mapHopsData[hop] = entity.IBCDenomList{val}
		}
	}

	resp := &vo.TokenDistributionResp{
		Chain: req.BaseDenomChain,
		Denom: req.BaseDenom,
		Hops:  0,
	}
	//hop get ibc denom
	hop := 1
	resp.Children = make([]*vo.GraphData, 0, 1)
	hopDenoms, ok := mapHopsData[hop]
	if !ok {
		return resp, nil
	}
	for _, val := range hopDenoms {

		children := vo.GraphData{
			Denom:  val.Denom,
			Chain:  val.Chain,
			Hops:   hop,
			Amount: mapChainData[val.Chain+val.Denom],
		}
		children = t.FindSource(mapChainData, mapHopsData, children)

		resp.Children = append(resp.Children, &children)
	}

	//todo caculate transfer amount from children to root

	return resp, nil
}

func (t *OverviewService) FindSource(mapChainData map[string]string, mapHopsData map[int]entity.IBCDenomList, ret vo.GraphData) vo.GraphData {
	hopDenoms, ok := mapHopsData[ret.Hops+1]
	if !ok {
		return ret
	}
	ret.Children = make([]*vo.GraphData, 0, 1)
	for _, val := range hopDenoms {
		if val.PrevDenom == ret.Denom && val.PrevChain == ret.Chain {
			children := vo.GraphData{
				Denom:  val.Denom,
				Chain:  val.Chain,
				Hops:   ret.Hops + 1,
				Amount: mapChainData[val.Chain+val.Denom],
			}
			children = t.FindSource(mapChainData, mapHopsData, children)
			ret.Children = append(ret.Children, &children)
		}
	}

	return ret
}
