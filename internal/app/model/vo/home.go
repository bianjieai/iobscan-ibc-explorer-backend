package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

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
		CurrentChainId string             `json:"current_chain_id"`
		ChainName      string             `json:"chain_name"`
		PrettyName     string             `json:"pretty_name"`
		Icon           string             `json:"icon"`
		Status         entity.ChainStatus `json:"status"`
	}

	IbcBaseDenomsResp struct {
		Items     []AuthDenomDto `json:"items"`
		TimeStamp int64          `json:"time_stamp"`
	}
	AuthDenomDto struct {
		Chain  string `json:"chain"`
		Denom  string `json:"denom"`
		Symbol string `json:"symbol"`
		Scale  int    `json:"scale"`
		Icon   string `json:"icon"`
	}
	IbcDenomsResp struct {
		Items     []IbcDenomDto `json:"items"`
		TimeStamp int64         `json:"time_stamp"`
	}

	IbcDenomDto struct {
		Chain          string `json:"chain"`
		Denom          string `json:"denom"`
		BaseDenom      string `json:"base_denom"`
		BaseDenomChain string `json:"base_denom_chain"`
		DenomPath      string `json:"denom_path"`
		Symbol         string `json:"symbol"`
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

func (dto AuthDenomDto) LoadDto(baseDenom *entity.AuthDenom) AuthDenomDto {
	return AuthDenomDto{
		Chain:  baseDenom.Chain,
		Denom:  baseDenom.Denom,
		Symbol: baseDenom.Symbol,
		Scale:  baseDenom.Scale,
		Icon:   baseDenom.Icon,
	}
}

func (dto IbcDenomDto) LoadDto(denom *entity.IBCDenom) IbcDenomDto {
	return IbcDenomDto{
		Chain:          denom.Chain,
		Denom:          denom.Denom,
		BaseDenom:      denom.BaseDenom,
		BaseDenomChain: denom.BaseDenomChain,
		DenomPath:      denom.DenomPath,
		Symbol:         denom.Symbol,
	}
}

func (dto StatisticsCntDto) LoadDto(statistic *entity.IbcStatistic) StatisticsCntDto {
	return StatisticsCntDto{
		Count:          statistic.Count + statistic.CountLatest,
		StatisticsName: statistic.StatisticsName,
	}
}

type SearchPointReq struct {
	Content string `json:"content"`
	Ip      string `json:"ip"`
}
