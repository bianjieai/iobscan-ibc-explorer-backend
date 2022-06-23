package entity

type IBCChannelStatistics struct {
	ChannelId          string `bson:"channel_id"`
	TransferBaseDenom  string `bson:"transfer_base_denom"`
	TransferAmount     string `bson:"transfer_amount"`
	TransferTotalValue string `bson:"transfer_total_value"`
	CreateAt           int64  `bson:"create_at"`
	UpdateAt           int64  `bson:"update_at"`
}

func (i IBCChannelStatistics) CollectionName() string {
	return "ibc_channel_statistics"
}
