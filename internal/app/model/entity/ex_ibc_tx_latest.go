package entity

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
)

type IbcTxStatus int

const (
	IbcTxStatusSuccess    IbcTxStatus = 1
	IbcTxStatusFailed     IbcTxStatus = 2
	IbcTxStatusProcessing IbcTxStatus = 3
	IbcTxStatusRefunded   IbcTxStatus = 4
	IbcTxStatusSetting    IbcTxStatus = 5
)

var IbcTxUsefulStatus = []IbcTxStatus{IbcTxStatusSuccess, IbcTxStatusFailed, IbcTxStatusProcessing, IbcTxStatusRefunded}

const (
	CollectionNameExIbcTx       = "ex_ibc_tx"
	CollectionNameExIbcTxLatest = "ex_ibc_tx_latest"
)

type (
	ExIbcTx struct {
		RecordId       string      `bson:"record_id"`
		TxTime         int64       `bson:"tx_time"`
		ScAddr         string      `bson:"sc_addr"`
		DcAddr         string      `bson:"dc_addr"`
		ScPort         string      `bson:"sc_port"`
		ScChannel      string      `bson:"sc_channel"`
		ScConnectionId string      `bson:"sc_connection_id"`
		ScClientId     string      `bson:"sc_client_id"`
		ScChainId      string      `bson:"sc_chain_id"`
		DcPort         string      `bson:"dc_port"`
		DcChannel      string      `bson:"dc_channel"`
		DcConnectionId string      `bson:"dc_connection_id"`
		DcClientId     string      `bson:"dc_client_id"`
		DcChainId      string      `bson:"dc_chain_id"`
		Sequence       string      `bson:"sequence"`
		Status         IbcTxStatus `bson:"status"`
		ScTxInfo       *TxInfo     `bson:"sc_tx_info"`
		DcTxInfo       *TxInfo     `bson:"dc_tx_info"`
		RefundedTxInfo *TxInfo     `bson:"refunded_tx_info"`
		//Log              *Log        `bson:"log"`
		Denoms           *Denoms `bson:"denoms"`
		BaseDenom        string  `bson:"base_denom"`
		BaseDenomChainId string  `bson:"base_denom_chain_id"`
		ProcessInfo      string  `bson:"process_info"`
		RetryTimes       int64   `bson:"retry_times"`
		NextTryTime      int64   `bson:"next_try_time"`
		CreateAt         int64   `bson:"create_at"`
		UpdateAt         int64   `bson:"update_at"`
	}
	Log struct {
		ScLog string `bson:"sc_log"`
	}
	Denoms struct {
		ScDenom string `bson:"sc_denom"`
		DcDenom string `bson:"dc_denom"`
	}
	TxInfo struct {
		Hash      string       `bson:"hash"`
		Status    TxStatus     `bson:"status"`
		Time      int64        `bson:"time"`
		Height    int64        `bson:"height"`
		Fee       *model.Fee   `bson:"fee"`
		MsgAmount *model.Coin  `bson:"msg_amount"`
		Msg       *model.TxMsg `bson:"msg"`
		Memo      string       `bson:"memo"`
		Signers   []string     `bson:"signers"`
		Log       string       `bson:"log"`
	}
)

func (i ExIbcTx) CollectionName(historyData bool) string {
	if historyData {
		return CollectionNameExIbcTx
	}
	return CollectionNameExIbcTxLatest
}
