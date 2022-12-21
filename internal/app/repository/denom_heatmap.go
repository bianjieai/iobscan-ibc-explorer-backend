package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IDenomHeatmap interface {
	InsertMany(batch []*entity.DenomHeatmap) error
}

var _ IDenomHeatmap = new(DenomHeatmap)

type DenomHeatmap struct {
}

func (repo *DenomHeatmap) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.DenomHeatmapCollName)
}

func (repo *DenomHeatmap) InsertMany(batch []*entity.DenomHeatmap) error {
	if len(batch) == 0 {
		return nil
	}
	if _, err := repo.coll().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}
