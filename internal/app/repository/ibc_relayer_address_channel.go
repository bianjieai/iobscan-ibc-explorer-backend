package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRelayerAddressChannelRepo interface {
	InsertOne(ac *entity.IBCRelayerAddressChannel) error
	InsertMany(batch []*entity.IBCRelayerAddressChannel) error
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
