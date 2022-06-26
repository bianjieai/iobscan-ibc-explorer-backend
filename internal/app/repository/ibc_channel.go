package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChannelRepo interface {
	UpdateOne(chainA, chainB, channelA, channelB string, updateTime, relayerCnt int64) error
}

var _ IChannelRepo = new(ChannelRepo)

type ChannelRepo struct {
}

func (repo *ChannelRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannel{}.CollectionName())
}

func (repo *ChannelRepo) UpdateOne(chainA, chainB, channelA, channelB string, updateTime, relayerCnt int64) error {
	filter := bson.M{
		"chain_a":   chainA,
		"chain_b":   chainB,
		"channel_a": channelA,
		"channel_b": channelB,
	}
	return repo.coll().UpdateOne(context.Background(), filter, bson.M{
		"$set": bson.M{
			"relayers":          relayerCnt,
			"channel_update_at": updateTime,
		},
	})
}
