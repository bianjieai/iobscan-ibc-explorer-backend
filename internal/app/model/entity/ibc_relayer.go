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
	return i.ChainA != "" && i.ChainB != "" && i.ChannelA != "" && i.ChannelB != "" && i.ChainBAddress != "" && i.ChainAAddress != ""
}

type ChannelPairInfoList []ChannelPairInfo

func (c ChannelPairInfoList) GetChainAddrs() []string {
	addrMap := make(map[string]struct{})
	for _, v := range c {
		addrMap[v.ChainAAddress] = struct{}{}
		addrMap[v.ChainBAddress] = struct{}{}
	}

	addrs := make([]string, 0, len(addrMap))
	for k, _ := range addrMap {
		addrs = append(addrs, k)
	}
	return addrs
}

func (c ChannelPairInfoList) GetChains() []string {
	chainMap := make(map[string]struct{})
	for _, v := range c {
		chainMap[v.ChainA] = struct{}{}
		chainMap[v.ChainB] = struct{}{}
	}

	chains := make([]string, 0, len(chainMap))
	for k, _ := range chainMap {
		chains = append(chains, k)
	}
	return chains
}
