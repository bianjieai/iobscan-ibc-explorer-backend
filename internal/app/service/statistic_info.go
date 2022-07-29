package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"time"
)

type IStatisticInfoService interface {
	IbcTxStatistic() (vo.StatisticInfoResp, errors.Error)
}

type StatisticInfoService struct {
	dto vo.IbcStatisticDto
}

func (svc *StatisticInfoService) IbcTxStatistic() (vo.StatisticInfoResp, errors.Error) {
	var resp vo.StatisticInfoResp
	rets, err := statisticRepo.FindBatchName([]string{constant.TxALlStatisticName, constant.TxFailedStatisticName})
	if err != nil {
		return resp, errors.Wrap(err)
	}
	for _, val := range rets {
		resp.Items = append(resp.Items, svc.dto.LoadDto(val))
	}
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
