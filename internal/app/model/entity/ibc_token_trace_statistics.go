package entity

const (
	IBCTokenTraceStatisticsCollName    = "ibc_token_trace_statistics"
	IBCTokenTraceStatisticsNewCollName = "ibc_token_trace_statistics_new"
)

type IBCTokenTraceStatistics struct {
	Denom            string `bson:"denom"`
	Chain            string `bson:"chain"`
	ReceiveTxs       int64  `bson:"receive_txs"`
	SegmentStartTime int64  `bson:"segment_start_time"`
	SegmentEndTime   int64  `bson:"segment_end_time"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCTokenTraceStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCTokenTraceStatisticsNewCollName
	}
	return IBCTokenTraceStatisticsCollName
}
