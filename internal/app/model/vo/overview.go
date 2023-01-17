package vo

type MarketHeatmapResp struct {
	Items     []HeatmapItem    `json:"items"`
	TotalInfo HeatmapTotalInfo `json:"total_info"`
}

type HeatmapItem struct {
	Price               float64 `json:"price"`
	PriceGrowthRate     float64 `json:"price_growth_rate"`
	PriceTrend          string  `json:"price_trend"`
	Denom               string  `json:"denom"`
	Chain               string  `json:"chain"`
	MarketCapValue      string  `json:"market_cap_value"`
	TransferVolumeValue string  `json:"transfer_volume_value"`
}

type HeatmapTotalInfo struct {
	StablecoinsMarketCap string  `json:"stablecoins_market_cap"`
	TotalMarketCap       string  `json:"total_market_cap"`
	TotalDenomNumber     int     `json:"total_denom_number"`
	MarketCapGrowthRate  float64 `json:"market_cap_growth_rate"`
	MarketCapTrend       string  `json:"market_cap_trend"`
	TotalTransferVolume  string  `json:"total_transfer_volume"`
	AtomPrice            float64 `json:"atom_price"`
	AtomDominance        float64 `json:"atom_dominance"`
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
	TokenDistributionResp struct {
		*GraphData
	}
	GraphData struct {
		Children []*GraphData `json:"children"`
		*GraphElem
	}

	GraphElem struct {
		Supply    string `json:"supply"`
		Amount    string `json:"amount"`
		Denom     string `json:"denom"`
		Chain     string `json:"chain"`
		Hops      int    `json:"hops"`
		PrevChain string `json:"-"`
		PrevDenom string `json:"-"`
	}
)

type ChainVolumeTrendReq struct {
	Chain string `json:"chain" form:"chain"`
}

type (
	ChainVolumeTrendResp struct {
		VolumeIn  []VolumeItem `json:"volume_in"`
		VolumeOut []VolumeItem `json:"volume_out"`
		Chain     string       `json:"chain"`
	}
)

type ChainVolumeReq struct {
}

type (
	ChainVolumeResp []ChainVolumeItem
	ChainVolumeItem struct {
		Chain               string `json:"chain"`
		TransferVolumeIn    string `json:"transfer_volume_in"`
		TransferVolumeOut   string `json:"transfer_volume_out"`
		TransferVolumeTotal string `json:"transfer_volume_total"`
	}
)
