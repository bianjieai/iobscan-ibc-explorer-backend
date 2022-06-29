package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IDenomCaculateRepo interface {
	FindByBaseDenom(baseDenom string) ([]*entity.IBCDenomCaculate, error)
}

var _ IDenomCaculateRepo = new(DenomCaculateRepo)

type DenomCaculateRepo struct {
}

func (repo *DenomCaculateRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenomCaculate{}.CollectionName())
}

func (repo *DenomCaculateRepo) FindByBaseDenom(baseDenom string) ([]*entity.IBCDenomCaculate, error) {
	var res []*entity.IBCDenomCaculate
	qurey := bson.M{"base_denom": baseDenom}
	err := repo.coll().Find(context.Background(), qurey).All(&res)
	return res, err
}
