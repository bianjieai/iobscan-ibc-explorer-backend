package entity

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
)

type TxStatus int64

const (
	TxStatusSuccess TxStatus = 1
	TxStatusFailed  TxStatus = 0
)

type (
	Tx struct {
		Time      int64          `bson:"time"`
		Height    int64          `bson:"height"`
		TxHash    string         `bson:"tx_hash"`
		Type      string         `bson:"type"` // parse from first msg
		Memo      string         `bson:"memo"`
		Status    TxStatus       `bson:"status"`
		Log       string         `bson:"log"`
		Fee       *model.Fee     `bson:"fee"`
		GasUsed   int64          `bson:"gas_used"`
		Types     []string       `bson:"types"`
		EventsNew []EventNew     `bson:"events_new"`
		Signers   []string       `bson:"signers"`
		DocTxMsgs []*model.TxMsg `bson:"msgs"`
		Addrs     []string       `bson:"addrs"`
		TxIndex   uint32         `bson:"tx_index"`
		Ext       interface{}    `bson:"ext"`
	}

	Event struct {
		Type       string   `bson:"type" json:"type"`
		Attributes []KvPair `bson:"attributes" json:"attributes"`
	}

	KvPair struct {
		Key   string `bson:"key" json:"key"`
		Value string `bson:"value" json:"value"`
	}

	EventNew struct {
		MsgIndex uint32  `bson:"msg_index" json:"msg_index"`
		Events   []Event `bson:"events"`
	}
)

func (i Tx) CollectionName(chainId string) string {
	return fmt.Sprintf("sync_%v_tx", chainId)
}
