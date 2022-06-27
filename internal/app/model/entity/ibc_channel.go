package entity

type ChannelStatus int

const (
	ChannelStatusOpened = 1
	ChannelStatusClosed = 2
)

type IBCChannel struct {
	ChannelId            string        `bson:"channel_id"`
	ChainA               string        `bson:"chain_a"`
	ChainB               string        `bson:"chain_b"`
	ChannelA             string        `bson:"channel_a"`
	ChannelB             string        `bson:"channel_b"`
	Status               ChannelStatus `bson:"status"`
	OperatingPeriod      int64         `bson:"operating_period"`
	LatestSettlementTime int64         `bson:"latest_settlement_time"`
	ChannelUpdateAt      int64         `bson:"channel_update_at"`
	Relayers             int           `bson:"relayers"`
	TransferTxs          int64         `bson:"transfer_txs"`
	TransferTxsValue     string        `bson:"transfer_txs_value"`
	CreateAt             int64         `bson:"create_at"`
	UpdateAt             int64         `bson:"update_at"`
}

func (i IBCChannel) CollectionName() string {
	return "ibc_channel"
}

type IBCChannelList []*IBCChannel

func (l IBCChannelList) ConvertToMap() map[string]*IBCChannel {
	res := make(map[string]*IBCChannel)
	for _, v := range l {
		res[v.ChannelId] = v
	}
	return res
}

func (l IBCChannelList) GetChannelIds() []string {
	res := make([]string, 0, len(l))
	for _, v := range l {
		res = append(res, v.ChannelId)
	}
	return res
}
