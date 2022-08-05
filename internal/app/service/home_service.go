package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"strings"
	"time"
)

type IHomeService interface {
	DailyChains() (vo.DailyChainsResp, errors.Error)
	IbcBaseDenoms() (vo.IbcBaseDenomsResp, errors.Error)
	IbcDenoms() (vo.IbcDenomsResp, errors.Error)
	Statistics() (vo.StatisticsCntResp, errors.Error)
}

var _ IHomeService = new(HomeService)

type HomeService struct {
	baseDenomdto    vo.IbcBaseDenomDto
	denomDto        vo.IbcDenomDto
	statisticCntDto vo.StatisticsCntDto
}

func (service HomeService) DailyChains() (vo.DailyChainsResp, errors.Error) {
	var resp vo.DailyChainsResp

	data, err := statisticRepo.FindOne(constant.Chains24hStatisticName)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	chainIds := strings.Split(data.StatisticsInfo, ",")
	activeChainsMap := make(map[string]struct{}, len(chainIds))
	for _, val := range chainIds {
		activeChainsMap[val] = struct{}{}
	}

	chainCfgs, err := chainCfgRepo.FindAllChainInfs()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	allChainsLen := len(chainCfgs)
	activeChainsLen := len(chainIds)
	allChains := make([]vo.DailyData, 0, len(chainCfgs))
	activeChains := make([]vo.DailyData, 0, len(chainIds))
	inActiveChains := make([]vo.DailyData, 0, allChainsLen-activeChainsLen)
	for _, one := range chainCfgs {
		item := vo.DailyData{
			ChainName: one.ChainName,
			ChainId:   one.ChainId,
			Icon:      one.Icon,
		}
		allChains = append(allChains, item)

		_, exist := activeChainsMap[one.ChainId]
		if exist {
			activeChains = append(activeChains, item)
		} else {
			inActiveChains = append(inActiveChains, item)
		}
	}
	resp.Items = vo.DailyChainDto{All: allChains, Active: activeChains, Inactive: inActiveChains}

	return resp, nil
}

func (service HomeService) IbcBaseDenoms() (vo.IbcBaseDenomsResp, errors.Error) {
	var resp vo.IbcBaseDenomsResp
	rets, err := baseDenomRepo.FindAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, service.baseDenomdto.LoadDto(val))
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (service HomeService) IbcDenoms() (vo.IbcDenomsResp, errors.Error) {
	var resp vo.IbcDenomsResp
	rets, err := denomRepo.FindSymbolDenoms()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, service.denomDto.LoadDto(val))
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (service HomeService) Statistics() (vo.StatisticsCntResp, errors.Error) {
	var resp vo.StatisticsCntResp
	rets, err := statisticRepo.FindBatchName(constant.HomeStatistics)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, service.statisticCntDto.LoadDto(val))
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
