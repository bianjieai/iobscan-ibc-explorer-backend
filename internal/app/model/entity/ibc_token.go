package entity

type TokenType string

const (
	TokenTypeAuthed TokenType = "Authed"
	TokenTypeOther  TokenType = "Other"
)

type IBCToken struct {
	BaseDenom      string    `bson:"base_denom"`
	ChainId        string    `bson:"chain_id"`
	Type           TokenType `bson:"type"`
	Price          float64   `bson:"price"`
	Currency       string    `bson:"currency"`
	Supply         string    `bson:"supply"`
	TransferAmount string    `bson:"transfer_amount"`
	TransferTxs    int64     `bson:"transfer_txs"`
	ChainsInvolved int       `bson:"chains_involved"`
}

func (i IBCToken) CollectionName() string {
	return "ibc_token"
}

type IBCTokenList []*IBCToken
type IBCTokenMap map[string]*IBCToken

func (l IBCTokenList) ConvertToMap() IBCTokenMap {
	res := make(map[string]*IBCToken)
	for _, v := range l {
		res[v.ChainId+v.BaseDenom] = v
	}
	return res
}
