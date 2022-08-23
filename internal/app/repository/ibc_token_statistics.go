package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenStatisticsRepo interface {
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCTokenStatistics) error
	BatchInsert(batch []*entity.IBCTokenStatistics) error
	Aggr() ([]*dto.CountBaseDenomTxsDTO, error)
}

var _ ITokenStatisticsRepo = new(TokenStatisticsRepo)

type TokenStatisticsRepo struct {
}

func (repo *TokenStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenStatistics{}.CollectionName())
}

func (repo *TokenStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCTokenStatistics) error {
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

func (repo *TokenStatisticsRepo) BatchInsert(batch []*entity.IBCTokenStatistics) error {
	if len(batch) == 0 {
		return nil
	}

	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *TokenStatisticsRepo) Aggr() ([]*dto.CountBaseDenomTxsDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":          "$base_denom",
				"base_denom_chain_id": "$base_denom_chain_id",
			},
			"count": bson.M{
				"$sum": "$transfer_txs",
			},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"_id":                 0,
			"base_denom":          "$_id.base_denom",
			"base_denom_chain_id": "$_id.base_denom_chain_id",
			"count":               "$count",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group, project)
	var res []*dto.CountBaseDenomTxsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
