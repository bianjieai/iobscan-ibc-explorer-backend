package entity

type ChannelStatus int

const (
	ChannelStatusOpened = 1
	ChannelStatusClosed = 2
)

type IBCChannel struct {
	ChannelId        string        `bson:"channel_id"`
	ChainA           string        `bson:"chain_a"`
	ChainB           string        `bson:"chain_b"`
	ChannelA         string        `bson:"channel_a"`
	ChannelB         string        `bson:"channel_b"`
	Status           ChannelStatus `bson:"status"`
	OperatingPeriod  int64         `bson:"operating_period"`
	ChannelUpdateAt  int64         `bson:"channel_update_at"`
	Relayers         int           `bson:"relayers"`
	TransferTxs      int64         `bson:"transfer_txs"`
	TransferTxsValue float64       `bson:"transfer_txs_value"`
	CreateAt         int64         `bson:"create_at"`
	UpdateAt         int64         `bson:"update_at"`
}

func (i IBCChannel) CollectionName() string {
	return "ibc_channel"
}
