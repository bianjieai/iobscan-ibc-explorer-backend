package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChannelConfigRepo interface {
	Find(chainA, channelA, chainB, channelB string) (*entity.IBCChannelConfig, error)
}

var _ IChannelConfigRepo = new(ChannelConfigRepo)

type ChannelConfigRepo struct {
}

func (repo *ChannelConfigRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannelConfig{}.CollectionName())
}

func (repo *ChannelConfigRepo) Find(chainA, channelA, chainB, channelB string) (*entity.IBCChannelConfig, error) {
	var res entity.IBCChannelConfig
	query := bson.M{
		"$or": []bson.M{
			{
				"chain_a":   chainA,
				"channel_a": channelA,
				"chain_b":   chainB,
				"channel_b": channelB,
			},
			{
				"chain_a":   chainB,
				"channel_a": channelB,
				"chain_b":   chainA,
				"channel_b": channelA,
			},
		},
	}

	err := repo.coll().Find(context.Background(), query).One(&res)
	return &res, err
}
