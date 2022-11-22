package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IRelayerRepo interface {
	FindAllRelayerAddrs(skip, limit int64) ([]*entity.IBCRelayerNew, error)
}

var _ IRelayerRepo = new(IbcRelayerRepo)

type IbcRelayerRepo struct {
}

func (repo *IbcRelayerRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerNew{}.CollectionName())
}

func (repo *IbcRelayerRepo) FindAllRelayerAddrs(skip, limit int64) ([]*entity.IBCRelayerNew, error) {
	var res []*entity.IBCRelayerNew
	err := repo.coll().Find(context.Background(), bson.M{}).
		Select(bson.M{
			"channel_pair_info.chain_a_address": 1,
			"channel_pair_info.chain_b_address": 1}).Skip(skip).Limit(limit).All(&res)
	return res, err
}
