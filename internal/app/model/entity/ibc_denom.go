package entity

type IBCDenom struct {
	Symbol           string `bson:"symbol"`
	ChainId          string `bson:"chain_id"`
	Denom            string `bson:"denom"`
	PrevDenom        string `bson:"prev_denom"`
	PrevChainId      string `bson:"prev_chain_id"`
	BaseDenom        string `bson:"base_denom"`
	BaseDenomChainId string `bson:"base_denom_chain_id"`
	DenomPath        string `bson:"denom_path"`
	IsSourceChain    bool   `bson:"is_source_chain"` // todo remove this
	IsBaseDenom      bool   `bson:"is_base_denom"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCDenom) CollectionName() string {
	return "ibc_denom"
}

type IBCDenomList []*IBCDenom
type IBCDenomMap map[string]*IBCDenom

func (l IBCDenomList) ConvertToMap() IBCDenomMap {
	res := make(map[string]*IBCDenom)
	for _, v := range l {
		res[v.ChainId+v.Denom] = v
	}
	return res
}
