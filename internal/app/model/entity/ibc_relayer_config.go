package entity

type IBCRelayerConfig struct {
	RelayerId   string `bson:"relayer_id"`
	RelayerName string `bson:"relayer_name"`
	Icon        string `bson:"icon"`
}

func (i IBCRelayerConfig) CollectionName() string {
	return "ibc_relayer_config"
}
