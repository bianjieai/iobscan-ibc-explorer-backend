package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type ITokenRepo interface {
}

var _ ITokenRepo = new(TokenRepo)

type TokenRepo struct {
}

func (repo *TokenRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCToken{}.CollectionName())
}
