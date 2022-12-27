package entity

import "fmt"

type AuthDenom struct {
	Chain          string `bson:"chain"`
	Denom          string `bson:"denom"`
	Symbol         string `bson:"symbol"`
	Scale          int    `bson:"scale"`
	Icon           string `bson:"icon"`
	IsStakingToken bool   `bson:"is_staking_token"`
	IsStableCoin   bool   `bson:"is_stable_coin"`
	//CreateAt    int64  `bson:"create_at"`
	//UpdateAt    int64  `bson:"update_at"`
	CoinId              string `bson:"coin_id"`
	IbcInfoHashCaculate string `bson:"ibc_info_hash_caculate"`
}

func (i AuthDenom) CollectionName() string {
	return "auth_denom"
}

type AuthDenomList []*AuthDenom
type IBCBaseDenomMap map[string]*AuthDenom

func (l AuthDenomList) ConvertToMap() IBCBaseDenomMap {
	res := make(map[string]*AuthDenom)
	for _, v := range l {
		res[fmt.Sprintf("%s%s", v.Chain, v.Denom)] = v
	}
	return res
}
