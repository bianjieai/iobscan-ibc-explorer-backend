package entity

type (
	IbcStatistic struct {
		StatisticsName string `bson:"statistics_name"`
		StatisticsInfo string `bson:"statistics_info"`
		Count          int64  `bson:"count"`
		UpdateAt       int64  `bson:"update_at"`
		CreateAt       int64  `bson:"create_at"`
	}
)

func (ibc IbcStatistic) CollectionName() string {
	return "ibc_statistics"
}
