package entity

const (
	IBCRelayerFeeStatisticsCollName    = "ibc_relayer_fee_statistics"
	IBCRelayerFeeStatisticsNewCollName = "ibc_relayer_fee_statistics_new"
)

type IBCRelayerFeeStatistics struct {
	StatisticChain   string   `bson:"statistics_chain"`
	RelayerAddress   string   `bson:"relayer_address"`
	ChainAddressComb string   `bson:"chain_address_comb"`
	TxStatus         TxStatus `bson:"tx_status"`
	TxType           TxType   `bson:"tx_type"`
	FeeDenom         string   `bson:"fee_denom"`
	FeeAmount        float64  `bson:"fee_amount"`
	RelayedTxs       int64    `bson:"relayed_txs"`
	SegmentStartTime int64    `bson:"segment_start_time"`
	SegmentEndTime   int64    `bson:"segment_end_time"`
	CreateAt         int64    `bson:"create_at"`
	UpdateAt         int64    `bson:"update_at"`
}

func (i IBCRelayerFeeStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCRelayerFeeStatisticsNewCollName
	}
	return IBCRelayerFeeStatisticsCollName
}
