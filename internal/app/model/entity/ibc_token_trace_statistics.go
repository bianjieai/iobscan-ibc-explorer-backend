package entity

type IBCTokenTraceStatistics struct {
	Denom            string `bson:"denom"`
	ChainId          string `bson:"chain_id"`
	ReceiveTxs       int64  `bson:"receive_txs"`
	SegmentStartTime int64  `bson:"segment_start_time"`
	SegmentEndTime   int64  `bson:"segment_end_time"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCTokenTraceStatistics) CollectionName() string {
	return "ibc_token_trace_statistics"
}
