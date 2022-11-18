package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainRegistryRepo interface {
	FindAll() ([]*entity.ChainRegistry, error)
	FindOne(chain string) (*entity.ChainRegistry, error)
}

var _ IChainRegistryRepo = new(ChainRegistryRepo)

type ChainRegistryRepo struct {
}

func (repo *ChainRegistryRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ChainRegistry{}.CollectionName())
}

func (repo *ChainRegistryRepo) FindAll() ([]*entity.ChainRegistry, error) {
	var res []*entity.ChainRegistry
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *ChainRegistryRepo) FindOne(chain string) (*entity.ChainRegistry, error) {
	var res *entity.ChainRegistry
	err := repo.coll().Find(context.Background(), bson.M{"chain": chain}).One(&res)
	return res, err
}
