package vo

type FailureStatisticsResp struct {
	Items []struct {
		FailureReason         string `json:"failure_reason"`
		FailureTransferNumber int    `json:"failure_transfer_number"`
	} `json:"items"`
	StatisticCaliber struct {
		TxTimeMin int `json:"tx_time_min"`
		TxTimeMax int `json:"tx_time_max"`
	} `json:"statistic_caliber"`
}

type TxReq struct {
	Chain          string `json:"chain"`
	Channel        string `json:"channel"`
	Port           string `json:"port"`
	PacketSequence string `json:"packet_sequence"`
}

type TxResp struct {
	TxDetail      TxDetail      `json:"tx_detail"`
	Ics20Transfer Ics20Transfer `json:"ics20_transfer"`
}

type TxDetail struct {
	Chain    string `json:"chain"`
	TxTime   int    `json:"tx_time"`
	Height   int    `json:"height"`
	TxHash   string `json:"tx_hash"`
	Memo     string `json:"memo"`
	Status   int    `json:"status"`
	ErrorLog string `json:"error_log"`
	Fee      struct {
		Amount []Token `json:"amount"`
	} `json:"fee"`
	GasUsed int64 `json:"gas_used"`
	Logs    []struct {
		MsgIndex int `json:"msg_index"`
		Events   []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
	} `json:"logs"`
	Msgs []struct {
		Type string `json:"type"`
		Msg  struct {
		} `json:"msg"`
	} `json:"msgs"`
	Signers []string `json:"signers"`
}

type SimpleTxInfo struct {
	Chain  string `json:"chain"`
	TxHash string `json:"tx_hash"`
	Height int    `json:"height"`
	TxTime int    `json:"tx_time"`
	Status int    `json:"status"`
}

type Ics20Transfer struct {
	Sender             string       `json:"sender"`
	Receiver           string       `json:"receiver"`
	SourcePort         string       `json:"source_port"`
	SourceChannel      string       `json:"source_channel"`
	SourceChain        string       `json:"source_chain"`
	DestinationPort    string       `json:"destination_port"`
	DestinationChannel string       `json:"destination_channel"`
	DestinationChain   string       `json:"destination_chain"`
	PacketSequence     int          `json:"packet_sequence"`
	Token              Token        `json:"token"`
	TransferTx         SimpleTxInfo `json:"transfer_tx"`
	RecvPacketTxs      []struct {
		SimpleTxInfo
		IsEffective bool `json:"is_effective"`
	} `json:"recv_packet_txs"`
	AckPacketTxs []struct {
		SimpleTxInfo
		IsEffective bool `json:"is_effective"`
	} `json:"ack_packet_txs"`
	TimeoutPacketTxs []struct {
		SimpleTxInfo
		IsEffective bool `json:"is_effective"`
	} `json:"timeout_packet_txs"`
}

type Token struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
