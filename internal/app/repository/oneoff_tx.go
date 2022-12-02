package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/oneoff"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITxNewRepo interface {
	GetTransferTx(chain string, height, limit int64) ([]*oneoff.TxNew, error)
}

var _ ITxNewRepo = new(TxNewRepo)

type TxNewRepo struct {
}

func (repo *TxNewRepo) GetTransferTx(chain string, height, limit int64) ([]*oneoff.TxNew, error) {
	var res []*oneoff.TxNew
	query := bson.M{
		"chain": chain,
		"height": bson.M{
			"$gt": height,
		},
	}
	err := repo.coll().Find(context.Background(), query).Sort("height").Limit(limit).All(&res)
	return res, err
}

func (repo *TxNewRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(oneoff.TxNew{}.CollectionName())
}
