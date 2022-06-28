package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
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
		relayerCfgMap[val.RelayerId] = *val
	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		if cfg, ok := relayerCfgMap[val.RelayerId]; ok {
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
