package vo

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"

type BaseInfoResp struct {
	Address         string `json:"address"`
	PubKey          string `json:"pub_key"`
	Chain           string `json:"chain"`
	AccountNumber   string `json:"account_number"`
	PubKeyType      string `json:"pub_key_type"`
	PubKeyAlgorithm string `json:"pub_key_algorithm"`
	AccountSequence string `json:"account_sequence"`
}

type AddrTokenListResp struct {
	Tokens     []AddrToken `json:"tokens"`
	TotalValue string      `json:"total_value"`
}

type AddrToken struct {
	Denom                string           `json:"denom"`
	Chain                string           `json:"chain"`
	BaseDenom            string           `json:"base_denom"`
	BaseDenomChain       string           `json:"base_denom_chain"`
	DenomType            entity.TokenType `json:"denom_type"`
	DenomAmount          string           `json:"denom_amount"`
	DenomAvaliableAmount string           `json:"denom_avaliable_amount"`
	Price                float64          `json:"price"`
	DenomValue           string           `json:"denom_value"`
}
