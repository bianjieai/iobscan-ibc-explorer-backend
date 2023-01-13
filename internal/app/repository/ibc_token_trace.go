package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenTraceRepo interface {
	DelByBaseDenom(baseDenom, BaseDenomChain string) error
	BatchSwap(batch []*entity.IBCTokenTrace, baseDenom, BaseDenomChain string) error
	AggregateIBCChain() ([]*dto.AggregateIBCChainDTO, error)
	List(req *vo.IBCTokenListReq) ([]*entity.IBCTokenTrace, error)
	CountList(req *vo.IBCTokenListReq) (int64, error)
	FindByBaseDenom(baseDenom, baseDenomChain string) ([]*entity.IBCTokenTrace, error)
}

var _ ITokenTraceRepo = new(TokenTraceRepo)

type TokenTraceRepo struct {
}

func (repo *TokenTraceRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCTokenTrace{}.CollectionName())
}

func (repo *TokenTraceRepo) DelByBaseDenom(baseDenom, BaseDenomChain string) error {
	query := bson.M{
		"base_denom":       baseDenom,
		"base_denom_chain": BaseDenomChain,
	}
	_, err := repo.coll().RemoveAll(context.Background(), query)
	return err
}

func (repo *TokenTraceRepo) BatchSwap(batch []*entity.IBCTokenTrace, baseDenom, BaseDenomChain string) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"base_denom":       baseDenom,
			"base_denom_chain": BaseDenomChain,
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
			"_id": "$chain",
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

func (repo *TokenTraceRepo) analyzeListParam(req *vo.IBCTokenListReq) map[string]interface{} {
	q := make(map[string]interface{})
	q["base_denom"] = req.BaseDenom
	q["base_denom_chain"] = req.BaseDenomChain

	if req.Chain != "" {
		q["chain"] = req.Chain
	}

	if req.TokenType != "" {
		q["type"] = req.TokenType
	}

	return q
}

func (repo *TokenTraceRepo) List(req *vo.IBCTokenListReq) ([]*entity.IBCTokenTrace, error) {
	param := repo.analyzeListParam(req)
	skip, limit := vo.ParseParamPage(req.PageNum, req.PageSize)
	var res []*entity.IBCTokenTrace
	err := repo.coll().Find(context.Background(), param).Limit(limit).Skip(skip).All(&res)
	return res, err
}

func (repo *TokenTraceRepo) CountList(req *vo.IBCTokenListReq) (int64, error) {
	param := repo.analyzeListParam(req)
	count, err := repo.coll().Find(context.Background(), param).Count()
	return count, err
}

func (repo *TokenTraceRepo) FindByBaseDenom(baseDenom, baseDenomChain string) ([]*entity.IBCTokenTrace, error) {
	q := bson.M{
		"base_denom":       baseDenom,
		"base_denom_chain": baseDenomChain,
	}
	var res []*entity.IBCTokenTrace
	err := repo.coll().Find(context.Background(), q).All(&res)
	return res, err
}
