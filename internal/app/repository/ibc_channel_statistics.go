package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IChannelStatisticsRepo interface {
}

var _ IChannelStatisticsRepo = new(ChannelStatisticsRepo)

type ChannelStatisticsRepo struct {
}

func (repo *ChannelStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannelStatistics{}.CollectionName())
}
