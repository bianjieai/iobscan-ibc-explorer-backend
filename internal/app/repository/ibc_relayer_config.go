package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

type IRelayerConfigRepo interface {
}

var _ IRelayerConfigRepo = new(RelayerConfigRepo)

type RelayerConfigRepo struct {
}

func (repo *RelayerConfigRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-relayer_name", "-relayer_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})
	indexes = append(indexes, options.IndexModel{
		Key: []string{"-relayer_id"},
	})

	ensureIndexes(entity.IBCRelayerConfig{}.CollectionName(), indexes)
}
func (repo *RelayerConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerConfig{}.CollectionName())
}

func (repo *RelayerConfigRepo) FindAll() ([]*entity.IBCRelayerConfig, error) {
	var res []*entity.IBCRelayerConfig
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}
