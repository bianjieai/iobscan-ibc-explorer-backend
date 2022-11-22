package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IUbaSearchRecordRepo interface {
	Insert(record *entity.UbaSearchRecord) error
}

var _ IUbaSearchRecordRepo = new(UbaSearchRecordRepo)

type UbaSearchRecordRepo struct {
}

func (repo *UbaSearchRecordRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.UbaSearchRecord{}.CollectionName())
}

func (repo *UbaSearchRecordRepo) Insert(record *entity.UbaSearchRecord) error {
	_, err := repo.coll().InsertOne(context.Background(), record)
	return err
}
