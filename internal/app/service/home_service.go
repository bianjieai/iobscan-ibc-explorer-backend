package service

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IHomeService interface {
	ChainsConnection() (vo.ChainsConnectionResp, errors.Error)
	DailyChains() (vo.DailyChainsResp, errors.Error)
	AuthDenoms() (vo.AuthDenomsResp, errors.Error)
	IbcDenoms() (vo.IbcDenomsResp, errors.Error)
	Statistics() (vo.StatisticsCntResp, errors.Error)
	SearchPoint(req *vo.SearchPointReq) errors.Error
}

var _ IHomeService = new(HomeService)

type HomeService struct {
	authDenomdto    vo.AuthDenomDto
	denomDto        vo.IbcDenomDto
	statisticCntDto vo.StatisticsCntDto
}

func (svc HomeService) ChainsConnection() (vo.ChainsConnectionResp, errors.Error) {
	var resp vo.ChainsConnectionResp

	if value, err := chainCache.GetChainsConnection(); err == nil {
		utils.UnmarshalJsonIgnoreErr([]byte(value), &resp.Items)
		resp.TimeStamp = time.Now().Unix()
		return resp, nil
	}

	chainCfgs, err := chainCfgRepo.FindAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}

	getIobConnectionChain := func(ibcInfo *entity.IbcInfo) vo.IobConnectionChain {
		status := entity.ChannelStatusClosed
		for _, val := range ibcInfo.Paths {
			if val.State == constant.ChannelStateOpen && val.Counterparty.State == constant.ChannelStateOpen {
				status = entity.ChannelStatusOpened
				break
			}
		}
		return vo.IobConnectionChain{
			ChainName:        ibcInfo.Chain,
			ConnectionStatus: status,
		}
	}

	iobChains := make([]vo.IobChainDto, 0, len(chainCfgs))
	for _, one := range chainCfgs {
		item := vo.IobChainDto{
			ChainName:      one.ChainName,
			PrettyName:     one.PrettyName,
			CurrentChainId: one.CurrentChainId,
			Icon:           fmt.Sprintf(constant.IBCConnectionChainsIconUri, one.ChainName),
		}
		connectionChains := make([]vo.IobConnectionChain, 0, len(one.IbcInfo))
		for _, val := range one.IbcInfo {
			connectChain := getIobConnectionChain(val)
			connectionChains = append(connectionChains, connectChain)
		}
		item.ConnectionChains = connectionChains

		iobChains = append(iobChains, item)
	}
	if len(iobChains) > 0 {
		_ = chainCache.SetChainsConnection(string(utils.MarshalJsonIgnoreErr(iobChains)))
	}
	resp.Items = iobChains
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
func (svc HomeService) DailyChains() (vo.DailyChainsResp, errors.Error) {
	var resp vo.DailyChainsResp

	data, err := statisticRepo.FindOne(constant.Chains24hStatisticName)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	chains := strings.Split(data.StatisticsInfo, ",")
	activeChainsMap := make(map[string]struct{}, len(chains))
	for _, val := range chains {
		activeChainsMap[val] = struct{}{}
	}

	chainCfgs, err := chainCfgRepo.FindAllChainInfos()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	allChainsLen := len(chainCfgs)
	activeChainsLen := len(chains)
	allChains := make([]vo.DailyData, 0, len(chainCfgs))
	activeChains := make([]vo.DailyData, 0, len(chains))
	inActiveChains := make([]vo.DailyData, 0, allChainsLen-activeChainsLen)
	for _, one := range chainCfgs {
		item := vo.DailyData{
			ChainName:      one.ChainName,
			PrettyName:     one.PrettyName,
			AddrPrefix:     []string{one.AddrPrefix},
			CurrentChainId: one.CurrentChainId,
			Icon:           one.Icon,
			Status:         one.Status,
		}
		allChains = append(allChains, item)

		_, exist := activeChainsMap[one.ChainName]
		if exist {
			activeChains = append(activeChains, item)
		} else {
			inActiveChains = append(inActiveChains, item)
		}
	}
	resp.Items = []vo.DailyChainDto{{All: allChains, Active: activeChains, Inactive: inActiveChains}}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc HomeService) AuthDenoms() (vo.AuthDenomsResp, errors.Error) {
	var resp vo.AuthDenomsResp
	rets, err := authDenomRepo.FindAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, svc.authDenomdto.LoadDto(val))
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc HomeService) IbcDenoms() (vo.IbcDenomsResp, errors.Error) {
	var resp vo.IbcDenomsResp
	rets, err := denomRepo.FindSymbolDenoms()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, svc.denomDto.LoadDto(val))
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc HomeService) Statistics() (vo.StatisticsCntResp, errors.Error) {
	var resp vo.StatisticsCntResp
	rets, err := statisticRepo.FindBatchName(constant.HomeStatistics)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, svc.statisticCntDto.LoadDto(val))
	}

	relayers, err := relayerRepo.CountAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}

	resp.Items = append(resp.Items, vo.StatisticsCntDto{
		StatisticsName: constant.RelayersStatisticName,
		Count:          relayers,
	})
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc HomeService) SearchPoint(req *vo.SearchPointReq) errors.Error {
	if err := exSearchRecordRepo.Insert(&entity.UbaSearchRecord{
		Ip:       req.Ip,
		Content:  req.Content,
		CreateAt: time.Now().Unix(),
	}); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
