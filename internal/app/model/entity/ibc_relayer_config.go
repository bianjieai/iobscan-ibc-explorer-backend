package entity

// todo remove this
type IBCRelayerConfig struct {
	RelayerPairId string `bson:"relayer_pair_id"`
	ChainA        string `bson:"chain_a"`
	ChainB        string `bson:"chain_b"`
	ChannelA      string `bson:"channel_a"`
	ChannelB      string `bson:"channel_b"`
	ChainAAddress string `bson:"chain_a_address"`
	ChainBAddress string `bson:"chain_b_address"`
	RelayerName   string `bson:"relayer_name"`
	Icon          string `bson:"icon"`
}

func (i IBCRelayerConfig) CollectionName() string {
	return "ibc_relayer_config"
}
