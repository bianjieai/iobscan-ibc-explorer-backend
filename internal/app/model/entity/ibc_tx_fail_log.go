package entity

const (
	IBCTxFailLogCollName = "ibc_tx_fail_log"
)

type TxFailCode string

const (
	TxFailCodeTimeout                     TxFailCode = "timeout"
	TxFailCodeOther                       TxFailCode = "other"
	TxFailCodeWrongSeq                    TxFailCode = "wrong_sequence"
	TxFailCodeOutOfGas                    TxFailCode = "out_of_gas"
	TxFailCodeInsufficientFunds           TxFailCode = "insufficient_funds"
	TxFailCodeClientNotActive             TxFailCode = "client_is_not_active"
	TxFailCodeParsePacketFowradingInfoErr TxFailCode = "cannot_parse_packet_fowrading_information"
	TxFailCodeChannelNotFound             TxFailCode = "channel_not_found"
	TxFailCodeDecodingBech32Failed        TxFailCode = "decoding_bech32_failed" // decoding bech32 failed
	TxFailCodeInvalidBech32Prefix         TxFailCode = "invalid_bech32_prefix"  // invalid Bech32 prefix
	TxFailCodeUnauthorized                TxFailCode = "unauthorized"
	TxFailCodeInvalidCoins                TxFailCode = "invalid_coins"
	TxFailCodeErrorHandlingPacket         TxFailCode = "error_handling_packet"                          // error handling packet
	TxFailCodeIncorrectAccountSequence    TxFailCode = "incorrect_account_sequence"                     // incorrect account sequence
	TxFailCodeDenominationTraceNotFound   TxFailCode = "denomination_trace_not_found"                   // denomination trace not found
	TxFailCodeParsedChannelNotMatchPacket TxFailCode = "parsed_channel_from_denom_doesn't_match_packet" // Parsed channel from denom (channel-42) doesn't match packet
)

type IBCTxFailLog struct {
	Chain            string     `bson:"chain"`
	Log              string     `bson:"log"`
	Code             TxFailCode `bson:"code"`
	TxsNumber        int64      `bson:"txs_number"`
	SegmentStartTime int64      `bson:"segment_start_time"`
	SegmentEndTime   int64      `bson:"segment_end_time"`
	CreateAt         int64      `bson:"create_at"`
	UpdateAt         int64      `bson:"update_at"`
}

func (i *IBCTxFailLog) CollectionName() string {
	return IBCTxFailLogCollName
}
