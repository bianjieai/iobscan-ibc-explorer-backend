package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainRepo interface {
}

var _ IChainRepo = new(IbcChainRepo)

type IbcChainRepo struct {
}

func (repo *IbcChainRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChain{}.CollectionName())
}

func (repo *IbcChainRepo) FindAll() ([]*entity.IBCChain, error) {
	var res []*entity.IBCChain
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}
