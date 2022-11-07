package entity

type IBCRelayerNew struct {
	RelayerId            string            `bson:"relayer_id"`
	RelayerName          string            `bson:"relayer_name"`
	RelayerIcon          string            `bson:"relayer_icon"`
	ServedChains         int64             `bson:"served_chains"`
	ChannelPairInfo      []ChannelPairInfo `bson:"channel_pair_info"`
	UpdateTime           int64             `bson:"update_time"`
	RelayedTotalTxs      int64             `bson:"relayed_total_txs"`
	RelayedSuccessTxs    int64             `bson:"relayed_success_txs"`
	RelayedTotalTxsValue string            `bson:"relayed_total_txs_value"`
	TotalFeeValue        string            `bson:"total_fee_value"`
	CreateAt             int64             `bson:"create_at"`
	UpdateAt             int64             `bson:"update_at"`
}

type ChannelPairInfo struct {
	PairId        string `bson:"pair_id"`
	ChainA        string `bson:"chain_a"`
	ChainB        string `bson:"chain_b"`
	ChannelA      string `bson:"channel_a"`
	ChannelB      string `bson:"channel_b"`
	ChainAAddress string `bson:"chain_a_address"`
	ChainBAddress string `bson:"chain_b_address"`
}

func (i IBCRelayerNew) CollectionName() string {
	return "ibc_relayer"
}

func (i ChannelPairInfo) Valid() bool {
	return i.ChainA != "" && i.ChainB != "" && i.ChannelA != "" && i.ChannelB != "" && (i.ChainBAddress != "" || i.ChainAAddress != "")
}
