package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type ITokenStatisticsRepo interface {
}

var _ ITokenStatisticsRepo = new(TokenStatisticsRepo)

type TokenStatisticsRepo struct {
}

func (repo *TokenStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenStatistics{}.CollectionName())
}
