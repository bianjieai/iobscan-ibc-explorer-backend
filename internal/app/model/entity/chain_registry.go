package entity

type ChainRegistry struct {
	ChainId      string `bson:"chain_id"`
	ChainJsonUrl string `bson:"chain_json_url"`
}

func (c ChainRegistry) CollectionName() string {
	return "chain_registry"
}
