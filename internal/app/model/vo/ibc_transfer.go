package vo

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type (
	TranaferTxsReq struct {
		Page
		UseCount       bool   `json:"use_count" form:"use_count"`
		DateRange      string `json:"date_range" form:"date_range"`
		Status         string `json:"status" form:"status"`
		Chain          string `json:"chain" form:"chain"`
		Denom          string `json:"denom" form:"denom"`
		BaseDenom      string `json:"base_denom" form:"base_denom"`
		BaseDenomChain string `json:"base_denom_chain" form:"base_denom_chain"`
	}
	TranaferTxsResp struct {
		Items     []IbcTxDto `json:"items"`
		PageInfo  PageInfo   `json:"page_info"`
		TimeStamp int64      `json:"time_stamp"`
	}

	TransferTxsCountResp struct {
		TxsCount int64  `json:"txs_count"`
		TxsValue string `json:"txs_value"`
	}

	TxsCountChanDTO struct {
		Count int64
		Err   error
	}

	TxsValueChanDTO struct {
		Value string
		Err   error
	}

	IbcTxDto struct {
		RecordId       string    `json:"record_id"`
		ScAddr         string    `json:"sc_addr"`
		DcAddr         string    `json:"dc_addr"`
		Status         int       `json:"status"`
		ScChain        string    `json:"sc_chain"`
		DcChain        string    `json:"dc_chain"`
		ScChannel      string    `json:"sc_channel"`
		DcChannel      string    `json:"dc_channel"`
		Sequence       string    `json:"sequence"`
		ScTxInfo       TxInfoDto `json:"sc_tx_info"`
		DcTxInfo       TxInfoDto `json:"dc_tx_info"`
		BaseDenom      string    `json:"base_denom"`
		BaseDenomChain string    `json:"base_denom_chain"`
		Denoms         Denoms    `json:"denoms"`
		TxTime         int64     `json:"tx_time"`
		EndTime        int64     `json:"end_time"`
	}

	TranaferTxDetailNewResp struct {
		Items       []IbcTxDto   `json:"items,omitempty"`
		IsList      bool         `json:"is_list"`
		ScInfo      *ChainInfo   `json:"sc_info"`
		DcInfo      *ChainInfo   `json:"dc_info"`
		TokenInfo   *TokenInfo   `json:"token_info"`
		RelayerInfo *RelayerInfo `json:"relayer_info"`
		IbcTxInfo   *IbcTxInfo   `json:"ibc_tx_info"`
		Status      int          `json:"status"`
		Sequence    string       `json:"sequence"`
		ErrorLog    string       `json:"error_log"`
		TimeStamp   int64        `json:"time_stamp"`
	}

	TraceSourceReq struct {
		Chain   string `json:"chain" form:"chain"`
		MsgType string `json:"msg_type" form:"msg_type"`
	}

	TraceSourceResp struct {
		Msg    interface{} `json:"msg"`
		Events interface{} `json:"events"`
	}

	SearchConditionResp struct {
		TxTimeMin int64 `json:"tx_time_min"`
	}

	Denoms struct {
		ScDenom string `json:"sc_denom"`
		DcDenom string `json:"dc_denom"`
	}
	ChainInfo struct {
		Address      string `json:"address"`
		Chain        string `json:"chain"`
		ChannelId    string `json:"channel_id"`
		PortId       string `json:"port_id"`
		ConnectionId string `json:"connection_id"`
		ClientId     string `json:"client_id"`
	}
	TokenInfo struct {
		BaseDenom      string      `json:"base_denom"`
		BaseDenomChain string      `json:"base_denom_chain"`
		SendToken      DetailToken `json:"send_token"`
		RecvToken      DetailToken `json:"recv_token"`
		Amount         string      `json:"amount"`
	}
	DetailToken struct {
		Denom     string `json:"denom"`
		DenomPath string `json:"denom_path"`
	}
	RelayerInfo struct {
		ScRelayer RelayerCfg `json:"sc_relayer,omitempty"`
		DcRelayer RelayerCfg `json:"dc_relayer,omitempty"`
	}
	RelayerCfg struct {
		RelayerName string `json:"relayer_name"`
		Icon        string `json:"icon"`
		RelayerAddr string `json:"relayer_addr"`
	}
	IbcTxInfo struct {
		ScTxInfo         *TxDetailDto `json:"sc_tx_info"`
		DcTxInfo         *TxDetailDto `json:"dc_tx_info"`
		AckTimeoutTxInfo *TxDetailDto `json:"ack_timeout_tx_info"`
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

	TxDetailDto struct {
		TxHash           string     `json:"tx_hash"`
		Status           int        `json:"status"`
		Time             int64      `json:"time"`
		Height           int64      `json:"height"`
		Type             string     `json:"type"`
		Memo             string     `json:"memo"`
		Fee              *model.Fee `json:"fee"`
		Signers          []string   `json:"signers"`
		TimeoutHeight    string     `json:"timeout_height"`
		TimeoutTimestamp int64      `json:"timeout_timestamp"`
		Ack              string     `json:"ack"`
		ProofHeight      string     `json:"proof_height"`
		NextSequenceRecv int64      `json:"next_sequence_recv"`
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
func loadChainInfo(tx *entity.ExIbcTx) (*ChainInfo, *ChainInfo) {
	return &ChainInfo{
			Address:      tx.ScAddr,
			Chain:        tx.ScChain,
			ChannelId:    tx.ScChannel,
			PortId:       tx.ScPort,
			ConnectionId: tx.ScConnectionId,
			ClientId:     tx.ScClientId,
		}, &ChainInfo{
			Address:      tx.DcAddr,
			Chain:        tx.DcChain,
			ChannelId:    tx.DcChannel,
			PortId:       tx.DcPort,
			ConnectionId: tx.DcConnectionId,
			ClientId:     tx.DcClientId,
		}
}
func loadTxDetailDto(info *entity.TxInfo) *TxDetailDto {
	if info == nil {
		return &TxDetailDto{}
	}
	dto := &TxDetailDto{
		TxHash:  info.Hash,
		Status:  int(info.Status),
		Time:    info.Time,
		Height:  info.Height,
		Fee:     info.Fee,
		Memo:    info.Memo,
		Signers: info.Signers,
	}

	proofHeightString := func(pfHeight model.ProofHeight) string {
		if pfHeight.RevisionHeight == 0 {
			return ""
		}
		return fmt.Sprintf("%v-%v", pfHeight.RevisionNumber, pfHeight.RevisionHeight)
	}
	timeHeightString := func(timeoutHeight model.TimeoutHeight) string {
		if timeoutHeight.RevisionHeight == 0 {
			return ""
		}
		return fmt.Sprintf("%v-%v", timeoutHeight.RevisionNumber, timeoutHeight.RevisionHeight)
	}
	if info.Msg != nil {
		dto.Type = info.Msg.Type
		switch info.Msg.Type {
		case constant.MsgTypeRecvPacket:
			//dto.Ack = info.Ack
			dto.ProofHeight = proofHeightString(info.Msg.RecvPacketMsg().ProofHeight)
		case constant.MsgTypeAcknowledgement:
			dto.Ack = info.Msg.AckPacketMsg().Acknowledgement
			dto.ProofHeight = proofHeightString(info.Msg.AckPacketMsg().ProofHeight)
		case constant.MsgTypeTimeoutPacket:
			dto.NextSequenceRecv = info.Msg.TimeoutPacketMsg().NextSequenceRecv
			dto.ProofHeight = proofHeightString(info.Msg.TimeoutPacketMsg().ProofHeight)
		case constant.MsgTypeTransfer:
			dto.TimeoutTimestamp = info.Msg.TransferMsg().TimeoutTimestamp
			dto.TimeoutHeight = timeHeightString(info.Msg.TransferMsg().TimeoutHeight)
		}
	}

	return dto
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
			if ibcTx.AckTimeoutTxInfo != nil {
				endTime = ibcTx.AckTimeoutTxInfo.Time
			}
		}
	case entity.IbcTxStatusRefunded:
		if ibcTx.AckTimeoutTxInfo != nil {
			endTime = ibcTx.AckTimeoutTxInfo.Time
		}
	}
	return IbcTxDto{
		RecordId:       ibcTx.RecordId,
		ScAddr:         ibcTx.ScAddr,
		DcAddr:         ibcTx.DcAddr,
		Status:         int(ibcTx.Status),
		ScChain:        ibcTx.ScChain,
		DcChain:        ibcTx.DcChain,
		ScChannel:      ibcTx.ScChannel,
		DcChannel:      ibcTx.DcChannel,
		Sequence:       ibcTx.Sequence,
		ScTxInfo:       loadTxInfoDto(ibcTx.ScTxInfo),
		DcTxInfo:       loadTxInfoDto(ibcTx.DcTxInfo),
		BaseDenom:      ibcTx.BaseDenom,
		BaseDenomChain: ibcTx.BaseDenomChain,
		Denoms:         Denoms{ScDenom: ibcTx.Denoms.ScDenom, DcDenom: ibcTx.Denoms.DcDenom},
		TxTime:         ibcTx.TxTime,
		EndTime:        endTime,
	}
}

func LoadTranaferTxDetail(ibcTx *entity.ExIbcTx) TranaferTxDetailNewResp {
	scChainInfo, dcChainInfo := loadChainInfo(ibcTx)
	var errLog string
	switch ibcTx.Status {
	case entity.IbcTxStatusFailed:
		errLog = ibcTx.ScTxInfo.Log
	case entity.IbcTxStatusRefunded:
		if ibcTx.DcTxInfo != nil {
			if ibcTx.DcTxInfo.Status == entity.TxStatusSuccess {
				if ibcTx.AckTimeoutTxInfo != nil {
					errLog = ibcTx.AckTimeoutTxInfo.Msg.AckPacketMsg().Acknowledgement
				}
			} else {
				errLog = ibcTx.DcTxInfo.Log
			}
		}
	}
	scTxInfo := loadTxDetailDto(ibcTx.ScTxInfo)
	ibcTxInfo := &IbcTxInfo{
		ScTxInfo: scTxInfo,
	}
	if ibcTx.DcTxInfo != nil {
		ibcTxInfo.DcTxInfo = loadTxDetailDto(ibcTx.DcTxInfo)
		if ibcTx.AckTimeoutTxInfo != nil && ibcTx.AckTimeoutTxInfo.Msg != nil {
			ibcTxInfo.DcTxInfo.Ack = ibcTx.AckTimeoutTxInfo.Msg.AckPacketMsg().Acknowledgement
		}
	}
	if ibcTx.AckTimeoutTxInfo != nil {
		ibcTxInfo.AckTimeoutTxInfo = loadTxDetailDto(ibcTx.AckTimeoutTxInfo)
	}
	return TranaferTxDetailNewResp{
		ErrorLog:  errLog,
		Status:    int(ibcTx.Status),
		Sequence:  ibcTx.Sequence,
		ScInfo:    scChainInfo,
		DcInfo:    dcChainInfo,
		IbcTxInfo: ibcTxInfo,
	}
}
