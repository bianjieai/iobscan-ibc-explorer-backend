package entity

type IBCTokenStatistics struct {
	Denom       string  `bson:"denom"`
	DenomPath   string  `bson:"denom_path"`
	BaseDenom   string  `bson:"base_denom"`
	ChainId     string  `bson:"chain_id"`
	Type        string  `bson:"type"`
	IBCHops     int     `bson:"ibc_hops"`
	DenomAmount string  `bson:"denom_amount"`
	DenomValue  float64 `bson:"denom_value"`
	ReceiveTxs  int64   `bson:"receive_txs"`
}

func (i IBCTokenStatistics) CollectionName() string {
	return "ibc_token_statistics"
}
