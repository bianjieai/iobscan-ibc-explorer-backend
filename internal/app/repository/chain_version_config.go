package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainVersionConfigRepo interface {
	FindOne(chainId string) (*entity.ChainVersionConfig, error)
	FindAll() ([]*entity.ChainVersionConfig, error)
	GetChainVerCfgMap() (map[string]string, error)
}

var _ IChainVersionConfigRepo = new(ChainVersionConfigRepo)

type ChainVersionConfigRepo struct {
}

func (repo *ChainVersionConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ChainVersionConfig{}.CollectionName())
}

func (repo *ChainVersionConfigRepo) FindOne(chainId string) (*entity.ChainVersionConfig, error) {
	var res *entity.ChainVersionConfig
	err := repo.coll().Find(context.Background(), bson.M{"chain_id": chainId}).One(&res)
	return res, err
}

func (repo *ChainVersionConfigRepo) FindAll() ([]*entity.ChainVersionConfig, error) {
	var res []*entity.ChainVersionConfig
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *ChainVersionConfigRepo) GetChainVerCfgMap() (map[string]string, error) {
	verCfgMap := make(map[string]string, 10)
	res, err := repo.FindAll()
	if err != nil {
		return nil, err
	}
	for _, val := range res {
		verCfgMap[val.ChainId] = val.ChainName
	}
	return verCfgMap, nil
}
