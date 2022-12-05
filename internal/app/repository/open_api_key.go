package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IOpenApiKeyRepo interface {
	FindByApiKey(apiKey string) (entity.OpenApiKey, error)
}

var _ IOpenApiKeyRepo = new(OpenApiKeyRepo)

type OpenApiKeyRepo struct {
}

func (repo *OpenApiKeyRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.OpenApiKeyCollName)
}

func (repo *OpenApiKeyRepo) FindByApiKey(apiKey string) (entity.OpenApiKey, error) {
	var res entity.OpenApiKey
	err := repo.coll().Find(context.Background(), bson.M{"api_key": apiKey}).One(&res)
	return res, err
}
