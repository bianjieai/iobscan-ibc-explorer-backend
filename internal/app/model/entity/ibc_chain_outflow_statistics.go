package entity

const IBCChainOutflowStatisticsCollName = "ibc_chain_outflow_statistics"

type IBCChainOutflowStatistics struct {
	Chain            string      `bson:"chain"`
	BaseDenom        string      `bson:"base_denom"`
	BaseDenomChain   string      `bson:"base_denom_chain"`
	Status           IbcTxStatus `bson:"status"`
	DenomAmount      string      `bson:"denom_amount"`
	TxsNumber        int64       `bson:"txs_number"`
	SegmentStartTime int64       `bson:"segment_start_time"`
	SegmentEndTime   int64       `bson:"segment_end_time"`
	CreateAt         int64       `bson:"create_at"`
	UpdateAt         int64       `bson:"update_at"`
}

func (i IBCChainOutflowStatistics) CollectionName() string {
	return IBCChainOutflowStatisticsCollName
}
