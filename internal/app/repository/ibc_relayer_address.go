package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRelayerAddressRepo interface {
	InsertOne(addr *entity.IBCRelayerAddress) error
	InsertMany(batch []*entity.IBCRelayerAddress) error
}

var _ IRelayerAddressRepo = new(RelayerAddressRepo)

type RelayerAddressRepo struct {
}

func (repo *RelayerAddressRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerAddressCollName)
}

func (repo *RelayerAddressRepo) InsertOne(addr *entity.IBCRelayerAddress) error {
	_, err := repo.coll().InsertOne(context.Background(), addr)
	return err
}

func (repo *RelayerAddressRepo) InsertMany(batch []*entity.IBCRelayerAddress) error {
	_, err := repo.coll().InsertMany(context.Background(), batch, insertIgnoreErrOpt)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}
	return err
}
