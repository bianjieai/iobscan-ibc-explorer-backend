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

type CountIBCTokenRecvTxsDTO struct {
	Denom string `bson:"_id"`
	Count int64  `bson:"count"`
}

type AggregateIBCChainDTO struct {
	ChainId    string  `bson:"_id"`
	DenomValue float64 `bson:"denom_value"`
	Count      int64   `bson:"count"`
}
