package service

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"

type IChainService interface {
	List(req *vo.ChainListReq) (vo.ChainListResp, error)
}

type ChainService struct {
	dto vo.ChainDto
}

var _ IChainService = new(ChainService)

func (svc *ChainService) List(req *vo.ChainListReq) (vo.ChainListResp, error) {
	var resp vo.ChainListResp
	//todo current no use request data
	rets, err := chainRepo.FindAll()
	if err != nil {
		return resp, err
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, svc.dto.LoadDto(val))
	}
	resp.PageInfo = vo.PageInfo{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	resp.Total = len(rets)
	return resp, nil
}
