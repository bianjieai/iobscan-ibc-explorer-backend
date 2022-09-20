package entity

import "fmt"

type IBCDenom struct {
	Symbol           string `bson:"symbol"`
	ChainId          string `bson:"chain_id"`
	Denom            string `bson:"denom"`
	PrevDenom        string `bson:"prev_denom"`
	PrevChainId      string `bson:"prev_chain_id"`
	BaseDenom        string `bson:"base_denom"`
	BaseDenomChainId string `bson:"base_denom_chain_id"`
	DenomPath        string `bson:"denom_path"`
	RootDenom        string `bson:"root_denom"`
	IsBaseDenom      bool   `bson:"is_base_denom"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCDenom) CollectionName(isNew bool) string {
	if isNew {
		return "ibc_denom_new"
	}
	return "ibc_denom"
}

type IBCDenomList []*IBCDenom
type IBCDenomMap map[string]*IBCDenom

func (l IBCDenomList) ConvertToMap() IBCDenomMap {
	res := make(map[string]*IBCDenom)
	for _, v := range l {
		res[fmt.Sprintf("%s%s", v.ChainId, v.Denom)] = v
	}
	return res
}
