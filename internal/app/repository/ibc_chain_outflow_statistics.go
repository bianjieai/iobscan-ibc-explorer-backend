package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IChainOutflowStatisticsRepo interface {
	InsertMany(batch []*entity.IBCChainOutflowStatistics) error
}

var _ IChainOutflowStatisticsRepo = new(ChainOutflowStatisticsRepo)

type ChainOutflowStatisticsRepo struct {
}

func (repo *ChainOutflowStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainOutflowStatisticsCollName)
}

func (repo *ChainOutflowStatisticsRepo) InsertMany(batch []*entity.IBCChainOutflowStatistics) error {
	if len(batch) == 0 {
		return nil
	}
	if _, err := repo.coll().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}
