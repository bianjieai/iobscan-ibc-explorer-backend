package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

type ITokenStatisticsRepo interface {
	FindByBaseDenom(baseDenom, originChainId string) ([]*entity.IBCTokenStatistics, error)
	BatchSwap(batch []*entity.IBCTokenStatistics, baseDenom, originChainId string) error
	AggregateIBCChain() ([]*dto.AggregateIBCChainDTO, error)
}

var _ ITokenStatisticsRepo = new(TokenStatisticsRepo)

type TokenStatisticsRepo struct {
}

func (repo *TokenStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenStatistics{}.CollectionName())
}

func (repo *TokenStatisticsRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"denom", "chain_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})

	ensureIndexes(entity.IBCTokenStatistics{}.CollectionName(), indexes)
}

func (repo *TokenStatisticsRepo) FindByBaseDenom(baseDenom, chainId string) ([]*entity.IBCTokenStatistics, error) {
	var res []*entity.IBCTokenStatistics
	qurey := bson.M{"origional_id": chainId,
		"base_denom": baseDenom,
	}
	err := repo.coll().Find(context.Background(), qurey).All(&res)
	return res, err
}

func (repo *TokenStatisticsRepo) BatchSwap(batch []*entity.IBCTokenStatistics, baseDenom, originChainId string) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"base_denom":  baseDenom,
			"original_id": originChainId,
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

func (repo *TokenStatisticsRepo) AggregateIBCChain() ([]*dto.AggregateIBCChainDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"type": bson.M{
				"$in": []entity.TokenStatisticsType{entity.TokenStatisticsTypeAuthed, entity.TokenStatisticsTypeOther},
			},
		},
	}

	group := bson.M{
		"$group": bson.M{
			"_id": "$chain_id",
			"denom_value": bson.M{
				"$sum": bson.M{
					"$toDouble": "$denom_value",
				},
			},
			"Count": bson.M{
				"$sum": 1,
			},
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group)
	var res []*dto.AggregateIBCChainDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
