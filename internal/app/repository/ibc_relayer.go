package repository

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	RelayerFieldelayerId     = "relayer_id"
	RelayerFieldLatestTxTime = "latest_tx_time"
)

type IRelayerRepo interface {
	FindLatestOne() (*entity.IBCRelayer, error)
	Insert(relayer []entity.IBCRelayer) error
	UpdateStatusAndTime(relayerId string, status int, updateTime, timePeriod int64) error
	UpdateTxsInfo(relayerId string, txs, txsSuccess int64, totalValue float64) error
	FindAll(skip, limit int64) ([]*entity.IBCRelayer, error)
	FindAllBycond(chainId string, status int, skip, limit int64, useCount bool) ([]*entity.IBCRelayer, int64, error)
	FindRelayersCnt(chainId string) (int64, error)
	CountChannelRelayers() ([]*dto.CountChannelRelayersDTO, error)
	FindRelayerId(chainId string, relayerAddr string) (*entity.IBCRelayer, error)
}

var _ IRelayerRepo = new(IbcRelayerRepo)

type IbcRelayerRepo struct {
}

func (repo *IbcRelayerRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-chain_a", "-channel_a", "-chain_a_address"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-chain_b", "-channel_b", "-chain_b_address"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})
	indexes = append(indexes, options.IndexModel{
		Key: []string{"-chain_b_address", "-chain_b"},
	})
	indexes = append(indexes, options.IndexModel{
		Key: []string{"-chain_a_address", "-chain_a"},
	})

	ensureIndexes(entity.IBCRelayer{}.CollectionName(), indexes)
}

func (repo *IbcRelayerRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayer{}.CollectionName())
}

func (repo *IbcRelayerRepo) FindAll(skip, limit int64) ([]*entity.IBCRelayer, error) {
	var res []*entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).Sort("+update_time").All(&res)
	return res, err
}

func (repo *IbcRelayerRepo) FindAllBycond(chainId string, status int, skip, limit int64, useCount bool) ([]*entity.IBCRelayer, int64, error) {
	var (
		res   []*entity.IBCRelayer
		total int64
	)
	filter := bson.M{}
	if chainId != "" {
		chains := strings.Split(chainId, ",")
		filter["$or"] = []bson.M{
			{"chain_a": bson.M{"$in": chains}},
			{"chain_b": bson.M{"$in": chains}},
		}
	}
	if status > 0 {
		filter["status"] = status
	}
	err := repo.coll().Find(context.Background(), filter).Skip(skip).Limit(limit).Sort("+update_time").All(&res)
	if useCount {
		total, err = repo.coll().Find(context.Background(), filter).Count()
	}
	return res, total, err
}

func (repo *IbcRelayerRepo) Insert(relayer []entity.IBCRelayer) error {
	if _, err := repo.coll().InsertMany(context.Background(), relayer); err != nil {
		return err
	}
	return nil
}

func (repo *IbcRelayerRepo) UpdateTxsInfo(relayerId string, txs, txsSuccess int64, totalValue float64) error {
	updateData := bson.M{
		"transfer_total_txs":       txs,
		"transfer_success_txs":     txsSuccess,
		"transfer_total_txs_value": "",
		"update_at":                time.Now().Unix(),
	}
	if totalValue > 0 {
		updateData["transfer_total_txs_value"] = fmt.Sprint(totalValue)
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": updateData})
}
func (repo *IbcRelayerRepo) UpdateStatusAndTime(relayerId string, status int, updateTime, timePeriod int64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": bson.M{
			"status":      status,
			"update_time": updateTime,
			"time_period": timePeriod,
			"update_at":   time.Now().Unix(),
		}})
}

func (repo *IbcRelayerRepo) FindLatestOne() (*entity.IBCRelayer, error) {
	var res *entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{}).Sort("-latest_tx_time").One(&res)
	return res, err
}

func (repo *IbcRelayerRepo) FindRelayersCnt(chainId string) (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{
		"$or": []bson.M{
			{"chain_a": chainId},
			{"chain_b": chainId},
		},
	}).Count()
}

func (repo *IbcRelayerRepo) FindRelayerId(chainId string, relayerAddr string) (*entity.IBCRelayer, error) {
	var res *entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{
		"$or": []bson.M{
			{"chain_a": chainId, "chain_a_address": relayerAddr},
			{"chain_b": chainId, "chain_b_address": relayerAddr},
		},
	}).Select(bson.M{RelayerFieldelayerId: 1}).One(&res)
	return res, err
}

func (repo *IbcRelayerRepo) CountChannelRelayers() ([]*dto.CountChannelRelayersDTO, error) {
	match := bson.M{
		"$match": bson.M{},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"chain_a":   "$chain_a",
				"chain_b":   "$chain_b",
				"channel_a": "$channel_a",
				"channel_b": "$channel_b",
			},
			"count": bson.M{
				"$sum": 1,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":       0,
			"chain_a":   "$_id.chain_a",
			"chain_b":   "$_id.chain_b",
			"channel_a": "$_id.channel_a",
			"channel_b": "$_id.channel_b",
			"count":     "$count",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.CountChannelRelayersDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
