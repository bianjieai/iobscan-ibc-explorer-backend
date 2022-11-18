package entity

const (
	CollectionNameChainVersionConfig = "chain_version_config"
)

type ChainVersionConfig struct {
	ChainId          string `bson:"chain_id"`
	ChainName        string `bson:"chain_name"`
	StartBlockHeight int64  `bson:"start_block_height"`
	EndBlockHeight   int64  `bson:"end_block_height"`
}

func (c ChainVersionConfig) CollectionName() string {
	return CollectionNameChainVersionConfig
}
