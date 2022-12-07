package entity

type ChainRegistry struct {
	Chain        string `bson:"chain"`
	ChainJsonUrl string `bson:"chain_json_url"`
}

func (c ChainRegistry) CollectionName() string {
	return "chain_registry"
}
