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

type ITokenTraceStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCTokenTraceStatistics) error
	BatchInsert(batch []*entity.IBCTokenTraceStatistics) error
	BatchInsertToNew(batch []*entity.IBCTokenTraceStatistics) error
	Aggr() ([]*dto.TokenTraceStatisticsDTO, error)
}

var _ ITokenTraceStatisticsRepo = new(TokenTraceStatisticsRepo)

type TokenTraceStatisticsRepo struct {
}

func (repo *TokenTraceStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenTraceStatisticsCollName)
}

func (repo *TokenTraceStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenTraceStatisticsNewCollName)
}

func (repo *TokenTraceStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true)
	key := []string{"denom", "chain", "-segment_start_time", "-segment_end_time"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})
}

func (repo *TokenTraceStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCTokenTraceStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCTokenTraceStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
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

func (repo *TokenTraceStatisticsRepo) BatchInsertToNew(batch []*entity.IBCTokenTraceStatistics) error {
	if len(batch) == 0 {
		return nil
	}

	_, err := repo.collNew().InsertMany(context.Background(), batch)
	return err
}

func (repo *TokenTraceStatisticsRepo) Aggr() ([]*dto.TokenTraceStatisticsDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"denom": "$denom",
				"chain": "$chain",
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
			"chain":       "$_id.chain",
			"receive_txs": "$receive_txs",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group, project)
	var res []*dto.TokenTraceStatisticsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
