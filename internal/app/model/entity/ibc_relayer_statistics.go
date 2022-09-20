package entity

const (
	IBCRelayerStatisticsCollName    = "ibc_relayer_statistics"
	IBCRelayerStatisticsNewCollName = "ibc_relayer_statistics_new"
)

type IBCRelayerStatistics struct {
	StatisticId       string `bson:"statistic_id"` // scChain|scChannel|dcChain|dcChannel
	Address           string `bson:"address"`
	TransferBaseDenom string `bson:"transfer_base_denom"`
	BaseDenomChainId  string `bson:"base_denom_chain_id"`
	TransferAmount    string `bson:"transfer_amount"`
	SuccessTotalTxs   int64  `bson:"success_total_txs"`
	TotalTxs          int64  `bson:"total_txs"`
	SegmentStartTime  int64  `bson:"segment_start_time"`
	SegmentEndTime    int64  `bson:"segment_end_time"`
	CreateAt          int64  `bson:"create_at"`
	UpdateAt          int64  `bson:"update_at"`
}

func (i IBCRelayerStatistics) CollectionName(isNew bool) string {
	if isNew {
		return IBCRelayerStatisticsNewCollName
	}
	return IBCRelayerStatisticsCollName
}
