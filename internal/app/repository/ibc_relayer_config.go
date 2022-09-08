package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IRelayerConfigRepo interface {
	FindAll() ([]*entity.IBCRelayerConfig, error)
	FindRelayerPairIds() ([]*dto.RelayerPairIdDTO, error)
	Insert(cfg *entity.IBCRelayerConfig) error
}

var _ IRelayerConfigRepo = new(RelayerConfigRepo)

type RelayerConfigRepo struct {
}

//func (repo *RelayerConfigRepo) EnsureIndexes() {
//	var indexes []options.IndexModel
//	indexes = append(indexes, options.IndexModel{
//		Key:          []string{"-relayer_channel_pair"},
//		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
//	})
//
//	ensureIndexes(entity.IBCRelayerConfig{}.CollectionName(), indexes)
//}
func (repo *RelayerConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerConfig{}.CollectionName())
}

func (repo *RelayerConfigRepo) FindAll() ([]*entity.IBCRelayerConfig, error) {
	var res []*entity.IBCRelayerConfig
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *RelayerConfigRepo) FindRelayerPairIds() ([]*dto.RelayerPairIdDTO, error) {
	var res []*dto.RelayerPairIdDTO
	err := repo.coll().Find(context.Background(), bson.M{}).Select(bson.M{"relayer_pair_id": 1}).All(&res)
	return res, err
}

func (repo *RelayerConfigRepo) Insert(cfg *entity.IBCRelayerConfig) error {
	_, err := repo.coll().InsertOne(context.Background(), cfg)
	return err
}
