package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IRelayerDenomStatisticsRepo interface {
}

var _ IRelayerDenomStatisticsRepo = new(RelayerDenomStatisticsRepo)

type RelayerDenomStatisticsRepo struct {
}

func (repo *RelayerDenomStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerDenomStatistics{}.CollectionName())
}
