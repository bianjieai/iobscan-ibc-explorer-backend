package entity

type IBCBaseDenom struct {
	ChainId     string `bson:"chain_id"`
	Denom       string `bson:"denom"`
	Symbol      string `bson:"symbol"`
	Scale       string `bson:"scale"`
	Icon        string `bson:"icon"`
	IsMainToken bool   `bson:"is_main_token"`
	CreateAt    string `bson:"create_at"`
	UpdateAt    string `bson:"update_at"`
	CoinId      string `bson:"coin_id"`
}

func (i IBCBaseDenom) CollectionName() string {
	return "ibc_base_denom"
}
