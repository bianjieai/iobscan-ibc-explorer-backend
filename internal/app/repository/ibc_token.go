package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenRepo interface {
	List(baseDenoms []string, chain string, tokenType entity.TokenType, skip, limit int64) (entity.IBCTokenList, error)
	CountList(baseDenoms []string, chain string, tokenType entity.TokenType) (int64, error)
	FindAll() (entity.IBCTokenList, error)
	InsertBatch(batch []*entity.IBCToken) error
	UpdateToken(token *entity.IBCToken) error
	Delete(baseDenom, chain string) error
}

var _ ITokenRepo = new(TokenRepo)

type TokenRepo struct {
}

func (repo *TokenRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCToken{}.CollectionName())
}

func (repo *TokenRepo) analyzeListParam(baseDenoms []string, chain string, tokenType entity.TokenType) map[string]interface{} {
	q := make(map[string]interface{})
	if len(baseDenoms) > 0 {
		q["base_denom"] = bson.M{
			"$in": baseDenoms,
		}
	}

	if chain != "" {
		q["chain"] = chain
	}

	if tokenType != "" {
		q["type"] = tokenType
	}

	return q
}

func (repo *TokenRepo) List(baseDenoms []string, chain string, tokenType entity.TokenType, skip, limit int64) (entity.IBCTokenList, error) {
	param := repo.analyzeListParam(baseDenoms, chain, tokenType)
	var res entity.IBCTokenList
	err := repo.coll().Find(context.Background(), param).Limit(limit).Skip(skip).Sort("-chains_involved").All(&res)
	return res, err
}

func (repo *TokenRepo) CountList(baseDenoms []string, chain string, tokenType entity.TokenType) (int64, error) {
	param := repo.analyzeListParam(baseDenoms, chain, tokenType)
	count, err := repo.coll().Find(context.Background(), param).Count()
	return count, err
}

func (repo *TokenRepo) FindAll() (entity.IBCTokenList, error) {
	var res entity.IBCTokenList
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *TokenRepo) InsertBatch(batch []*entity.IBCToken) error {
	if len(batch) == 0 {
		return nil
	}
	now := time.Now().Unix()
	for _, v := range batch {
		v.UpdateAt = now
		v.CreateAt = now
	}
	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *TokenRepo) UpdateToken(token *entity.IBCToken) error {
	token.UpdateAt = time.Now().Unix()
	query := bson.M{
		"base_denom": token.BaseDenom,
		"chain":      token.Chain,
	}
	update := bson.M{
		"$set": bson.M{
			"type":            token.Type,
			"price":           token.Price,
			"currency":        token.Currency,
			"supply":          token.Supply,
			"transfer_amount": token.TransferAmount,
			"transfer_txs":    token.TransferTxs,
			"chains_involved": token.ChainsInvolved,
			"update_at":       time.Now().Unix(),
		},
	}
	return repo.coll().UpdateOne(context.Background(), query, update)
}

func (repo *TokenRepo) Delete(baseDenom, chain string) error {
	return repo.coll().Remove(context.Background(), bson.M{"base_denom": baseDenom, "chain": chain})
}
