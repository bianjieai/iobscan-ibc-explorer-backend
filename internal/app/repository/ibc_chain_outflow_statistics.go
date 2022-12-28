package repository

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IChainOutflowStatisticsRepo interface {
	InsertMany(batch []*entity.IBCChainOutflowStatistics) error
	InsertManyToNew(batch []*entity.IBCChainOutflowStatistics) error
	AggrTrend(chain string, segmentStartTime, segmentEndTime int64) ([]*dto.AggrChainOutflowTrendDTO, error)
	CreateNew() error
	SwitchColl() error
	BatchSwapNew(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainOutflowStatistics) error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainOutflowStatistics) error
}

var _ IChainOutflowStatisticsRepo = new(ChainOutflowStatisticsRepo)

type ChainOutflowStatisticsRepo struct {
}

func (repo *ChainOutflowStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainOutflowStatisticsCollName)
}

func (repo *ChainOutflowStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainOutflowStatisticsNewCollName)
}

func (repo *ChainOutflowStatisticsRepo) InsertManyToNew(batch []*entity.IBCChainOutflowStatistics) error {
	if len(batch) == 0 {
		return nil
	}
	_, err := repo.collNew().InsertMany(context.Background(), batch)
	return err
}

func (repo *ChainOutflowStatisticsRepo) InsertMany(batch []*entity.IBCChainOutflowStatistics) error {
	if len(batch) == 0 {
		return nil
	}
	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *ChainOutflowStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCChainOutflowStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCChainOutflowStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}

func (repo *ChainOutflowStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index()
	key := []string{"segment_start_time", "segment_end_time"}
	if err := repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts}); err != nil {
		return err
	}

	return nil
}

func (repo *ChainOutflowStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainOutflowStatistics) error {
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

func (repo *ChainOutflowStatisticsRepo) BatchSwapNew(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainOutflowStatistics) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"segment_start_time": segmentStartTime,
			"segment_end_time":   segmentEndTime,
		}
		if _, err := repo.collNew().RemoveAll(sessCtx, query); err != nil {
			return nil, err
		}

		if len(batch) == 0 {
			return nil, nil
		}

		if _, err := repo.collNew().InsertMany(sessCtx, batch); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := mgo.DoTransaction(context.Background(), callback)
	return err
}

func (repo *ChainOutflowStatisticsRepo) AggrTrend(chain string, segmentStartTime, segmentEndTime int64) ([]*dto.AggrChainOutflowTrendDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"chain": chain,
			"status": bson.M{
				"$in": []entity.IbcTxStatus{entity.IbcTxStatusSuccess, entity.IbcTxStatusProcessing},
			},
			"segment_start_time": bson.M{"$gte": segmentStartTime, "$lte": segmentEndTime},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":         "$base_denom",
				"base_denom_chain":   "$base_denom_chain",
				"segment_start_time": "$segment_start_time",
			},
			"denom_amount": bson.M{
				"$sum": "$denom_amount",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":                0,
			"base_denom":         "$_id.base_denom",
			"base_denom_chain":   "$_id.base_denom_chain",
			"segment_start_time": "$_id.segment_start_time",
			"denom_amount":       "$denom_amount",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.AggrChainOutflowTrendDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
