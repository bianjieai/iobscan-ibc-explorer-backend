package entity

const (
	IBCChannelStatisticsCollName    = "ibc_channel_statistics"
	IBCChannelStatisticsNewCollName = "ibc_channel_statistics_new"
)

type IBCChannelStatistics struct {
	ChannelId        string      `bson:"channel_id"`
	BaseDenom        string      `bson:"base_denom"`
	BaseDenomChain   string      `bson:"base_denom_chain"`
	TransferTxs      int64       `bson:"transfer_txs"`
	TransferAmount   string      `bson:"transfer_amount"`
	Status           IbcTxStatus `bson:"status"`
	SegmentStartTime int64       `bson:"segment_start_time"`
	SegmentEndTime   int64       `bson:"segment_end_time"`
	CreateAt         int64       `bson:"create_at"`
	UpdateAt         int64       `bson:"update_at"`
}

func (i IBCChannelStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCChannelStatisticsNewCollName
	}
	return IBCChannelStatisticsCollName
}
