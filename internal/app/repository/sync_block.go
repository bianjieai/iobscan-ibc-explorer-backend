package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ISyncBlockRepo interface {
	FindLatestBlock(chainId string) (*entity.SyncBlock, error)
}

var _ ISyncBlockRepo = new(SyncBlockRepo)

type SyncBlockRepo struct {
}

func (repo *SyncBlockRepo) coll(chainId string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.SyncBlock{}.CollectionName(chainId))
}

func (repo *SyncBlockRepo) FindLatestBlock(chainId string) (*entity.SyncBlock, error) {
	var res entity.SyncBlock
	err := repo.coll(chainId).Find(context.Background(), bson.M{}).Sort("-height").Limit(1).One(&res)
	return &res, err
}
