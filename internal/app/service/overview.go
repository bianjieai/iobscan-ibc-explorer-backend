package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
)

type IOverviewService interface {
	MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error)
	TokenDistribution(req *vo.TokenDistributionReq) (*vo.TokenDistributionResp, errors.Error)
}

var _ IOverviewService = new(OverviewService)

type OverviewService struct {
}

func (t *OverviewService) MarketHeatmap() (*vo.MarketHeatmapResp, errors.Error) {
	return nil, nil
}

func (t *OverviewService) TokenDistribution(req *vo.TokenDistributionReq) (*vo.TokenDistributionResp, errors.Error) {
	return nil, nil
}