package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	RelayerFieldelayerId     = "relayer_id"
	RelayerFieldLatestTxTime = "latest_tx_time"
)

type IRelayerRepo interface {
	FindLatestOne() (*entity.IBCRelayer, error)
	Insert(relayer []entity.IBCRelayer) error
	Update(relayerId string, data bson.M) error
	FindAll(skip, limit int64) ([]*entity.IBCRelayer, error)
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

func (repo *IbcRelayerRepo) Insert(relayer []entity.IBCRelayer) error {
	if _, err := repo.coll().InsertMany(context.Background(), relayer); err != nil {
		return err
	}
	return nil
}

func (repo *IbcRelayerRepo) Update(relayerId string, data bson.M) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, data)
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
