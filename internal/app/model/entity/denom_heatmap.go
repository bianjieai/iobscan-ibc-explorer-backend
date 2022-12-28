package entity

const DenomHeatmapCollName = "denom_heatmap"

type DenomHeatmap struct {
	Denom             string  `bson:"denom"`
	Chain             string  `bson:"chain"`
	StatisticsTime    string  `bson:"statistics_time"`
	Price             float64 `bson:"price"`
	Supply            float64 `bson:"supply"`
	MarketCap         float64 `bson:"market_cap"`
	TransferVolume24h float64 `bson:"transfer_volume_24h"`
	CreateAt          int64   `bson:"create_at"`
	UpdateAt          int64   `bson:"update_at"`
}

func (d DenomHeatmap) CollectionName() string {
	return DenomHeatmapCollName
}
