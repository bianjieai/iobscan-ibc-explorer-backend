package entity

type IBCToken struct {
	BaseDenom      string  `bson:"base_denom"`
	ChainId        string  `bson:"chain_id"`
	Type           string  `bson:"type"`
	Price          float64 `bson:"price"`
	Currency       string  `bson:"currency"`
	Supply         string  `bson:"supply"`
	TransferAmount string  `bson:"transfer_amount"`
	TransferTxs    int64   `bson:"transfer_txs"`
	ChainsInvolved int64   `bson:"chains_involved"`
}

func (i IBCToken) CollectionName() string {
	return "ibc_token"
}
