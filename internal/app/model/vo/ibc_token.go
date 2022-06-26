package vo

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"

type TokenListResp struct {
	Items    []TokenItem `json:"items"`
	PageInfo PageInfo    `json:"page_info"`
}

type TokenItem struct {
	BaseDenom         string  `json:"base_denom"`
	ChainId           string  `json:"chain_id"`
	Supply            string  `json:"supply"`
	Currency          string  `json:"currency"`
	Price             float64 `json:"price"`
	ChainsInvolved    int64   `json:"chains_involved"`
	IBCTransferTxs    int64   `json:"ibc_transfer_txs"`
	IBCTransferAmount string  `json:"ibc_transfer_amount"`
}

type IBCTokenListResp struct {
	Items    []IBCTokenItem `json:"items"`
	PageInfo PageInfo       `json:"page_info"`
}

type IBCTokenItem struct {
	Denom      string                     `json:"denom"`
	DenomPath  string                     `json:"denom_path"`
	ChainId    string                     `json:"chain_id"`
	TokenType  entity.TokenStatisticsType `json:"token_type"`
	IBCHops    int                        `json:"ibc_hops"`
	Amount     string                     `json:"amount"`
	ReceiveTxs int64                      `json:"receive_txs"`
}
