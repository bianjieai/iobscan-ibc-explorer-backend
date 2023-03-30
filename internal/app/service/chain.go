package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
)

type IChainService interface {
	List() (*vo.ChainListResp, errors.Error)
	ChainExists(chain string) (bool, errors.Error)
	ActiveChainNum() (*vo.ActiveChainNumResp, errors.Error)
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

func (svc *ChainService) ActiveChainNum() (*vo.ActiveChainNumResp, errors.Error) {
	res, err := ibcChainRepo.CountActiveChainNum()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var resp vo.ActiveChainNumResp
	resp.ActiveChainNumber = res
	return &resp, nil
}
