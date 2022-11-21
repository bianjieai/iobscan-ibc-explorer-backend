package entity

const IBCRelayerAddressChannelCollName = "ibc_relayer_address_channel"

type IBCRelayerAddressChannel struct {
	RelayerAddress      string `bson:"relayer_address"`
	Channel             string `bson:"channel"`
	Chain               string `bson:"chain"`
	CounterPartyChannel string `bson:"counter_party_channel"`
	CreateAt            int64  `bson:"create_at"`
	UpdateAt            int64  `bson:"update_at"`
}

func (i IBCRelayerAddressChannel) CollectionName() string {
	return IBCRelayerAddressChannelCollName
}
