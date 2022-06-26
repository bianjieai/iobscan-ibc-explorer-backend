package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITxRepo interface {
	GetRelayerScChainAddr(packetId, chainId string) ([]*dto.GetRelayerScChainAddreeDTO, error)
	GetTimePeriodByUpdateClient(chainId string) (int64, int64, error)
}

var _ ITxRepo = new(TxRepo)

type TxRepo struct {
}

func (repo *TxRepo) coll(chainId string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.Tx{}.CollectionName(chainId))
}

func (repo *TxRepo) GetRelayerScChainAddr(packetId, chainId string) ([]*dto.GetRelayerScChainAddreeDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"msgs.msg.packet_id": packetId,
			"msgs.type": bson.M{
				"$in": []string{"acknowledge_packet", "timeout_packet"},
			},
		},
	}
	unwind := bson.M{
		"$unwind": "$msgs",
	}
	group := bson.M{
		"$group": bson.M{
			"_id": "$msgs.msg.signer",
		},
	}
	project := bson.M{
		"$project": bson.M{
			"sc_chain_address": "$_id",
		},
	}
	//sort for latest relayer address on sc_chain
	sort := bson.M{
		"$sort": bson.M{
			"height": -1,
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, unwind, group, project, sort)
	var res []*dto.GetRelayerScChainAddreeDTO
	err := repo.coll(chainId).Aggregate(context.Background(), pipe).All(&res)
	return res, err

}

// return value description
//1: latest update_client tx_time
//2: time_period
//3: error
func (repo *TxRepo) GetTimePeriodByUpdateClient(chainId string) (int64, int64, error) {
	var res []*entity.Tx
	err := repo.coll(chainId).Find(context.Background(),
		bson.M{"msgs.type": "update_client"}).
		Select(bson.M{"time": 1}).Sort("-height").Limit(2).All(&res)
	if err != nil {
		return 0, 0, err
	}
	if len(res) == 2 {
		return res[0].Time, res[0].Time - res[1].Time, nil
	}
	if len(res) == 1 {
		return res[0].Time, -1, nil
	}
	return 0, -1, nil
}