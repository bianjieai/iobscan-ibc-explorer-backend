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
		arrs := strings.Split(val.RelayerChannelPair, ":")
		if len(arrs) == 6 {
			srcChainInfo := strings.Join([]string{arrs[0], arrs[1], arrs[2]}, ":")
			dcChainInfo := strings.Join([]string{arrs[3], arrs[4], arrs[5]}, ":")
			relayerCfgMap[srcChainInfo] = *val
			relayerCfgMap[dcChainInfo] = *val
		}

	}
	for _, val := range rets {
		item := svc.dto.LoadDto(val)
		srcChainInfo := strings.Join([]string{val.ChainA, val.ChannelA, val.ChainAAddress}, ":")
		dcChainInfo := strings.Join([]string{val.ChainB, val.ChannelB, val.ChainBAddress}, ":")
		if cfg, ok := relayerCfgMap[srcChainInfo]; ok {
			item.RelayerName = cfg.RelayerName
			item.RelayerIcon = cfg.Icon
		} else if cfg, ok = relayerCfgMap[dcChainInfo]; ok {
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

func (svc *RelayerService) ListCount(req *vo.RelayerListReq) (int64, errors.Error) {
	total, err := relayerRepo.CountBycond(req.Chain, req.Status)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return total, nil
}
