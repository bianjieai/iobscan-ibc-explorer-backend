package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/qiniu/qmgo"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IChainService interface {
	List() (*vo.ChainListResp, errors.Error)
	ChainExists(chain string) (bool, errors.Error)
	IbcChainsNum() (*vo.IbcChainsNumResp, errors.Error)
	IbcChainsVolume(chainName string) (*vo.IbcChainsVolumeResp, errors.Error)
	IbcChainsActive() (*vo.IbcChainsActiveResp, errors.Error)
}

var _ IChainService = new(ChainService)

type ChainService struct {
}

func (svc *ChainService) List() (*vo.ChainListResp, errors.Error) {
	chainList, err := chainConfigRepo.FindAllChainInfos()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	registryList, err := chainRegistryRepo.FindAll()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var registryMap = make(map[string]string, len(registryList))
	for _, v := range registryList {
		registryMap[v.Chain] = v.ChainJsonUrl
	}
	chainItems := make([]vo.ChainItem, 0, len(chainList))
	for _, c := range chainList {
		chainItems = append(chainItems, vo.ChainItem{
			Chain:            c.ChainName,
			ChainRegistryUrl: registryMap[c.ChainName],
		})
	}

	return &vo.ChainListResp{
		Items:   chainItems,
		Comment: constant.ContactUs,
	}, nil
}

func (svc *ChainService) ChainExists(chain string) (bool, errors.Error) {
	_, err := chainConfigRepo.FindOne(chain)
	if err == qmgo.ErrNoSuchDocuments {
		return false, nil
	}
	if err != nil {
		return true, errors.Wrap(err)
	}
	return true, nil
}

func (svc *ChainService) IbcChainsNum() (*vo.IbcChainsNumResp, errors.Error) {
	res, err := ibcChainRepo.CountIbcChainsNum()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var resp vo.IbcChainsNumResp
	resp.IbcChainsNumber = res
	return &resp, nil
}

func (svc *ChainService) IbcChainsVolume(chainName string) (*vo.IbcChainsVolumeResp, errors.Error) {
	segmentStartTime, _ := utils.TodayUnix()
	res, err := ibcChainInflowStatisticsRepo.GetLatestUpdate(segmentStartTime)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			res.UpdateAt = time.Now().Unix()
		} else {
			return nil, errors.Wrap(err)
		}
	}

	if chainName == "" {
		chainInVolumesMap, err := chainFlowCacheRepo.GetAllInflowVolume(constant.ChainFlowVolumeDays)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		chainOutVolumesMap, err := chainFlowCacheRepo.GetAllOutflowVolume(constant.ChainFlowVolumeDays)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		chainsCfg, err := chainConfigRepo.FindAllChainInfos()
		if err != nil {
			return nil, errors.Wrap(err)
		}

		items := make([]vo.IbcChainVolume, 0, len(chainsCfg))
		for _, val := range chainsCfg {
			inVolume := chainInVolumesMap[val.ChainName]
			outVolume := chainOutVolumesMap[val.ChainName]
			totalVolume := inVolume + outVolume
			item := vo.IbcChainVolume{
				ChainName:              val.ChainName,
				IbcVolumeIn:            strconv.FormatFloat(inVolume, 'f', 4, 64),
				IbcVolumeOut:           strconv.FormatFloat(outVolume, 'f', 4, 64),
				IbcTransferVolumeTotal: strconv.FormatFloat(totalVolume, 'f', 4, 64),
			}
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			iv, _ := strconv.ParseFloat(items[i].IbcTransferVolumeTotal, 64)
			jv, _ := strconv.ParseFloat(items[j].IbcTransferVolumeTotal, 64)
			return iv > jv
		})
		var resp vo.IbcChainsVolumeResp
		resp.Chains = items
		resp.TimeStamp = res.UpdateAt
		return &resp, nil
	} else {
		inVolume, err := chainFlowCacheRepo.GetInflowVolume(constant.ChainFlowVolumeDays, chainName)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		outVolume, err := chainFlowCacheRepo.GetOutflowVolume(constant.ChainFlowVolumeDays, chainName)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		inVolumeF, _ := strconv.ParseFloat(inVolume, 64)
		outVolumeF, _ := strconv.ParseFloat(outVolume, 64)
		var items []vo.IbcChainVolume
		item := vo.IbcChainVolume{
			ChainName:              chainName,
			IbcVolumeIn:            inVolume,
			IbcVolumeOut:           outVolume,
			IbcTransferVolumeTotal: strconv.FormatFloat(inVolumeF+outVolumeF, 'f', 4, 64),
		}
		items = append(items, item)
		var resp vo.IbcChainsVolumeResp
		resp.Chains = items
		resp.TimeStamp = res.UpdateAt
		return &resp, nil
	}
}

func (svc *ChainService) IbcChainsActive() (*vo.IbcChainsActiveResp, errors.Error) {
	data, err := statisticRepo.FindOne(constant.Chains24hStatisticName)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var resp vo.IbcChainsActiveResp
	if data.StatisticsInfo != "" {
		resp.ChainNameList = strings.Split(data.StatisticsInfo, ",")
	} else {
		resp.ChainNameList = []string{}
	}
	resp.TotalActiveChainsNumber = len(resp.ChainNameList)
	return &resp, nil
}
