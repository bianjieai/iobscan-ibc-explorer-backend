package entity

type IBCDenomCalculate struct {
	Symbol    string `bson:"symbol"`
	BaseDenom string `bson:"base_denom"`
	Denom     string `bson:"denom"`
	DenomPath string `bson:"denom_path"`
	ChainId   string `bson:"chain_id"`
	ScChainId string `bson:"sc_chain_id"`
	CreateAt  int64  `bson:"create_at"`
	UpdateAt  int64  `bson:"update_at"`
}

func (i IBCDenomCalculate) CollectionName() string {
	return "ibc_denom_caculate"
}
