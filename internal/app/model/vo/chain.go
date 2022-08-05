package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type ChainListReq struct {
	Page
	UseCount bool `json:"use_count" form:"use_count"`
}
type ChainDto struct {
	ChainId          string `json:"chain_id"`
	ConnectedChains  int64  `json:"connected_chains"`
	Channels         int64  `json:"channels"`
	Relayers         int64  `json:"relayers"`
	IbcTokens        int64  `json:"ibc_tokens"`
	IbcTokensValue   string `json:"ibc_tokens_value"`
	TransferTxs      int64  `json:"transfer_txs"`
	TransferTxsValue string `json:"transfer_txs_value"`
	Currency         string `json:"currency"`
}

type ChainListResp struct {
	Items     []ChainDto `json:"items"`
	PageInfo  PageInfo   `json:"page_info"`
	TimeStamp int64      `json:"time_stamp"`
}

func (dto ChainDto) LoadDto(chain *entity.IBCChain) ChainDto {
	return ChainDto{
		ChainId:          chain.ChainId,
		ConnectedChains:  chain.ConnectedChains,
		Channels:         chain.Channels,
		Relayers:         chain.Relayers,
		IbcTokens:        chain.IbcTokens,
		IbcTokensValue:   chain.IbcTokensValue,
		TransferTxs:      chain.TransferTxs,
		TransferTxsValue: chain.TransferTxsValue,
		Currency:         constant.DefaultCurrency,
	}
}

type (
	DailyChainsResp struct {
		Items     []DailyChainDto `json:"items"`
		TimeStamp int64           `json:"time_stamp"`
	}
	DailyChainDto struct {
		All      []DailyData `json:"all"`
		Active   []DailyData `json:"active"`
		Inactive []DailyData `json:"inactive"`
	}
	DailyData struct {
		ChainId   string `json:"chain_id"`
		ChainName string `json:"chain_name"`
		Icon      string `json:"icon"`
	}

	IbcBaseDenomsResp struct {
		Items     []IbcBaseDenomDto `json:"items"`
		TimeStamp int64             `json:"time_stamp"`
	}
	IbcBaseDenomDto struct {
		ChainId string `json:"chain_id"`
		Denom   string `json:"denom"`
		Symbol  string `json:"symbol"`
		Scale   int    `json:"scale"`
		Icon    string `json:"icon"`
	}
	IbcDenomsResp struct {
		Items     []IbcDenomDto `json:"items"`
		TimeStamp int64         `json:"time_stamp"`
	}

	IbcDenomDto struct {
		ChainId   string `json:"chain_id"`
		Denom     string `json:"denom"`
		BaseDenom string `json:"base_denom"`
		DenomPath string `json:"denom_path"`
		Symbol    string `json:"symbol"`
	}

	StatisticsCntResp struct {
		Items     []StatisticsCntDto `json:"items"`
		TimeStamp int64              `json:"time_stamp"`
	}
	StatisticsCntDto struct {
		StatisticsName string `json:"statistics_name"`
		Count          int64  `json:"count"`
	}
)

func (dto IbcBaseDenomDto) LoadDto(baseDenom *entity.IBCBaseDenom) IbcBaseDenomDto {
	return IbcBaseDenomDto{
		ChainId: baseDenom.ChainId,
		Denom:   baseDenom.Denom,
		Symbol:  baseDenom.Symbol,
		Scale:   baseDenom.Scale,
		Icon:    baseDenom.Icon,
	}
}

func (dto IbcDenomDto) LoadDto(denom *entity.IBCDenom) IbcDenomDto {
	return IbcDenomDto{
		ChainId:   denom.ChainId,
		Denom:     denom.Denom,
		BaseDenom: denom.BaseDenom,
		DenomPath: denom.DenomPath,
		Symbol:    denom.Symbol,
	}
}

func (dto StatisticsCntDto) LoadDto(statistic *entity.IbcStatistic) StatisticsCntDto {
	return StatisticsCntDto{
		Count:          statistic.Count,
		StatisticsName: statistic.StatisticsName,
	}
}
