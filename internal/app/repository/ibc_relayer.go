package repository

import (
	"context"
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
