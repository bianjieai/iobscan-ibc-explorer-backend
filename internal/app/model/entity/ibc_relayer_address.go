package entity

const IBCRelayerAddressCollName = "ibc_relayer_address"

type GatherStatus int

const (
	GatherStatusTODO     GatherStatus = 1
	GatherStatusPubKey   GatherStatus = 2
	GatherStatusRegistry GatherStatus = 3
)

type IBCRelayerAddress struct {
	Address      string       `bson:"address"`
	Chain        string       `bson:"chain"`
	PubKey       string       `bson:"pub_key"`
	GatherStatus GatherStatus `bson:"gather_status"`
	CreateAt     int64        `bson:"create_at"`
	UpdateAt     int64        `bson:"update_at"`
}

func (i IBCRelayerAddress) CollectionName() string {
	return IBCRelayerAddressCollName
}
