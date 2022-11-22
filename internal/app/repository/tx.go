package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITxRepo interface {
	GetTxByHash(chainId, txHash string) (entity.Tx, error)
	GetRelayerTxs(chainId string, skip, limit int64) ([]entity.Tx, error)
	GetActiveAccountsOfDay(chain string, startTime, endTime int64) ([]*dto.Aggr24hActiveAddrOfDayDto, error)
}

var _ ITxRepo = new(TxRepo)

type TxRepo struct {
}

func (repo *TxRepo) coll(chainId string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.Tx{}.CollectionName(chainId))
}

//========api support=========
func (repo *TxRepo) GetRelayerTxs(chain string, skip, limit int64) ([]entity.Tx, error) {
	var res []entity.Tx
	query := bson.M{
		"msgs.type": bson.M{
			"$in": []string{constant.MsgTypeRecvPacket, constant.MsgTypeTimeoutPacket, constant.MsgTypeAcknowledgement},
		},
	}
	err := repo.coll(chain).Find(context.Background(), query).
		Select(bson.M{"tx_hash": 1, "signers": 1, "fee": 1}).Skip(skip).Limit(limit).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (repo *TxRepo) GetTxByHash(chainId, txHash string) (entity.Tx, error) {
	var res entity.Tx
	query := bson.M{
		"tx_hash": txHash,
	}
	err := repo.coll(chainId).Find(context.Background(), query).One(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

//need index: time_-1_msgs.type_-1
func (repo *TxRepo) GetActiveAccountsOfDay(chainId string, startTime, endTime int64) ([]*dto.Aggr24hActiveAddrOfDayDto, error) {
	pipe := repo.AggrActiveAddrsOfDayPipe(startTime, endTime)
	var res []*dto.Aggr24hActiveAddrOfDayDto
	err := repo.coll(chainId).Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *TxRepo) AggrActiveAddrsOfDayPipe(startTime int64, endTime int64) []bson.M {
	match := bson.M{
		"$match": bson.M{
			"time": bson.M{
				"$gte": startTime,
				"$lt":  endTime,
			},
			"msgs.type": bson.M{
				"$in": []string{constant.MsgTypeTransfer, constant.MsgTypeRecvPacket, constant.MsgTypeTimeoutPacket, constant.MsgTypeAcknowledgement},
			},
		},
	}
	unwind := bson.M{
		"$unwind": "$addrs",
	}
	group := bson.M{
		"$group": bson.M{
			"_id": "$addrs",
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":     0,
			"address": "$_id",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, unwind, group, project)
	return pipe
}
