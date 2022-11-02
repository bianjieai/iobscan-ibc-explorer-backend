package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IRelayerFeeStatisticsRepo interface {
}

var _ IRelayerFeeStatisticsRepo = new(RelayerFeeStatisticsRepo)

type RelayerFeeStatisticsRepo struct {
}

func (repo *RelayerFeeStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerFeeStatistics{}.CollectionName())
}
