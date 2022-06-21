package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IExIbcTxRepo interface {
}

var _ IExIbcTxRepo = new(ExIbcTxRepo)

type ExIbcTxRepo struct {
}

func (repo *ExIbcTxRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExIbcTx{}.CollectionName(false))
}

func (repo *ExIbcTxRepo) collHistory() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExIbcTx{}.CollectionName(true))
}

func (repo *ExIbcTxRepo) FindAll(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.coll().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) FindAllHistory(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.collHistory().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}
