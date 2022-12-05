package vo

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

type TxReq struct {
	Chain          string `json:"chain" form:"chain"`
	Channel        string `json:"channel" form:"channel"`
	Port           string `json:"port" form:"port"`
	PacketSequence string `json:"packet_sequence" form:"packet_sequence"`
}

type TxResp struct {
	TxDetail      *TxDetail      `json:"tx_detail"`
	Ics20Transfer *Ics20Transfer `json:"ics20_transfer"`
}

type TxFee struct {
	Amount []model.Coin `json:"amount"`
}

type TxDetail struct {
	Chain    string            `json:"chain"`
	TxTime   int64             `json:"tx_time"`
	Height   int64             `json:"height"`
	TxHash   string            `json:"tx_hash"`
	Memo     string            `json:"memo"`
	Status   entity.TxStatus   `json:"status"`
	ErrorLog string            `json:"error_log"`
	Fee      TxFee             `json:"fee"`
	GasUsed  int64             `json:"gas_used"`
	Logs     []entity.EventNew `json:"logs"`
	Msgs     []*model.TxMsg    `json:"msgs"`
	Signers  []string          `json:"signers"`
}

func BuildTxDetail(chain string, tx *entity.Tx) TxDetail {
	var coins []model.Coin
	for _, v := range tx.Fee.Amount {
		coins = append(coins, model.Coin{
			Denom:  v.Denom,
			Amount: v.Amount,
		})
	}

	return TxDetail{
		Chain:    chain,
		TxTime:   tx.Time,
		Height:   tx.Height,
		TxHash:   tx.TxHash,
		Memo:     tx.Memo,
		Status:   tx.Status,
		ErrorLog: tx.Log,
		Fee:      TxFee{Amount: coins},
		GasUsed:  tx.Fee.Gas,
		Logs:     tx.EventsNew,
		Msgs:     tx.DocTxMsgs,
		Signers:  tx.Signers,
	}
}

type SimpleTx struct {
	Chain  string          `json:"chain"`
	TxHash string          `json:"tx_hash"`
	Height int64           `json:"height"`
	TxTime int64           `json:"tx_time"`
	Status entity.TxStatus `json:"status"`
	Msg    *model.TxMsg    `json:"-"`
}

type SimpleTxExt struct {
	SimpleTx
	IsEffective bool `json:"is_effective"`
}

type Ics20Transfer struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	IBCPacket
	Token            *model.Coin   `json:"token"`
	TransferTx       *SimpleTx     `json:"transfer_tx"`
	RecvPacketTxs    []SimpleTxExt `json:"recv_packet_txs"`
	AckPacketTxs     []SimpleTxExt `json:"ack_packet_txs"`
	TimeoutPacketTxs []SimpleTxExt `json:"timeout_packet_txs"`
}

type IBCPacket struct {
	SourcePort         string `json:"source_port"`
	SourceChannel      string `json:"source_channel"`
	SourceChain        string `json:"source_chain"`
	DestinationPort    string `json:"destination_port"`
	DestinationChannel string `json:"destination_channel"`
	DestinationChain   string `json:"destination_chain"`
	PacketSequence     int64  `json:"packet_sequence"`
}

func BuildIBCPacket(scChain, dcChain string, modelPacket model.Packet) IBCPacket {
	return IBCPacket{
		SourcePort:         modelPacket.SourcePort,
		SourceChannel:      modelPacket.SourceChannel,
		SourceChain:        scChain,
		DestinationPort:    modelPacket.DestinationPort,
		DestinationChannel: modelPacket.DestinationChannel,
		DestinationChain:   dcChain,
		PacketSequence:     modelPacket.Sequence,
	}
}

type (
	FailureStatisticsReq struct {
		StartDate string `json:"start_date" form:"start_date"`
		EndDate   string `json:"end_date" form:"end_date"`
	}

	FailureStatisticsResp struct {
		Items            []FailureStatisticsItem `json:"items"`
		StatisticCaliber FailureStatisticCaliber `json:"statistic_caliber"`
	}

	FailureStatisticsItem struct {
		FailureReason         string `json:"failure_reason"`
		FailureTransferNumber int64  `json:"failure_transfer_number"`
	}

	FailureStatisticCaliber struct {
		TxTimeMin int64 `json:"tx_time_min"`
		TxTimeMax int64 `json:"tx_time_max"`
	}
)
