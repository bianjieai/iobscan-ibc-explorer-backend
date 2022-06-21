package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
)

type IDenomRepo interface {
}

var _ IDenomRepo = new(DenomRepo)

type DenomRepo struct {
}

func (repo *DenomRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenom{}.CollectionName())
}
