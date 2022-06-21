package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IDenomRepo interface {
	FindBaseDenom() (entity.IBCDenomList, error)
}

var _ IDenomRepo = new(DenomRepo)

type DenomRepo struct {
}

func (repo *DenomRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenom{}.CollectionName())
}

func (repo *DenomRepo) FindBaseDenom() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"is_base_denom": true}).All(&res)
	return res, err
}
