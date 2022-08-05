package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
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

	TranaferTxsReq struct {
		Page
		UseCount  bool   `json:"use_count" form:"use_count"`
		DateRange string `json:"date_range" form:"date_range"`
		Status    string `json:"status" form:"status"`
		ChainId   string `json:"chain_id" form:"chain_id"`
		Symbol    string `json:"symbol" form:"symbol"`
		Denom     string `json:"denom" form:"denom"`
	}
	TranaferTxsResp struct {
		Items     []IbcTxDto `json:"items"`
		PageInfo  PageInfo   `json:"page_info"`
		TimeStamp int64      `json:"time_stamp"`
	}

	TranaferTxDetailResp struct {
		Items     []IbcTxDetailDto `json:"items"`
		TimeStamp int64            `json:"time_stamp"`
	}
	IbcTxDto struct {
		ScAddr    string    `json:"sc_addr"`
		DcAddr    string    `json:"dc_addr"`
		Status    int       `json:"status"`
		ScChainId string    `json:"sc_chain_id"`
		DcChainId string    `json:"dc_chain_id"`
		ScChannel string    `json:"sc_channel"`
		DcChannel string    `json:"dc_channel"`
		Sequence  string    `json:"sequence"`
		ScTxInfo  TxInfoDto `json:"sc_tx_info"`
		DcTxInfo  TxInfoDto `json:"dc_tx_info"`
		BaseDenom string    `json:"base_denom"`
		Denoms    Denoms    `json:"denoms"`
		TxTime    int64     `json:"tx_time"`
		EndTime   int64     `json:"end_time"`
	}
	Denoms struct {
		ScDenom string `json:"sc_denom"`
		DcDenom string `json:"dc_denom"`
	}

	IbcTxDetailDto struct {
		ScSigners        []string  `json:"sc_signers"`
		DcSigners        []string  `json:"dc_signers"`
		ScAddr           string    `json:"sc_addr"`
		DcAddr           string    `json:"dc_addr"`
		Status           int       `json:"status"`
		ScChainId        string    `json:"sc_chain_id"`
		ScChannel        string    `json:"sc_channel"`
		ScPort           string    `json:"sc_port"`
		ScConnection     string    `json:"sc_connection"`
		DcChainId        string    `json:"dc_chain_id"`
		DcChannel        string    `json:"dc_channel"`
		DcPort           string    `json:"dc_port"`
		DcConnection     string    `json:"dc_connection"`
		Sequence         string    `json:"sequence"`
		ScTxInfo         TxInfoDto `json:"sc_tx_info"`
		DcTxInfo         TxInfoDto `json:"dc_tx_info"`
		BaseDenom        string    `json:"base_denom"`
		Denoms           Denoms    `json:"denoms"`
		TxTime           int64     `json:"tx_time"`
		Ack              string    `json:"ack"`
		TimeoutTimestamp string    `json:"timeout_timestamp"`
	}
	TxInfoDto struct {
		Hash      string       `json:"hash,omitempty"`
		Status    int          `json:"status,omitempty"`
		Time      int64        `json:"time,omitempty"`
		Height    int64        `json:"height,omitempty"`
		Fee       *model.Fee   `json:"fee,omitempty"`
		MsgAmount *model.Coin  `json:"msg_amount,omitempty"`
		Msg       *model.TxMsg `json:"msg,omitempty"`
	}
)

func loadTxInfoDto(info *entity.TxInfo) TxInfoDto {
	if info == nil {
		return TxInfoDto{}
	}
	return TxInfoDto{
		Hash:      info.Hash,
		Status:    int(info.Status),
		Time:      info.Time,
		Height:    info.Height,
		Fee:       info.Fee,
		MsgAmount: info.MsgAmount,
		Msg:       info.Msg,
	}
}
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

func (dto IbcTxDto) LoadDto(ibcTx *entity.ExIbcTx) IbcTxDto {
	endTime := int64(0)
	switch ibcTx.Status {
	case entity.IbcTxStatusSuccess:
		if ibcTx.DcTxInfo != nil {
			endTime = ibcTx.DcTxInfo.Time
		}
	case entity.IbcTxStatusFailed:
		if ibcTx.ScTxInfo.Status == entity.TxStatusFailed {
			if ibcTx.ScTxInfo != nil {
				endTime = ibcTx.ScTxInfo.Time
			}
		} else {
			if ibcTx.RefundedTxInfo != nil {
				endTime = ibcTx.RefundedTxInfo.Time
			}
		}
	case entity.IbcTxStatusRefunded:
		if ibcTx.RefundedTxInfo != nil {
			endTime = ibcTx.RefundedTxInfo.Time
		}
	}
	return IbcTxDto{
		ScAddr:    ibcTx.ScAddr,
		DcAddr:    ibcTx.DcAddr,
		Status:    int(ibcTx.Status),
		ScChainId: ibcTx.ScChainId,
		DcChainId: ibcTx.DcChainId,
		ScChannel: ibcTx.ScChannel,
		DcChannel: ibcTx.DcChannel,
		Sequence:  ibcTx.Sequence,
		ScTxInfo:  loadTxInfoDto(ibcTx.ScTxInfo),
		DcTxInfo:  loadTxInfoDto(ibcTx.DcTxInfo),
		BaseDenom: ibcTx.BaseDenom,
		Denoms:    Denoms{ScDenom: ibcTx.Denoms.ScDenom, DcDenom: ibcTx.Denoms.DcDenom},
		TxTime:    ibcTx.TxTime,
		EndTime:   endTime,
	}
}

func (dto IbcTxDetailDto) LoadDto(ibcTx *entity.ExIbcTx) IbcTxDetailDto {
	return IbcTxDetailDto{
		ScAddr:    ibcTx.ScAddr,
		DcAddr:    ibcTx.DcAddr,
		Status:    int(ibcTx.Status),
		ScChainId: ibcTx.ScChainId,
		ScChannel: ibcTx.ScChannel,
		ScPort:    ibcTx.ScPort,
		DcChainId: ibcTx.DcChainId,
		DcChannel: ibcTx.DcChannel,
		DcPort:    ibcTx.DcPort,
		Sequence:  ibcTx.Sequence,
		ScTxInfo:  loadTxInfoDto(ibcTx.ScTxInfo),
		DcTxInfo:  loadTxInfoDto(ibcTx.DcTxInfo),
		BaseDenom: ibcTx.BaseDenom,
		Denoms:    Denoms{ScDenom: ibcTx.Denoms.ScDenom, DcDenom: ibcTx.Denoms.DcDenom},
		TxTime:    ibcTx.TxTime,
	}
}
