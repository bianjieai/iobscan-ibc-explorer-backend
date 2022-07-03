package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenTraceStatisticsRepo interface {
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCTokenTraceStatistics) error
	BatchInsert(batch []*entity.IBCTokenTraceStatistics) error
	Aggr() ([]*dto.TokenTraceStatisticsDTO, error)
}

var _ ITokenTraceStatisticsRepo = new(TokenTraceStatisticsRepo)

type TokenTraceStatisticsRepo struct {
}

func (repo *TokenTraceStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenTraceStatistics{}.CollectionName())
}

func (repo *TokenTraceStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCTokenTraceStatistics) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"segment_start_time": segmentStartTime,
			"segment_end_time":   segmentEndTime,
		}
		if _, err := repo.coll().RemoveAll(sessCtx, query); err != nil {
			return nil, err
		}

		if len(batch) == 0 {
			return nil, nil
		}

		if _, err := repo.coll().InsertMany(sessCtx, batch); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := mgo.DoTransaction(context.Background(), callback)
	return err
}

func (repo *TokenTraceStatisticsRepo) BatchInsert(batch []*entity.IBCTokenTraceStatistics) error {
	if len(batch) == 0 {
		return nil
	}

	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *TokenTraceStatisticsRepo) Aggr() ([]*dto.TokenTraceStatisticsDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"denom":    "$denom",
				"chain_id": "$chain_id",
			},
			"receive_txs": bson.M{
				"$sum": "$receive_txs",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":         0,
			"denom":       "$_id.denom",
			"chain_id":    "$_id.chain_id",
			"receive_txs": "$receive_txs",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group, project)
	var res []*dto.TokenTraceStatisticsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
