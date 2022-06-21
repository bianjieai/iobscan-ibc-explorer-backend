package entity

type IBCChain struct {
	ChainId          string `bson:"chain_id"`
	ConnectedChains  int64  `bson:"connected_chains"`
	Channels         int64  `bson:"channels"`
	Relayers         int64  `bson:"relayers"`
	IbcTokens        int64  `bson:"ibc_tokens"`
	IbcTokensValue   string `bson:"ibc_tokens_value"`
	TransferTxs      int64  `bson:"transfer_txs"`
	TransferTxsValue string `bson:"transfer_txs_value"`
	CreateAt         string `bson:"create_at"`
	UpdateAt         string `bson:"update_at"`
}

func (i IBCChain) CollectionName() string {
	return "ibc_chain"
}
