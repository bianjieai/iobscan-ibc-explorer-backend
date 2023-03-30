package entity

const (
	IBCAddressStatisticsCollName    = "ibc_address_statistics"
	IBCAddressStatisticsNewCollName = "ibc_address_statistics_new"
)

type IBCAddressStatistics struct {
	ChainName        string `bson:"chain_name"`
	ActiveAddressNum int64  `bson:"active_address_num"`
	SegmentStartTime int64  `bson:"segment_start_time"`
	SegmentEndTime   int64  `bson:"segment_end_time"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCAddressStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCAddressStatisticsNewCollName
	}
	return IBCAddressStatisticsCollName
}
