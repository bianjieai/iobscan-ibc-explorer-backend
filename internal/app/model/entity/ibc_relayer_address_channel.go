package entity

const IBCRelayerAddressChannelCollName = "ibc_relayer_address_channel"

type IBCRelayerAddressChannel struct {
	RelayerAddress      string `json:"relayer_address"`
	Channel             string `json:"channel"`
	Chain               string `json:"chain"`
	CounterPartyChannel string `json:"counter_party_channel"`
	CreateAt            int64  `bson:"create_at"`
	UpdateAt            int64  `bson:"update_at"`
}

func (i IBCRelayerAddressChannel) CollectionName() string {
	return IBCRelayerAddressChannelCollName
}
