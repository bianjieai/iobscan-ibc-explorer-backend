package entity

type IBCBaseDenom struct {
	ChainId     string `bson:"chain_id"`
	Denom       string `bson:"denom"`
	Symbol      string `bson:"symbol"`
	Scale       int    `bson:"scale"`
	Icon        string `bson:"icon"`
	IsMainToken bool   `bson:"is_main_token"`
	//CreateAt    int64  `bson:"create_at"`
	//UpdateAt    int64  `bson:"update_at"`
	CoinId              string `bson:"coin_id"`
	IbcInfoHashCaculate string `bson:"ibc_info_hash_caculate"`
}

func (i IBCBaseDenom) CollectionName() string {
	return "ibc_base_denom"
}

type IBCBaseDenomList []*IBCBaseDenom
type IBCBaseDenomMap map[string]*IBCBaseDenom

func (l IBCBaseDenomList) ConvertToMap() IBCBaseDenomMap {
	res := make(map[string]*IBCBaseDenom)
	for _, v := range l {
		res[v.Denom] = v
	}
	return res
}
