package entity

import "fmt"

type TokenType string

const (
	TokenTypeAuthed TokenType = "Authed"
	TokenTypeOther  TokenType = "Other"
)

type IBCToken struct {
	BaseDenom      string    `bson:"base_denom"`
	Chain          string    `bson:"chain"`
	Type           TokenType `bson:"type"`
	Price          float64   `bson:"price"`
	Currency       string    `bson:"currency"`
	Supply         string    `bson:"supply"`
	TransferAmount string    `bson:"transfer_amount"`
	TransferTxs    int64     `bson:"transfer_txs"`
	ChainsInvolved int64     `bson:"chains_involved"`
	CreateAt       int64     `bson:"create_at"`
	UpdateAt       int64     `bson:"update_at"`
}

func (i IBCToken) CollectionName() string {
	return "ibc_token"
}

type IBCTokenList []*IBCToken
type IBCTokenMap map[string]*IBCToken

func (l IBCTokenList) ConvertToMap() IBCTokenMap {
	res := make(map[string]*IBCToken)
	for _, v := range l {
		res[fmt.Sprintf("%s%s", v.Chain, v.BaseDenom)] = v
	}
	return res
}
