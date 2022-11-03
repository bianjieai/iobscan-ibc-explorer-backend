package entity

type TxType string

const (
	TxTypeTransfer      TxType = "transfer"
	TxTypeRecvPacket    TxType = "recv_packet"
	TxTypeTimeoutPacket TxType = "timeout_packet"
	TxTypeAckPacket     TxType = "acknowledge_packet"

	IBCRelayerDenomStatisticsCollName    = "ibc_relayer_denom_statistics"
	IBCRelayerDenomStatisticsNewCollName = "ibc_relayer_denom_statistics_new"
)

type IBCRelayerDenomStatistics struct {
	StatisticChain   string   `bson:"statistics_chain"`
	RelayerAddress   string   `bson:"relayer_address"`
	TxStatus         TxStatus `bson:"tx_status"`
	TxType           TxType   `bson:"tx_type"`
	BaseDenom        string   `bson:"base_denom"`
	BaseDenomChainId string   `bson:"base_denom_chain_id"`
	RelayedAmount    float64  `bson:"relayed_amount"`
	RelayedTxs       int64    `bson:"relayed_txs"`
	SegmentStartTime int64    `bson:"segment_start_time"`
	SegmentEndTime   int64    `bson:"segment_end_time"`
	CreateAt         int64    `bson:"create_at"`
	UpdateAt         int64    `bson:"update_at"`
}

func (i IBCRelayerDenomStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCRelayerDenomStatisticsNewCollName
	}
	return IBCRelayerDenomStatisticsCollName
}
