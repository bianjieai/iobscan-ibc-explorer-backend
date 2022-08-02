package entity

type (
	IbcStatistic struct {
		StatisticsName string `bson:"statistics_name"`
		Count          int64  `bson:"count"`
		Data           string `bson:"data"`
		UpdateAt       int64  `bson:"update_at"`
		CreateAt       int64  `bson:"create_at"`
	}
)

func (ibc IbcStatistic) CollectionName() string {
	return "ibc_statistics"
}
