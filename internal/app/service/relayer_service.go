package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"strings"
	"time"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, errors.Error)
}

type RelayerService struct {
	dto vo.RelayerDto
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
	relayerCfgMap := make(map[string]entity.IBCRelayerConfig, len(relayerCfgs))
	for _, val := range relayerCfgs {
		relayerCfgMap[val.RelayerChannelPair] = *val
	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		relayerChannelPairA, relayerChannelPairB := repository.CreateRelayerChannelPair(val.ChainA, val.ChainB, val.ChannelA, val.ChannelB, val.ChainAAddress, val.ChainBAddress)
		if cfg, ok := relayerCfgMap[relayerChannelPairA]; ok {
			item.RelayerName = cfg.RelayerName
			item.RelayerIcon = cfg.Icon
		} else if cfg, ok = relayerCfgMap[relayerChannelPairB]; ok {
			item.RelayerName = cfg.RelayerName
			item.RelayerIcon = cfg.Icon
		}
		resp.Items = append(resp.Items, item)
	}
	page := vo.BuildPageInfo(total, req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
