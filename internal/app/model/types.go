package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Fee struct {
		Amount []*Coin `bson:"amount"`
		Gas    int64   `bson:"gas"`
	}

	Coin struct {
		Denom  string `bson:"denom" json:"denom"`
		Amount string `bson:"amount" json:"amount"`
	}

	TxMsg struct {
		Type string `bson:"type"`
		Msg  bson.M `bson:"msg"`
	}

	CommonMsg struct {
		ClientId string `bson:"client_id"`
		PacketId string `bson:"packet_id" json:"packet_id"`
		Signer   string `bson:"signer" json:"signer"`
	}

	TransferTxMsg struct {
		PacketId      string `bson:"packet_id" json:"packet_id"`
		SourcePort    string `bson:"source_port" json:"source_port"`
		SourceChannel string `bson:"source_channel" json:"source_channel"`
		Token         *Coin  `bson:"token" json:"token"`
		Sender        string `bson:"sender" json:"sender"`
		Receiver      string `bson:"receiver" json:"receiver"`
		TimeoutHeight struct {
			RevisionNumber int `bson:"revision_number" json:"revision_number"`
			RevisionHeight int `bson:"revision_height" json:"revision_height"`
		} `bson:"timeout_height" json:"timeout_height"`
		TimeoutTimestamp int64 `bson:"timeout_timestamp" json:"timeout_timestamp"`
	}
)

func (m TxMsg) CommonMsg() CommonMsg {
	var msg CommonMsg
	bz, _ := json.Marshal(m.Msg)
	_ = json.Unmarshal(bz, &msg)
	return msg
}

func (m TxMsg) TransferMsg() TransferTxMsg {
	var msg TransferTxMsg
	bz, _ := json.Marshal(m.Msg)
	_ = json.Unmarshal(bz, &msg)

	return msg
}

type TransferTxPacketData struct {
	Amount   string `json:"amount"`
	Denom    string `json:"denom"`
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
}
