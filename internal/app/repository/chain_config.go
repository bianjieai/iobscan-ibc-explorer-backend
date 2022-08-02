package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainConfigRepo interface {
	FindAll() ([]*entity.ChainConfig, error)
	FindOne(chainId string) (*entity.ChainConfig, error)
	FindAllChainIds() ([]*entity.ChainConfig, error)
}

var _ IChainConfigRepo = new(ChainConfigRepo)

type ChainConfigRepo struct {
}

func (repo *ChainConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ChainConfig{}.CollectionName())
}

func (repo *ChainConfigRepo) FindAll() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindOne(chainId string) (*entity.ChainConfig, error) {
	var res *entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{"chain_id": chainId}).One(&res)
	return res, err
}

func (repo *ChainConfigRepo) FindAllChainIds() ([]*entity.ChainConfig, error) {
	var res []*entity.ChainConfig
	err := repo.coll().Find(context.Background(), bson.M{}).Select(bson.M{"chain_id": 1, "addr_prefix": 1, "chain_name": 1}).All(&res)
	return res, err
}
