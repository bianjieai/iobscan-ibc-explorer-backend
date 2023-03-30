package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/qiniu/qmgo/operator"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IExIbcTxRepo interface {
	GetMinTxTime(isTargetHistory bool) (int64, error)
	// FindFailLog filter by status = 2 or 4
	FindFailLog(startTime, endTime, skip, limit int64, isTargetHistory bool) ([]*entity.ExIbcTx, error)
	AggrChainAddress(startTime, endTime int64, isTargetHistory bool, isTargetDcChain bool) ([]*dto.ChainActiveAddressesDTO, error)
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

func (repo *ExIbcTxRepo) AggrChainAddress(startTime, endTime int64, isTargetHistory bool, isTargetDcChain bool) ([]*dto.ChainActiveAddressesDTO, error) {
	match := bson.M{
		operator.Match: bson.M{
			"tx_time": bson.M{
				operator.Gte: startTime,
				operator.Lte: endTime,
			},
			"status": bson.M{
				operator.In: entity.IbcTxUsefulStatus,
			},
		},
	}

	group := bson.M{}
	if isTargetDcChain {
		group = bson.M{
			operator.Group: bson.M{
				"_id": "$dc_chain",
				"addresses": bson.M{
					operator.AddToSet: "$dc_addr",
				},
			},
		}
	} else {
		group = bson.M{
			operator.Group: bson.M{
				"_id": "$sc_chain",
				"addresses": bson.M{
					operator.AddToSet: "$sc_addr",
				},
			},
		}
	}

	project := bson.M{
		operator.Project: bson.M{
			"_id":       0,
			"chain":     "$_id",
			"addresses": "$addresses",
		},
	}

	var res []*dto.ChainActiveAddressesDTO
	var err error
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	if isTargetHistory {
		err = repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	} else {
		err = repo.coll().Aggregate(context.Background(), pipe).All(&res)
	}

	return res, err
}
