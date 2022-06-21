package entity

type IBCRelayer struct {
	RelayerId             string `bson:"relayer_id"`
	ChainA                string `bson:"chain_a"`
	ChainB                string `bson:"chain_b"`
	ChannelA              string `bson:"channel_a"`
	ChannelB              string `bson:"channel_b"`
	ChainAAddress         string `bson:"chain_a_address"`
	ChainBAddress         string `bson:"chain_b_address"`
	TimePeriod            int    `bson:"time_period"`
	Status                int    `bson:"status"`
	ChannelAUpdateAt      int64  `bson:"channel_a_update_at"`
	ChannelBUpdateAt      int64  `bson:"channel_b_update_at"`
	TransferTotalTxs      int64  `bson:"transfer_total_txs"`
	TransferSuccessTxs    int64  `bson:"transfer_success_txs"`
	TransferTotalTxsValue string `bson:"transfer_total_txs_value"`
	CreateAt              string `bson:"create_at"`
	UpdateAt              string `bson:"update_at"`
}

func (i IBCRelayer) CollectionName() string {
	return "ibc_relayer"
}
