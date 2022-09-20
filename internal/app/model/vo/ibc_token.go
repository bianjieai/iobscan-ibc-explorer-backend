package vo

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"

type TokenListReq struct {
	Page
	BaseDenom string           `json:"base_denom" form:"base_denom"`
	Chain     string           `json:"chain" form:"chain"`
	TokenType entity.TokenType `json:"token_type" form:"token_type"`
	UseCount  bool             `json:"use_count" form:"use_count"`
}

type TokenListResp struct {
	Items    []TokenItem `json:"items"`
	PageInfo PageInfo    `json:"page_info"`
}

type TokenItem struct {
	BaseDenom         string           `json:"base_denom"`
	ChainId           string           `json:"chain_id"`
	TokenType         entity.TokenType `json:"token_type"`
	Supply            string           `json:"supply"`
	Currency          string           `json:"currency"`
	Price             float64          `json:"price"`
	ChainsInvolved    int64            `json:"chains_involved"`
	IBCTransferTxs    int64            `json:"ibc_transfer_txs"`
	IBCTransferAmount string           `json:"ibc_transfer_amount"`
}

type IBCTokenListReq struct {
	Page
	Chain     string                `json:"chain" form:"chain"`
	TokenType entity.TokenTraceType `json:"token_type" form:"token_type"`
	UseCount  bool                  `json:"use_count" form:"use_count"`
}

type IBCTokenListResp struct {
	Items    []IBCTokenItem `json:"items"`
	PageInfo PageInfo       `json:"page_info"`
}

type IBCTokenItem struct {
	Denom      string                `json:"denom"`
	DenomPath  string                `json:"denom_path"`
	ChainId    string                `json:"chain_id"`
	TokenType  entity.TokenTraceType `json:"token_type"`
	IBCHops    int                   `json:"ibc_hops"`
	Amount     string                `json:"amount"`
	ReceiveTxs int64                 `json:"receive_txs"`
}
