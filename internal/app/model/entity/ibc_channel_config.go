package entity

type IBCChannelConfig struct {
	ChainA        string `bson:"chain_a"`
	ChainB        string `bson:"chain_b"`
	ChannelA      string `bson:"channel_a"`
	ChannelB      string `bson:"channel_b"`
	ChannelOpenAt int64  `bson:"channel_open_at"`
}

func (i IBCChannelConfig) CollectionName() string {
	return "ibc_channel_config"
}
