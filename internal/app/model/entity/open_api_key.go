package entity

const OpenApiKeyCollName = "open_api_key"

type OpenApiKey struct {
	ApiKey    string `bson:"api_key" json:"api_key"`
	ApiSecret string `bson:"api_secret" json:"api_secret"`
}

func (o *OpenApiKey) CollectionName() string {
	return "open_api_key"
}
