package dto

type CountBaseDenomTransferAmountDTO struct {
	BaseDenom string `bson:"base_denom"`
	ScChainId string `bson:"sc_chain_id"`
	DcChainId string `bson:"dc_chain_id"`
	Count     int64  `bson:"count"`
}

type GetDenomGroupByBaseDenomDTO struct {
	BaseDenom string   `bson:"_id"`
	Denom     []string `bson:"denom"`
}

type GetRelayerInfoDTO struct {
	DcChainAddress string `bson:"dc_chain_address"`
	ScChainId      string `bson:"sc_chain_id"`
	ScChannel      string `bson:"sc_channel"`
	DcChainId      string `bson:"dc_chain_id"`
	DcChannel      string `bson:"dc_channel"`
}

type GetRelayerScChainAddreeDTO struct {
	ScChainAddress string `bson:"sc_chain_address"`
}
