package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IChainService interface {
	List() (*vo.ChainListResp, errors.Error)
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
