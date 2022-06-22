package entity

type IBCDenomCaculate struct {
	Symbol    string `bson:"symbol"`
	CreateAt  int    `bson:"create_at"`
	UpdateAt  int    `bson:"update_at"`
	BaseDenom string `bson:"base_denom"`
	Denom     string `bson:"denom"`
	DenomPath string `bson:"denom_path"`
	ChainId   string `bson:"chain_id"`
	ScChainId string `bson:"sc_chain_id"`
}

func (i IBCDenomCaculate) CollectionName() string {
	return "ibc_denom_caculate"
}
