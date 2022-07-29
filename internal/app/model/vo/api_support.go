package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type (
	StatisticInfoResp struct {
		Items     []IbcStatisticDto `json:"items"`
		TimeStamp int64             `json:"time_stamp"`
	}

	IbcStatisticDto struct {
		StatisticsName string `json:"statistics_name"`
		Count          int64  `json:"count"`
	}
)

func (dto IbcStatisticDto) LoadDto(data *entity.IbcStatistic) IbcStatisticDto {
	return IbcStatisticDto{
		StatisticsName: data.StatisticsName,
		Count:          data.Count,
	}
}
