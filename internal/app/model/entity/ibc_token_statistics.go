package entity

const (
	IBCTokenStatisticsCollName    = "ibc_token_statistics"
	IBCTokenStatisticsNewCollName = "ibc_token_statistics_new"
)

type IBCTokenStatistics struct {
	BaseDenom        string `bson:"base_denom"`
	BaseDenomChain   string `bson:"base_denom_chain"`
	TransferTxs      int64  `bson:"transfer_txs"`
	SegmentStartTime int64  `bson:"segment_start_time"`
	SegmentEndTime   int64  `bson:"segment_end_time"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCTokenStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCTokenStatisticsNewCollName
	}
	return IBCTokenStatisticsCollName
}
