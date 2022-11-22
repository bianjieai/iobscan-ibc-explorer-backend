package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IExIbcTxRepo interface {
	FindAllByStatus(stats []entity.IbcTxStatus, skip, limit int64) ([]*entity.ExIbcTx, error)
	FindAllHistoryByStatus(stats []entity.IbcTxStatus, skip, limit int64) ([]*entity.ExIbcTx, error)
	FindRelayerTxs(skip, limit int64) ([]*entity.ExIbcTx, error)
	FindHistoryRelayerTxs(skip, limit int64) ([]*entity.ExIbcTx, error)
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

//=========================api_support===============================================
func (repo *ExIbcTxRepo) FindAllByStatus(stats []entity.IbcTxStatus, skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.coll().Find(context.Background(), bson.M{"status": bson.M{
		"$in": stats,
	}}).Sort("-create_at").Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) FindAllHistoryByStatus(stats []entity.IbcTxStatus, skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.collHistory().Find(context.Background(), bson.M{"status": bson.M{
		"$in": stats,
	}}).Sort("-create_at").Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) FindRelayerTxs(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.coll().Find(context.Background(), bson.M{
		"status":            bson.M{"$in": []entity.IbcTxStatus{entity.IbcTxStatusSuccess, entity.IbcTxStatusRefunded, entity.IbcTxStatusFailed}},
		"sc_tx_info.status": entity.TxStatusSuccess,
	}).Sort("-create_at").Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) FindHistoryRelayerTxs(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.collHistory().Find(context.Background(), bson.M{
		"status":            bson.M{"$in": []entity.IbcTxStatus{entity.IbcTxStatusSuccess, entity.IbcTxStatusRefunded, entity.IbcTxStatusFailed}},
		"sc_tx_info.status": entity.TxStatusSuccess,
	}).Sort("-create_at").Skip(skip).Limit(limit).All(&res)
	return res, err
}
