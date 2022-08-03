package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ISyncTaskRepo interface {
	CheckFollowingStatus(chainId string) (bool, error)
}

var _ ISyncTaskRepo = new(SyncTaskRepo)

type SyncTaskRepo struct {
}

func (repo *SyncTaskRepo) coll(chainId string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.SyncTask{}.CollectionName(chainId))
}

func (repo *SyncTaskRepo) CheckFollowingStatus(chainId string) (bool, error) {
	count, err := repo.coll(chainId).Find(context.Background(), bson.M{"status": entity.SyncTaskStatusUnderway, "end_height": 0}).Count()
	return count == 1, err
}
