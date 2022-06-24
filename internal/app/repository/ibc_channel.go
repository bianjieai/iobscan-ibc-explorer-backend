package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IChannelRepo interface {
	UpdateOne(filter interface{}, update interface{}) error
}

var _ IChannelRepo = new(ChannelRepo)

type ChannelRepo struct {
}

func (repo *ChannelRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannel{}.CollectionName())
}

func (repo *ChannelRepo) UpdateOne(filter interface{}, update interface{}) error {
	return repo.coll().UpdateOne(context.Background(), filter, update)
}
