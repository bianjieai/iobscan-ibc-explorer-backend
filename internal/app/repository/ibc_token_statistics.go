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

// ITokenStatisticsRepo
// Warning: ITokenStatisticsRepo is deprecated
type ITokenStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCTokenStatistics) error
	BatchInsert(batch []*entity.IBCTokenStatistics) error
	BatchInsertToNew(batch []*entity.IBCTokenStatistics) error
	Aggr() ([]*dto.CountBaseDenomTxsDTO, error)
	FindEmptyBaseDenomChainItems(skip, limit int64) ([]*entity.IBCTokenStatistics, error)
}

var _ ITokenStatisticsRepo = new(TokenStatisticsRepo)

// TokenStatisticsRepo
// Warning: TokenStatisticsRepo is deprecated
type TokenStatisticsRepo struct {
}

func (repo *TokenStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenStatisticsCollName)
}

func (repo *TokenStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenStatisticsNewCollName)
}

func (repo *TokenStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true)
	key := []string{"base_denom", "base_denom_chain", "-segment_start_time", "-segment_end_time"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})
}

func (repo *TokenStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCTokenStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCTokenStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
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

func (repo *TokenStatisticsRepo) BatchInsertToNew(batch []*entity.IBCTokenStatistics) error {
	if len(batch) == 0 {
		return nil
	}

	_, err := repo.collNew().InsertMany(context.Background(), batch)
	return err
}

func (repo *TokenStatisticsRepo) Aggr() ([]*dto.CountBaseDenomTxsDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":       "$base_denom",
				"base_denom_chain": "$base_denom_chain",
			},
			"count": bson.M{
				"$sum": "$transfer_txs",
			},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"_id":              0,
			"base_denom":       "$_id.base_denom",
			"base_denom_chain": "$_id.base_denom_chain",
			"count":            "$count",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group, project)
	var res []*dto.CountBaseDenomTxsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *TokenStatisticsRepo) FindEmptyBaseDenomChainItems(skip, limit int64) ([]*entity.IBCTokenStatistics, error) {
	var res []*entity.IBCTokenStatistics
	err := repo.coll().Find(context.Background(), bson.M{"base_denom_chain": ""}).Skip(skip).Limit(limit).All(&res)
	return res, err
}
