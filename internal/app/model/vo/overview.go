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

type (
	TokenDistributionResp GraphData
	GraphData             struct {
		Children []*GraphData `json:"children"`
		Amount   string       `json:"amount"`
		Denom    string       `json:"denom"`
		Chain    string       `json:"chain"`
		Hops     int          `json:"hops"`
	}
)
