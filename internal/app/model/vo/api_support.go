package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
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

	FailTxsListReq struct {
		Page
	}

	FailTxsListResp struct {
		Items     []FailTxsListDto `json:"items"`
		PageInfo  PageInfo         `json:"page_info"`
		TimeStamp int64            `json:"time_stamp"`
	}
	FailTxsListDto struct {
		TxHash     string `json:"tx_hash"`
		Chain      string `json:"chain"`
		TxErrorLog string `json:"tx_error_log"`
		SendChain  string `json:"send_chain"`
		RecvChain  string `json:"recv_chain"`
	}

	RelayerTxFeesReq struct {
		Page
		TxHash string `form:"tx_hash" json:"tx_hash"`
		Chain  string `form:"chain" json:"chain"`
	}

	RelayerTxFeesResp struct {
		Items     []RelayerTxFeeDto `json:"items"`
		PageInfo  PageInfo          `json:"page_info"`
		TimeStamp int64             `json:"time_stamp"`
	}
	RelayerTxFeeDto struct {
		TxHash      string     `json:"tx_hash"`
		Chain       string     `json:"chain"`
		Fee         *model.Fee `json:"fee"`
		RelayerAddr string     `json:"relayer_addr"`
	}

	AccountsDailyResp struct {
		Items     []AccountsDailyDto `json:"items"`
		TimeStamp int64              `json:"time_stamp"`
		DateTime  string             `json:"date_time"`
	}
	AccountsDailyDto struct {
		ChainName string   `json:"chain_name"`
		Address   []string `json:"address"`
	}
)

func (dto IbcStatisticDto) LoadDto(data *entity.IbcStatistic) IbcStatisticDto {
	return IbcStatisticDto{
		StatisticsName: data.StatisticsName,
		Count:          data.Count + data.CountLatest,
	}
}
