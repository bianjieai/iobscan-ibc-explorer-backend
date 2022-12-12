package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRelayerAddressChannelRepo interface {
	InsertOne(ac *entity.IBCRelayerAddressChannel) error
	InsertMany(batch []*entity.IBCRelayerAddressChannel) error
	FindChannels(arrs []string) ([]*entity.IBCRelayerAddressChannel, error)
	FindByAddressChain(address, chain string) ([]*entity.IBCRelayerAddressChannel, error)
	DistinctAddr() ([]*dto.ChainAddressDTO, error)
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

func (repo *RelayerAddressChannelRepo) FindByAddressChain(address, chain string) ([]*entity.IBCRelayerAddressChannel, error) {
	var res []*entity.IBCRelayerAddressChannel
	query := bson.M{
		"relayer_address": address,
		"chain":           chain,
	}
	err := repo.coll().Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *RelayerAddressChannelRepo) DistinctAddr() ([]*dto.ChainAddressDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"chain":   "$chain",
				"address": "$relayer_address",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":     0,
			"chain":   "$_id.chain",
			"address": "$_id.address",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group, project)
	var res []*dto.ChainAddressDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
