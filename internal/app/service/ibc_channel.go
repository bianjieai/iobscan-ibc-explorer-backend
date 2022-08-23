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
	ListCount(req *vo.ChannelListReq) (int64, errors.Error)
}

var _ IChannelService = new(ChannelService)

type ChannelService struct {
}

func (svc *ChannelService) List(req *vo.ChannelListReq) (*vo.ChannelListResp, errors.Error) {
	chainA, chainB, err := svc.analyzeChain(req.Chain)
	if err != nil {
		return nil, errors.Wrap(err)
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

func (svc *ChannelService) analyzeChain(chain string) (string, string, error) {
	if chain == "" {
		return constant.AllChain, constant.AllChain, nil
	} else {
		split := strings.Split(chain, ",")
		if len(split) != 2 {
			return "", "", fmt.Errorf("chain parameter format error")
		} else {
			return split[0], split[1], nil
		}
	}
}

func (svc *ChannelService) ListCount(req *vo.ChannelListReq) (int64, errors.Error) {
	chainA, chainB, err := svc.analyzeChain(req.Chain)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	totalItem, err := channelRepo.CountList(chainA, chainB, req.Status)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return totalItem, nil
}
