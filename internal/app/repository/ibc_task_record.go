package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITaskRecordRepo interface {
	FindByTaskName(taskName string) (*entity.IbcTaskRecord, error)
	Insert(record *entity.IbcTaskRecord) error
	UpdateHeight(taskName string, height int64) error
}

var _ ITaskRecordRepo = new(TaskRecordRepo)

type TaskRecordRepo struct {
}

func (repo *TaskRecordRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IbcTaskRecord{}.CollectionName())
}

func (repo *TaskRecordRepo) FindByTaskName(taskName string) (*entity.IbcTaskRecord, error) {
	var res entity.IbcTaskRecord
	err := repo.coll().Find(context.Background(), bson.M{"task_name": taskName}).One(&res)
	return &res, err
}

func (repo *TaskRecordRepo) Insert(record *entity.IbcTaskRecord) error {
	if record == nil {
		return nil
	}

	_, err := repo.coll().InsertOne(context.Background(), record)
	return err
}

func (repo *TaskRecordRepo) UpdateHeight(taskName string, height int64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{"task_name": taskName}, bson.M{
		"$set": bson.M{
			"height":    height,
			"update_at": time.Now().Unix(),
		},
	})
}
