package entity

const (
	IBCChainFeeStatisticsCollName    = "ibc_chain_fee_statistics"
	IBCChainFeeStatisticsNewCollName = "ibc_chain_fee_statistics_new"
)

type IBCChainFeeStatistics struct {
	ChainName        string   `bson:"chain_name"`
	TxStatus         TxStatus `bson:"tx_status"`
	TxType           TxType   `bson:"tx_type"`
	FeeDenom         string   `bson:"fee_denom"`
	FeeAmount        float64  `bson:"fee_amount"`
	SegmentStartTime int64    `bson:"segment_start_time"`
	SegmentEndTime   int64    `bson:"segment_end_time"`
	CreateAt         int64    `bson:"create_at"`
	UpdateAt         int64    `bson:"update_at"`
}

func (i IBCChainFeeStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCChainFeeStatisticsNewCollName
	}
	return IBCChainFeeStatisticsCollName
}
