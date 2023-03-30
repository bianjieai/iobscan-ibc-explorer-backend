package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
)

type IAddressService interface {
	ChainAddressStatistics(chain string, startTime, endTime int64) (*vo.ChainAddressStatisticsResp, errors.Error)
}

var _ IAddressService = new(AddressService)

type AddressService struct {
}

func (svc *AddressService) ChainAddressStatistics(chain string, startTime, endTime int64) (*vo.ChainAddressStatisticsResp, errors.Error) {
	var resp vo.ChainAddressStatisticsResp
	res, err := ibcChainAddressStatisticsRepo.AddressStatistics(chain, startTime, endTime)
	if err == qmgo.ErrNoSuchDocuments {
		resp.ActiveAddressNumber = 0
		return &resp, nil
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}
	resp.ActiveAddressNumber = res.AddressAmount
	return &resp, nil
}
