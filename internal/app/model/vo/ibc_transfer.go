package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type (
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
		RecordId         string    `json:"record_id"`
		ScAddr           string    `json:"sc_addr"`
		DcAddr           string    `json:"dc_addr"`
		Status           int       `json:"status"`
		ScChainId        string    `json:"sc_chain_id"`
		DcChainId        string    `json:"dc_chain_id"`
		ScChannel        string    `json:"sc_channel"`
		DcChannel        string    `json:"dc_channel"`
		Sequence         string    `json:"sequence"`
		ScTxInfo         TxInfoDto `json:"sc_tx_info"`
		DcTxInfo         TxInfoDto `json:"dc_tx_info"`
		BaseDenom        string    `json:"base_denom"`
		BaseDenomChainId string    `json:"base_denom_chain_id"`
		Denoms           Denoms    `json:"denoms"`
		TxTime           int64     `json:"tx_time"`
		EndTime          int64     `json:"end_time"`
	}
	Denoms struct {
		ScDenom string `json:"sc_denom"`
		DcDenom string `json:"dc_denom"`
	}

	IbcTxDetailDto struct {
		RecordId         string    `json:"record_id"`
		ScSigners        []string  `json:"sc_signers"`
		DcSigners        []string  `json:"dc_signers"`
		ScAddr           string    `json:"sc_addr"`
		DcAddr           string    `json:"dc_addr"`
		Status           int       `json:"status"`
		ScChainId        string    `json:"sc_chain_id"`
		ScChannel        string    `json:"sc_channel"`
		ScPort           string    `json:"sc_port"`
		ScConnect        string    `json:"sc_connect"`
		DcChainId        string    `json:"dc_chain_id"`
		DcChannel        string    `json:"dc_channel"`
		DcPort           string    `json:"dc_port"`
		DcConnect        string    `json:"dc_connect"`
		Sequence         string    `json:"sequence"`
		ScTxInfo         TxInfoDto `json:"sc_tx_info"`
		DcTxInfo         TxInfoDto `json:"dc_tx_info"`
		BaseDenom        string    `json:"base_denom"`
		BaseDenomChainId string    `json:"base_denom_chain_id"`
		Denoms           Denoms    `json:"denoms"`
		TxTime           int64     `json:"tx_time"`
		Ack              string    `json:"ack"`
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
		RecordId:         ibcTx.RecordId,
		ScAddr:           ibcTx.ScAddr,
		DcAddr:           ibcTx.DcAddr,
		Status:           int(ibcTx.Status),
		ScChainId:        ibcTx.ScChainId,
		DcChainId:        ibcTx.DcChainId,
		ScChannel:        ibcTx.ScChannel,
		DcChannel:        ibcTx.DcChannel,
		Sequence:         ibcTx.Sequence,
		ScTxInfo:         loadTxInfoDto(ibcTx.ScTxInfo),
		DcTxInfo:         loadTxInfoDto(ibcTx.DcTxInfo),
		BaseDenom:        ibcTx.BaseDenom,
		BaseDenomChainId: ibcTx.BaseDenomChainId,
		Denoms:           Denoms{ScDenom: ibcTx.Denoms.ScDenom, DcDenom: ibcTx.Denoms.DcDenom},
		TxTime:           ibcTx.TxTime,
		EndTime:          endTime,
	}
}

func (dto IbcTxDetailDto) LoadDto(ibcTx *entity.ExIbcTx) IbcTxDetailDto {
	return IbcTxDetailDto{
		RecordId:         ibcTx.RecordId,
		ScAddr:           ibcTx.ScAddr,
		DcAddr:           ibcTx.DcAddr,
		Status:           int(ibcTx.Status),
		ScChainId:        ibcTx.ScChainId,
		ScChannel:        ibcTx.ScChannel,
		ScPort:           ibcTx.ScPort,
		DcChainId:        ibcTx.DcChainId,
		DcChannel:        ibcTx.DcChannel,
		DcPort:           ibcTx.DcPort,
		Sequence:         ibcTx.Sequence,
		ScTxInfo:         loadTxInfoDto(ibcTx.ScTxInfo),
		DcTxInfo:         loadTxInfoDto(ibcTx.DcTxInfo),
		BaseDenom:        ibcTx.BaseDenom,
		BaseDenomChainId: ibcTx.BaseDenomChainId,
		Denoms:           Denoms{ScDenom: ibcTx.Denoms.ScDenom, DcDenom: ibcTx.Denoms.DcDenom},
		TxTime:           ibcTx.TxTime,
	}
}
