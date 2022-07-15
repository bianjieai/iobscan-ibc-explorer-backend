package model

type (
	Fee struct {
		Amount []*Coin `bson:"amount"`
		Gas    int64   `bson:"gas"`
	}

	Coin struct {
		Denom  string `bson:"denom" json:"denom"`
		Amount string `bson:"amount" json:"amount"`
	}

	TxMsg struct {
		Type string `bson:"type"`
		Msg  struct {
			PacketId string `bson:"packet_id"`
			Signer   string `bson:"signer"`
		} `bson:"msg"`
	}
)
