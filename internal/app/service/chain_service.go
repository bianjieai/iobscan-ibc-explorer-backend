package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"time"
)

type IChainService interface {
	List(req *vo.ChainListReq) (vo.ChainListResp, errors.Error)
}

type ChainService struct {
	dto vo.ChainDto
}

var _ IChainService = new(ChainService)

func (svc *ChainService) List(req *vo.ChainListReq) (vo.ChainListResp, errors.Error) {
	var resp vo.ChainListResp
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	rets, err := chainRepo.FindAll(skip, limit)
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, svc.dto.LoadDto(val))
	}
	page := vo.BuildPageInfo(int64(len(rets)), req.PageNum, req.PageSize)
	resp.PageInfo = page
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
