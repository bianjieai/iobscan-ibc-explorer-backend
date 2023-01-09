package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IAuthDenomRepo interface {
	FindAll() (entity.AuthDenomList, error)
	FindBySymbol(symbol string) (entity.AuthDenom, error)
	FindStableCoins() (entity.AuthDenomList, error)
}

var _ IAuthDenomRepo = new(AuthDenomRepo)

type AuthDenomRepo struct {
}

func (repo *AuthDenomRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.AuthDenom{}.CollectionName())
}

func (repo *AuthDenomRepo) FindAll() (entity.AuthDenomList, error) {
	var res entity.AuthDenomList
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *AuthDenomRepo) FindBySymbol(symbol string) (entity.AuthDenom, error) {
	var res entity.AuthDenom
	err := repo.coll().Find(context.Background(), bson.M{"symbol": symbol}).One(&res)
	return res, err
}

func (repo *AuthDenomRepo) FindStableCoins() (entity.AuthDenomList, error) {
	var res entity.AuthDenomList
	err := repo.coll().Find(context.Background(), bson.M{"is_stable_coin": true}).All(&res)
	return res, err
}
