package entity

const (
	IBCTxFailLogCollName = "ibc_tx_fail_log"
)

type TxFailCode string

const (
	TxFailCodeTimeout  TxFailCode = "timeout"
	TxFailCodeOther    TxFailCode = "other"
	TxFailCodeWrongSeq TxFailCode = "wrong_sequence"
)

type IBCTxFailLog struct {
	Chain            string     `bson:"chain"`
	Log              string     `bson:"log"`
	Code             TxFailCode `bson:"code"`
	TxsNumber        int64      `bson:"txs_number"`
	SegmentStartTime int64      `bson:"segment_start_time"`
	SegmentEndTime   int64      `bson:"segment_end_time"`
	CreateAt         int64      `bson:"create_at"`
	UpdateAt         int64      `bson:"update_at"`
}

func (i *IBCTxFailLog) CollectionName() string {
	return IBCTxFailLogCollName
}
