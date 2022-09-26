package entity

type ExSearchRecord struct {
	Ip       string `bson:"ip"`
	Content  string `bson:"content"`
	CreateAt int64  `bson:"create_at"`
}

func (e ExSearchRecord) CollectionName() string {
	return "ex_search_record"
}
