package service

import (
	"fmt"
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IChannelService interface {
	List(req *vo.ChannelListReq) (*vo.ChannelListResp, errors.Error)
}

var _ IChannelService = new(ChannelService)

type ChannelService struct {
}

func (svc *ChannelService) List(req *vo.ChannelListReq) (*vo.ChannelListResp, errors.Error) {
	var chainA, chainB string
	if req.Chain == "" {
		chainA = constant.AllChain
		chainB = constant.AllChain
	} else {
		split := strings.Split(req.Chain, ",")
		if len(split) != 2 {
			return nil, errors.WrapBadRequest(fmt.Errorf("chain parameter format error"))
		} else {
			chainA = split[0]
			chainB = split[1]
		}
	}

	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	list, err := channelRepo.List(chainA, chainB, req.Status, skip, limit)
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
	if req.UseCount {
		totalItem, err = channelRepo.CountList(chainA, chainB, req.Status)
		if err != nil {
			return nil, errors.Wrap(err)
		}
	}

	page := vo.BuildPageInfo(totalItem, req.PageNum, req.PageSize)
	return &vo.ChannelListResp{
		Items:    items,
		PageInfo: page,
	}, nil
}
