package repository

import (
	"context"
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"
)

type IChainInflowStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InsertMany(batch []*entity.IBCChainInflowStatistics) error
	InsertManyToNew(batch []*entity.IBCChainInflowStatistics) error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainInflowStatistics) error
	BatchSwapNew(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainInflowStatistics) error
	AggrTrend(chain string, segmentStartTime, segmentEndTime int64) ([]*dto.AggrChainInflowTrendDTO, error)
}

var _ IChainInflowStatisticsRepo = new(ChainInflowStatisticsRepo)

type ChainInflowStatisticsRepo struct {
}

func (repo *ChainInflowStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainInflowStatisticsCollName)
}

func (repo *ChainInflowStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainInflowStatisticsNewCollName)
}

func (repo *ChainInflowStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index()
	key := []string{"segment_start_time", "segment_end_time"}
	if err := repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts}); err != nil {
		return err
	}

	return nil
}

func (repo *ChainInflowStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCChainInflowStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCChainInflowStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
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

func (repo *ChainInflowStatisticsRepo) InsertManyToNew(batch []*entity.IBCChainInflowStatistics) error {
	if len(batch) == 0 {
		return nil
	}
	if _, err := repo.collNew().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *ChainInflowStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainInflowStatistics) error {
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

func (repo *ChainInflowStatisticsRepo) BatchSwapNew(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChainInflowStatistics) error {
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

func (repo *ChainInflowStatisticsRepo) AggrTrend(chain string, segmentStartTime, segmentEndTime int64) ([]*dto.AggrChainInflowTrendDTO, error) {
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
			//"txs_number": bson.M{
			//	"$sum": "$txs_number",
			//},
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
			//"txs_number":         "$txs_number",
			"denom_amount": "$denom_amount",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.AggrChainInflowTrendDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
