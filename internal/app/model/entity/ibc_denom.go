package entity

type IBCDenom struct {
	Symbol        string `bson:"symbol"`
	CreateAt      int64  `bson:"create_at"`
	UpdateAt      int64  `bson:"update_at"`
	TxTime        int64  `bson:"tx_time"`
	RealDenom     bool   `bson:"real_denom"`
	ChainId       string `bson:"chain_id"`
	Denom         string `bson:"denom"`
	BaseDenom     string `bson:"base_denom"`
	DenomPath     string `bson:"denom_path"`
	IsSourceChain bool   `bson:"is_source_chain"`
	IsBaseDenom   bool   `bson:"is_base_denom"`
}

func (i IBCDenom) CollectionName() string {
	return "ibc_denom"
}
