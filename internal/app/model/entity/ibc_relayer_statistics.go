package entity

type IBCRelayerStatistics struct {
	RelayerId          string  `bson:"relayer_id"`
	ChainId            string  `bson:"chain_id"`
	TransferBaseDenom  string  `bson:"transfer_base_denom"`
	TransferAmount     string  `bson:"transfer_amount"`
	TransferTotalValue float64 `bson:"transfer_total_value"`
	CreateAt           int64   `bson:"create_at"`
	UpdateAt           int64   `bson:"update_at"`
}

func (i IBCRelayerStatistics) CollectionName() string {
	return "ibc_relayer_statistics"
}
