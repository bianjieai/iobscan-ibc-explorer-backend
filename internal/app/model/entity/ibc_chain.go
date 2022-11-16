package entity

type IBCChain struct {
	Chain            string `bson:"chain"`
	ConnectedChains  int64  `bson:"connected_chains"`
	Channels         int64  `bson:"channels"`
	Relayers         int64  `bson:"relayers"`
	IbcTokens        int64  `bson:"ibc_tokens"`
	IbcTokensValue   string `bson:"ibc_tokens_value"`
	TransferTxs      int64  `bson:"transfer_txs"`
	TransferTxsValue string `bson:"transfer_txs_value"`
	CreateAt         int64  `bson:"create_at"`
	UpdateAt         int64  `bson:"update_at"`
}

func (i IBCChain) CollectionName() string {
	return "ibc_chain"
}
