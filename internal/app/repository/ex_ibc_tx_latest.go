package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IExIbcTxRepo interface {
	FindAll(skip, limit int64) ([]*entity.ExIbcTx, error)
	FindAllHistory(skip, limit int64) ([]*entity.ExIbcTx, error)
	CountBaseDenomTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error)
	CountBaseDenomHistoryTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error)
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

func (repo *ExIbcTxRepo) FindAll(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.coll().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) FindAllHistory(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.collHistory().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) countBaseDenomTransferTxsPipe() []bson.M {
	match := bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": entity.IbcTxUsefulStatus,
			},
		},
	}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":  "$base_denom",
				"sc_chain_id": "$sc_chain_id",
				"dc_chain_id": "$dc_chain_id",
			},
			"count": bson.M{
				"$sum": 1,
			},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"_id":         0,
			"base_denom":  "$_id.base_denom",
			"sc_chain_id": "$_id.sc_chain_id",
			"dc_chain_id": "$_id.dc_chain_id",
			"count":       "$count",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	return pipe
}

func (repo *ExIbcTxRepo) CountBaseDenomTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error) {
	pipe := repo.countBaseDenomTransferTxsPipe()
	var res []*dto.CountBaseDenomTransferAmountDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountBaseDenomHistoryTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error) {
	pipe := repo.countBaseDenomTransferTxsPipe()
	var res []*dto.CountBaseDenomTransferAmountDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
