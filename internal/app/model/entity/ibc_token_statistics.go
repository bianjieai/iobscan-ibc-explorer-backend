package entity

type TokenStatisticsType string

const (
	TokenStatisticsTypeGenesis = "Genesis"
	TokenStatisticsTypeAuthed  = "Authed"
	TokenStatisticsTypeOther   = "Other"
)

type IBCTokenStatistics struct {
	Denom       string              `bson:"denom"`
	DenomPath   string              `bson:"denom_path"`
	BaseDenom   string              `bson:"base_denom"`
	ChainId     string              `bson:"chain_id"`
	OriginalId  string              `bson:"original_id"`
	Type        TokenStatisticsType `bson:"type"`
	IBCHops     int                 `bson:"ibc_hops"`
	DenomAmount string              `bson:"denom_amount"`
	DenomValue  float64             `bson:"denom_value"`
	ReceiveTxs  int64               `bson:"receive_txs"`
	CreateAt    int64               `bson:"create_at"`
	UpdateAt    int64               `bson:"update_at"`
}

func (i IBCTokenStatistics) CollectionName() string {
	return "ibc_token_statistics"
}
