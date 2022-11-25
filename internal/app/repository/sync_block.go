package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ISyncBlockRepo interface {
	FindLatestBlock(chain string) (*entity.SyncBlock, error)
}

var _ ISyncBlockRepo = new(SyncBlockRepo)

type SyncBlockRepo struct {
}

func (repo *SyncBlockRepo) coll(chain string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.SyncBlock{}.CollectionName(chain))
}

func (repo *SyncBlockRepo) FindLatestBlock(chain string) (*entity.SyncBlock, error) {
	var res entity.SyncBlock
	err := repo.coll(chain).Find(context.Background(), bson.M{}).Sort("-height").Limit(1).One(&res)
	return &res, err
}
