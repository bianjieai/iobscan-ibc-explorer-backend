package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Fee struct {
		Amount []*Coin `bson:"amount" json:"amount"`
		Gas    int64   `bson:"gas" json:"gas"`
	}

	Coin struct {
		Denom  string `bson:"denom" json:"denom"`
		Amount string `bson:"amount" json:"amount"`
	}

	TxMsg struct {
		Type string `bson:"type" json:"type"`
		Msg  bson.M `bson:"msg" json:"msg"`
	}

	CommonMsg struct {
		ClientId string `bson:"client_id"`
		PacketId string `bson:"packet_id" json:"packet_id"`
		Signer   string `bson:"signer" json:"signer"`
	}

	TransferTxMsg struct {
		PacketId         string        `bson:"packet_id" json:"packet_id"`
		SourcePort       string        `bson:"source_port" json:"source_port"`
		SourceChannel    string        `bson:"source_channel" json:"source_channel"`
		Token            *Coin         `bson:"token" json:"token"`
		Sender           string        `bson:"sender" json:"sender"`
		Receiver         string        `bson:"receiver" json:"receiver"`
		TimeoutHeight    TimeoutHeight `bson:"timeout_height" json:"timeout_height"`
		TimeoutTimestamp int64         `bson:"timeout_timestamp" json:"timeout_timestamp"`
	}

	TimeoutHeight struct {
		RevisionNumber int64 `json:"revision_number" bson:"destination_channel"`
		RevisionHeight int64 `json:"revision_height" bson:"revision_height"`
	}

	ProofHeight struct {
		RevisionNumber int64 `json:"revision_number" bson:"revision_number"`
		RevisionHeight int64 `json:"revision_height" bson:"revision_height"`
	}

	Packet struct {
		Sequence           int64  `json:"sequence" bson:"sequence"`
		SourcePort         string `json:"source_port" bson:"source_port"`
		SourceChannel      string `json:"source_channel" bson:"source_channel"`
		DestinationPort    string `json:"destination_port" bson:"destination_port"`
		DestinationChannel string `json:"destination_channel" bson:"destination_channel"`
		Data               struct {
			Denom    string      `json:"denom" bson:"denom"`
			Amount   interface{} `json:"amount" bson:"amount"`
			Sender   string      `json:"sender" bson:"sender"`
			Receiver string      `json:"receiver" bson:"receiver"`
		} `json:"data" bson:"data"`
		TimeoutHeight    TimeoutHeight `json:"timeout_height" bson:"timeout_height"`
		TimeoutTimestamp int64         `json:"timeout_timestamp" bson:"timeout_timestamp"`
	}

	RecvPacketMsg struct {
		PacketId        string      `json:"packet_id" bson:"packet_id"`
		Packet          Packet      `json:"packet" bson:"packet"`
		ProofCommitment string      `json:"proof_commitment" bson:"proof_commitment"`
		ProofHeight     ProofHeight `json:"proof_height" bson:"proof_height"`
		Signer          string      `json:"signer" bson:"signer"`
	}

	TimeoutPacketMsg struct {
		PacketId         string      `json:"packet_id" bson:"packet_id"`
		Packet           Packet      `json:"packet" bson:"packet"`
		ProofUnreceived  string      `json:"proof_unreceived" bson:"proof_unreceived"`
		ProofHeight      ProofHeight `json:"proof_height" bson:"proof_height"`
		NextSequenceRecv int64       `json:"next_sequence_recv" bson:"next_sequence_recv"`
		Signer           string      `json:"signer" bson:"signer"`
	}

	AckPacketMsg struct {
		PacketId        string      `json:"packet_id" bson:"packet_id"`
		Packet          Packet      `json:"packet" bson:"packet"`
		Acknowledgement string      `json:"acknowledgement" bson:"acknowledgement"`
		ProofAcked      string      `json:"proof_acked" bson:"proof_acked"`
		ProofHeight     ProofHeight `json:"proof_height" bson:"proof_height"`
		Signer          string      `json:"signer" bson:"signer"`
	}

	PacketDataMsg struct {
		PacketId string `json:"packet_id" bson:"packet_id"`
		Packet   Packet `json:"packet" bson:"packet"`
	}
)

func (m TxMsg) CommonMsg() CommonMsg {
	var msg CommonMsg
	bz, _ := json.Marshal(m.Msg)
	_ = json.Unmarshal(bz, &msg)
	return msg
}

func (m TxMsg) PacketDataMsg() PacketDataMsg {
	var msg PacketDataMsg
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

func (m TxMsg) RecvPacketMsg() RecvPacketMsg {
	var msg RecvPacketMsg
	bz, _ := json.Marshal(m.Msg)
	_ = json.Unmarshal(bz, &msg)

	return msg
}

func (m TxMsg) TimeoutPacketMsg() TimeoutPacketMsg {
	var msg TimeoutPacketMsg
	bz, _ := json.Marshal(m.Msg)
	_ = json.Unmarshal(bz, &msg)

	return msg
}

func (m TxMsg) AckPacketMsg() AckPacketMsg {
	var msg AckPacketMsg
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
