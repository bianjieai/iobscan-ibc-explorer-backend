package entity

type TokenTraceType string

const (
	TokenTraceTypeGenesis TokenTraceType = "Genesis"
	TokenTraceTypeAuthed  TokenTraceType = "Authed"
	TokenTraceTypeOther   TokenTraceType = "Other"
)

type IBCTokenTrace struct {
	Denom       string         `bson:"denom"`
	DenomPath   string         `bson:"denom_path"`
	BaseDenom   string         `bson:"base_denom"`
	ChainId     string         `bson:"chain_id"`
	OriginalId  string         `bson:"original_id"`
	Type        TokenTraceType `bson:"type"`
	IBCHops     int            `bson:"ibc_hops"`
	DenomAmount string         `bson:"denom_amount"`
	DenomValue  string         `bson:"denom_value"`
	ReceiveTxs  int64          `bson:"receive_txs"`
	CreateAt    int64          `bson:"create_at"`
	UpdateAt    int64          `bson:"update_at"`
}

func (i IBCTokenTrace) CollectionName() string {
	return "ibc_token_trace"
}
