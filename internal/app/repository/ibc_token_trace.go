package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenTraceRepo interface {
	FindByBaseDenom(baseDenom, originChainId string) ([]*entity.IBCTokenTrace, error)
	BatchSwap(batch []*entity.IBCTokenTrace, baseDenom, originChainId string) error
	AggregateIBCChain() ([]*dto.AggregateIBCChainDTO, error)
	List(baseDenom, chainId string, tokenType entity.TokenTraceType, skip, limit int64) ([]*entity.IBCTokenTrace, error)
	CountList(baseDenom, chainId string, tokenType entity.TokenTraceType) (int64, error)
}

var _ ITokenTraceRepo = new(TokenTraceRepo)

type TokenTraceRepo struct {
}

func (repo *TokenTraceRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenTrace{}.CollectionName())
}

//func (repo *TokenTraceRepo) EnsureIndexes() {
//	var indexes []options.IndexModel
//	indexes = append(indexes, options.IndexModel{
//		Key:          []string{"denom", "chain_id"},
//		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
//	})
//
//	ensureIndexes(entity.IBCTokenTrace{}.CollectionName(), indexes)
//}

func (repo *TokenTraceRepo) FindByBaseDenom(baseDenom, chainId string) ([]*entity.IBCTokenTrace, error) {
	var res []*entity.IBCTokenTrace
	qurey := bson.M{"origional_id": chainId,
		"base_denom": baseDenom,
	}
	err := repo.coll().Find(context.Background(), qurey).All(&res)
	return res, err
}

func (repo *TokenTraceRepo) BatchSwap(batch []*entity.IBCTokenTrace, baseDenom, originChainId string) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"base_denom": baseDenom,
		}
		if _, err := repo.coll().RemoveAll(sessCtx, query); err != nil {
			return nil, err
		}

		if len(batch) == 0 {
			return nil, nil
		}

		for _, v := range batch {
			v.CreateAt = time.Now().Unix()
			v.UpdateAt = time.Now().Unix()
		}
		if _, err := repo.coll().InsertMany(sessCtx, batch); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := mgo.DoTransaction(context.Background(), callback)
	return err
}

func (repo *TokenTraceRepo) AggregateIBCChain() ([]*dto.AggregateIBCChainDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"type": bson.M{
				"$in": []entity.TokenTraceType{entity.TokenTraceTypeAuthed, entity.TokenTraceTypeOther},
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

func (repo *TokenTraceRepo) analyzeListParam(baseDenom, chainId string, tokenType entity.TokenTraceType) map[string]interface{} {
	q := make(map[string]interface{})
	q["base_denom"] = baseDenom

	if chainId != "" {
		q["chain_id"] = chainId
	}

	if tokenType != "" {
		q["type"] = tokenType
	}

	return q
}

func (repo *TokenTraceRepo) List(baseDenom, chainId string, tokenType entity.TokenTraceType, skip, limit int64) ([]*entity.IBCTokenTrace, error) {
	param := repo.analyzeListParam(baseDenom, chainId, tokenType)
	var res []*entity.IBCTokenTrace
	err := repo.coll().Find(context.Background(), param).Limit(limit).Skip(skip).Sort("-denom_amount").All(&res)
	return res, err
}

func (repo *TokenTraceRepo) CountList(baseDenom, chainId string, tokenType entity.TokenTraceType) (int64, error) {
	param := repo.analyzeListParam(baseDenom, chainId, tokenType)
	count, err := repo.coll().Find(context.Background(), param).Count()
	return count, err
}
