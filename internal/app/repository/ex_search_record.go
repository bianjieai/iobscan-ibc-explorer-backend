package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IExSearchRecordRepo interface {
	Insert(record *entity.ExSearchRecord) error
}

var _ IExSearchRecordRepo = new(ExSearchRecordRepo)

type ExSearchRecordRepo struct {
}

func (repo *ExSearchRecordRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExSearchRecord{}.CollectionName())
}

func (repo *ExSearchRecordRepo) Insert(record *entity.ExSearchRecord) error {
	_, err := repo.coll().InsertOne(context.Background(), record)
	return err
}
