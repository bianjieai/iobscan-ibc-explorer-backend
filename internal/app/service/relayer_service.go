package service

import (
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error)
	ListCount(req *vo.RelayerListReq) (int64, errors.Error)
	Collect(OperatorFile string) errors.Error
}

type RelayerService struct {
	dto            vo.RelayerDto
	relayerHandler RelayerHandler
}

var _ IRelayerService = new(RelayerService)

func (svc *RelayerService) List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error) {
	var resp vo.RelayerListResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	chains := strings.Split(req.Chain, ",")
	//unsupport more than two chains
	if len(chains) > 2 {
		return resp, nil
	}
	rets, total, err := relayerRepo.FindAllBycond(req.Chain, req.Status, skip, limit, req.UseCount)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	relayerCfgs, err := relayerCfgRepo.FindAll()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	relayerCfgMap := make(map[string]*entity.IBCRelayerConfig, len(relayerCfgs))
	for _, val := range relayerCfgs {
		relayerCfgMap[val.RelayerPairId] = val
	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		chainA := strings.ReplaceAll(val.ChainA, "_", "-")
		chainB := strings.ReplaceAll(val.ChainB, "_", "-")
		pairId := entity.GenerateRelayerPairId(chainA, val.ChannelA, val.ChainAAddress, chainB, val.ChannelB, val.ChainBAddress)
		config, ok := relayerCfgMap[pairId]
		if ok {
			item.RelayerName = config.RelayerName
			item.RelayerIcon = config.Icon
		}
		resp.Items = append(resp.Items, item)
	}
	page := vo.BuildPageInfo(total, req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}

func (svc *RelayerService) ListCount(req *vo.RelayerListReq) (int64, errors.Error) {
	total, err := relayerRepo.CountBycond(req.Chain, req.Status)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return total, nil
}

func (svc *RelayerService) Collect(OperatorFile string) errors.Error {
	go svc.relayerHandler.Collect(OperatorFile)
	return nil
}
