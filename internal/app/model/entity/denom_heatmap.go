package entity

import "time"

const DenomHeatmapCollName = "denom_heatmap"

type DenomHeatmap struct {
	Denom             string    `bson:"denom"`
	Chain             string    `bson:"chain"`
	StatisticsTime    time.Time `bson:"statistics_time"`
	Price             float64   `bson:"price"`
	Supply            string    `bson:"supply"`
	MarketCap         string    `bson:"market_cap"`
	TransferVolume24h string    `bson:"transfer_volume_24h"`
	CreateAt          int64     `bson:"create_at"`
	UpdateAt          int64     `bson:"update_at"`
}

func (d DenomHeatmap) CollectionName() string {
	return DenomHeatmapCollName
}
