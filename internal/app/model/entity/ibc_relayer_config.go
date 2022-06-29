package entity

type IBCRelayerConfig struct {
	RelayerChannelPair string `bson:"relayer_channel_pair"`
	RelayerName        string `bson:"relayer_name"`
	Icon               string `bson:"icon"`
}

func (i IBCRelayerConfig) CollectionName() string {
	return "ibc_relayer_config"
}
