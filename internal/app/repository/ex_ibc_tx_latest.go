package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IExIbcTxRepo interface {
	GetMinTxTime(isTargetHistory bool) (int64, error)
	// FindFailLog filter by status = 2 or 4
	FindFailLog(startTime, endTime, skip, limit int64, isTargetHistory bool) ([]*entity.ExIbcTx, error)
}

var _ IExIbcTxRepo = new(ExIbcTxRepo)

type ExIbcTxRepo struct {
}

func (repo *ExIbcTxRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExIbcTx{}.CollectionName(false))
}

func (repo *ExIbcTxRepo) collHistory() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExIbcTx{}.CollectionName(true))
}

func (repo *ExIbcTxRepo) GetMinTxTime(isTargetHistory bool) (int64, error) {
	var res *entity.ExIbcTx
	var err error
	if isTargetHistory {
		err = repo.collHistory().Find(context.Background(), bson.M{}).Select(bson.M{"tx_time": 1}).Sort("tx_time").One(&res)
	} else {
		err = repo.coll().Find(context.Background(), bson.M{}).Select(bson.M{"tx_time": 1}).Sort("tx_time").One(&res)
	}

	if err != nil {
		return 0, err
	}
	return res.TxTime, nil
}

func (repo *ExIbcTxRepo) FindFailLog(startTime, endTime, skip, limit int64, isTargetHistory bool) ([]*entity.ExIbcTx, error) {
	query := bson.M{
		"status": bson.M{
			"$in": []entity.IbcTxStatus{entity.IbcTxStatusFailed, entity.IbcTxStatusRefunded},
		},
		"tx_time": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}
	selector := bson.M{
		"sc_chain":            1,
		"dc_chain":            1,
		"status":              1,
		"sc_tx_info.log":      1,
		"ack_timeout_tx_info": 1,
	}

	var res []*entity.ExIbcTx
	var err error
	if isTargetHistory {
		err = repo.collHistory().Find(context.Background(), query).Select(selector).Sort("tx_time").Skip(skip).Limit(limit).All(&res)
	} else {
		err = repo.coll().Find(context.Background(), query).Select(selector).Sort("tx_time").Skip(skip).Limit(limit).All(&res)
	}
	return res, err
}
