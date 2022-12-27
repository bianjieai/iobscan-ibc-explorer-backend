package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IDenomHeatmap interface {
	InsertMany(batch []*entity.DenomHeatmap) error
	FindLastStatisticsTime(time time.Time) (time.Time, error)
	FindByStatisticsTime(time time.Time) ([]*entity.DenomHeatmap, error)
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

func (repo *DenomHeatmap) FindLastStatisticsTime(time time.Time) (time.Time, error) {
	var res entity.DenomHeatmap
	err := repo.coll().Find(context.Background(), bson.M{
		"statistics_time": bson.M{
			"$lte": time,
		},
	}).Sort("-statistics_time").Select(bson.M{"statistics_time": 1}).One(&res)
	return res.StatisticsTime, err
}

func (repo *DenomHeatmap) FindByStatisticsTime(time time.Time) ([]*entity.DenomHeatmap, error) {
	var res []*entity.DenomHeatmap
	err := repo.coll().Find(context.Background(), bson.M{
		"statistics_time": time,
	}).All(&res)
	return res, err
}
