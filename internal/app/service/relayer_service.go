package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IRelayerService interface {
	List(req *vo.RelayerListReq) (vo.RelayerListResp, error)
}

type RelayerService struct {
	dto vo.RelayerDto
}

var _ IRelayerService = new(RelayerService)

func (svc *RelayerService) List(req *vo.RelayerListReq) (vo.RelayerListResp, error) {
	var resp vo.RelayerListResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	rets, total, err := relayerRepo.FindAllBycond(req.Chain, req.Status, skip, limit, req.UseCount)
	if err != nil {
		return resp, err
	}
	relayerCfgs, err := relayerCfgRepo.FindAll()
	if err != nil {
		return resp, err
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
	resp.PageInfo = vo.PageInfo{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	resp.Total = total
	return resp, nil
}
