package repository

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IRelayerConfigRepo interface {
	FindAll() ([]*entity.IBCRelayerConfig, error)
}

var _ IRelayerConfigRepo = new(RelayerConfigRepo)

type RelayerConfigRepo struct {
}

func CreateRelayerChannelPair(chainA, chainB, channelA, channelB, addressA, AddressB string) (string, string) {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", chainA, channelA, addressA, chainB, channelB, AddressB),
		fmt.Sprintf("%s:%s:%s:%s:%s:%s", chainB, channelB, AddressB, chainA, channelA, addressA)
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
