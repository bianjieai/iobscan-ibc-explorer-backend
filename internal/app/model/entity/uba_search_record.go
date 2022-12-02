package entity

type UbaSearchRecord struct {
	Ip       string `bson:"ip"`
	Content  string `bson:"content"`
	CreateAt int64  `bson:"create_at"`
}

func (e UbaSearchRecord) CollectionName() string {
	return "uba_search_record"
}
