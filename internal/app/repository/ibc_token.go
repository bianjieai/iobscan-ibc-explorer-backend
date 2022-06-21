package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenRepo interface {
	FindAll() (entity.IBCTokenList, error)
}

var _ ITokenRepo = new(TokenRepo)

type TokenRepo struct {
}

func (repo *TokenRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCToken{}.CollectionName())
}

func (repo *TokenRepo) FindAll() (entity.IBCTokenList, error) {
	var res entity.IBCTokenList
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}
