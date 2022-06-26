package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IChannelService interface {
	List(chainA, chainB string, status entity.ChannelStatus, useCount bool, pageNum, pageSize int64) (*vo.ChannelListResp, errors.Error)
}

type ChannelService struct {
}

func (svc *ChannelService) List(chainA, chainB string, status entity.ChannelStatus, useCount bool, pageNum, pageSize int64) (*vo.ChannelListResp, errors.Error) {
	list, err := channelRepo.List(chainA, chainB, status, pageNum, pageSize)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	items := make([]vo.ChannelItem, 0, len(list))
	for _, v := range list {
		items = append(items, vo.ChannelItem{
			ChainA:              v.ChainA,
			ChannelA:            v.ChannelA,
			ChainB:              v.ChainB,
			ChannelB:            v.ChannelB,
			OperatingPeriod:     v.OperatingPeriod,
			Relayers:            v.Relayers,
			LastUpdated:         v.ChannelUpdateAt,
			IbcTransferTxsValue: v.TransferTxsValue,
			IbcTransferTxs:      v.TransferTxs,
			Currency:            constant.DefaultCurrency,
			Status:              v.Status,
		})
	}

	var totalItem int64
	if useCount {
		totalItem, err = channelRepo.CountList(chainA, chainB, status)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	page := vo.BuildPageInfo(totalItem, pageNum, pageSize)
	return &vo.ChannelListResp{
		Items:    items,
		PageInfo: page,
	}, nil
}
