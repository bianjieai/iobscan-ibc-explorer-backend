package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"time"
)

type IStatisticInfoService interface {
	IbcTxStatistic() (vo.StatisticInfoResp, errors.Error)
	AccountsDailyStatistic() (vo.AccountsDailyResp, errors.Error)
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

func (svc *StatisticInfoService) AccountsDailyStatistic() (vo.AccountsDailyResp, errors.Error) {
	var resp vo.AccountsDailyResp
	ret, err := statisticRepo.FindOne(constant.AccountsDailyStatisticName)
	if err != nil {
		return resp, errors.Wrap(err)
	}

	var data map[string][]string
	utils.UnmarshalJsonIgnoreErr([]byte(ret.Data), &data)

	datas, err := chainConfigRepo.FindAllChainIds()
	if err != nil {
		return resp, errors.Wrap(err)
	}
	chainNameMap := make(map[string]string, len(data))
	for _, val := range datas {
		chainNameMap[val.ChainId] = val.ChainName
	}

	resp.Items = make([]vo.AccountsDailyDto, 0, len(data))
	for chainId, val := range data {
		item := vo.AccountsDailyDto{
			Address:   val,
			ChainName: chainNameMap[chainId],
		}
		resp.Items = append(resp.Items, item)
	}
	resp.DateTime, _ = cache.GetRedisClient().Get(cache.DailyAccountsDate)
	resp.TimeStamp = time.Now().Unix()
	return resp, nil
}
