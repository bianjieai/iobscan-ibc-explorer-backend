package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainConfigRepo interface {
	FindAllChains() ([]*entity.ChainConfig, error)
}

var _ IChainConfigRepo = new(ChainConfigRepo)

type ChainConfigRepo struct {
}

func (repo *ChainConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ChainConfig{}.CollectionName())
}

func (repo *ChainConfigRepo) FindAllChains() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{}).Select(bson.M{"addr_prefix": 1, "chain_name": 1}).All(&res)
	return res, err
}
