package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IChainInflowStatisticsRepo interface {
	InsertMany(batch []*entity.IBCChainInflowStatistics) error
}

var _ IChainInflowStatisticsRepo = new(ChainInflowStatisticsRepo)

type ChainInflowStatisticsRepo struct {
}

func (repo *ChainInflowStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainInflowStatisticsCollName)
}

func (repo *ChainInflowStatisticsRepo) InsertMany(batch []*entity.IBCChainInflowStatistics) error {
	if len(batch) == 0 {
		return nil
	}
	if _, err := repo.coll().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}
