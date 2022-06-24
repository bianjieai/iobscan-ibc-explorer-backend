package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

type ITokenRepo interface {
	FindAll() (entity.IBCTokenList, error)
	InsertBatch(batch []*entity.IBCToken) error
	UpdateToken(token *entity.IBCToken) error
}

var _ ITokenRepo = new(TokenRepo)

type TokenRepo struct {
}

func (repo *TokenRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"base_denom", "chain_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})

	ensureIndexes(entity.IBCToken{}.CollectionName(), indexes)
}

func (repo *TokenRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCToken{}.CollectionName())
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
		"chain_id":   token.ChainId,
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
