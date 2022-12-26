package vo

type MarketHeatmapResp struct {
}

type VolumeItem struct {
	Datetime string `json:"datetime"`
	Value    string `json:"value"`
}

type TokenDistributionReq struct {
	BaseDenom      string `json:"base_denom" form:"base_denom"`
	BaseDenomChain string `json:"base_denom_chain" form:"base_denom_chain"`
}

type TokenDistributionResp struct {
	Children []TokenDistributionResp `json:"children"`
	Amount   string                  `json:"amount"`
	Denom    string                  `json:"denom"`
	Symbol   string                  `json:"symbol"`
	Chain    string                  `json:"chain"`
}
