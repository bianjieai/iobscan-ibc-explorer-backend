package service

import (
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error)
	ListCount(req *vo.RelayerListReq) (int64, errors.Error)
	Collect(OperatorFile string) errors.Error
	Detail(relayerId string) (vo.RelayerDetailResp, errors.Error)
}

type RelayerService struct {
	dto            vo.RelayerDto
	relayerHandler RelayerHandler
}

var _ IRelayerService = new(RelayerService)

func (svc *RelayerService) List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error) {
	var resp vo.RelayerListResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	rets, total, err := relayerRepo.FindAllBycond(req.RelayerName, req.RelayerAddress, skip, limit, req.UseCount)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		resp.Items = append(resp.Items, item)
	}
	page := vo.BuildPageInfo(total, req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc *RelayerService) ListCount(req *vo.RelayerListReq) (int64, errors.Error) {
	total, err := relayerRepo.CountBycond(req.RelayerName, req.RelayerAddress)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return total, nil
}

func (svc *RelayerService) Collect(OperatorFile string) errors.Error {
	go svc.relayerHandler.Collect(OperatorFile)
	return nil
}

func (svc *RelayerService) Detail(relayerId string) (vo.RelayerDetailResp, errors.Error) {
	var resp vo.RelayerDetailResp
	one, err := relayerRepo.FindOneByRelayerId(relayerId)
	if err != nil {
		return resp, errors.Wrap(err)
	}

	channelPairs, err := channelRepo.FindAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	channelPairStatusMap := make(map[string]int, len(channelPairs))
	for _, val := range channelPairs {
		channelPairStatusMap[val.ChainA+val.ChannelA+val.ChainB+val.ChannelB] = int(val.Status)
	}

	resp = vo.LoadRelayerDetailDto(one, channelPairStatusMap)

	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
