package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITokenRepo interface {
	FindAll() (entity.IBCTokenList, error)
	InsertBatch(batch []*entity.IBCToken) error
	UpdateToken(token *entity.IBCToken) error
}

var _ ITokenRepo = new(TokenRepo)

type TokenRepo struct {
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
	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *TokenRepo) UpdateToken(token *entity.IBCToken) error {
	query := bson.M{
		"base_denom": token.BaseDenom,
		"chain_id":   token.ChainId,
	}
	return repo.coll().UpdateOne(context.Background(), query, token)
}
