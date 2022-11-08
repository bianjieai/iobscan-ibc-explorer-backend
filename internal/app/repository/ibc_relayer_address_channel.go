package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRelayerAddressChannelRepo interface {
	InsertOne(ac *entity.IBCRelayerAddressChannel) error
	InsertMany(batch []*entity.IBCRelayerAddressChannel) error
	FindChannels(arrs []string) ([]*entity.IBCRelayerAddressChannel, error)
}

var _ IRelayerAddressChannelRepo = new(RelayerAddressChannelRepo)

type RelayerAddressChannelRepo struct {
}

func (repo *RelayerAddressChannelRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerAddressChannelCollName)
}

func (repo *RelayerAddressChannelRepo) InsertOne(ac *entity.IBCRelayerAddressChannel) error {
	_, err := repo.coll().InsertOne(context.Background(), ac)
	return err
}

func (repo *RelayerAddressChannelRepo) InsertMany(batch []*entity.IBCRelayerAddressChannel) error {
	_, err := repo.coll().InsertMany(context.Background(), batch, insertIgnoreErrOpt)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}
	return err
}

func (repo *RelayerAddressChannelRepo) FindChannels(arrs []string) ([]*entity.IBCRelayerAddressChannel, error) {
	var res []*entity.IBCRelayerAddressChannel
	query := bson.M{
		"relayer_address": bson.M{"$in": arrs},
	}
	err := repo.coll().Find(context.Background(), query).All(&res)
	return res, err
}
