package vo

type PopularSymbolsReq struct {
	MinHops       string `json:"min_hops" form:"min_hops" binding:"required"`
	MinReceiveTxs string `json:"min_receive_txs" form:"min_receive_txs" binding:"required"`
}

type PopularSymbolsResp struct {
	TimeStamp int64          `json:"timestamp"`
	Symbols   []SymbolDetail `json:"symbols"`
}

type SymbolDetail struct {
	Symbol          string `json:"symbol"`
	TotalReceiveTxs int64  `json:"total_receive_txs"`
}
