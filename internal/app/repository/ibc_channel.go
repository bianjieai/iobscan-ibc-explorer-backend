package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

type IChannelRepo interface {
	UpdateOne(filter interface{}, update interface{}) error
	FindAll() (entity.IBCChannelList, error)
	InsertBatch(batch []*entity.IBCChannel) error
	UpdateChannel(channel *entity.IBCChannel) error
}

var _ IChannelRepo = new(ChannelRepo)

type ChannelRepo struct {
}

func (repo *ChannelRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannel{}.CollectionName())
}

func (repo *ChannelRepo) UpdateOne(filter interface{}, update interface{}) error {
	return repo.coll().UpdateOne(context.Background(), filter, update)
}

func (repo *ChannelRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"channel_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})

	ensureIndexes(entity.IBCChannel{}.CollectionName(), indexes)
}

func (repo *ChannelRepo) FindAll() (entity.IBCChannelList, error) {
	var res entity.IBCChannelList
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *ChannelRepo) InsertBatch(batch []*entity.IBCChannel) error {
	if len(batch) == 0 {
		return nil
	}
	now := time.Now().Unix()
	for _, v := range batch {
		v.UpdateAt = now
		v.CreateAt = now
	}
	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *ChannelRepo) UpdateChannel(channel *entity.IBCChannel) error {
	query := bson.M{
		"channel_id": channel.ChannelId,
	}
	update := bson.M{
		"$set": bson.M{
			"status":                 channel.Status,
			"operating_period":       channel.OperatingPeriod,
			"latest_settlement_time": channel.LatestSettlementTime,
			"transfer_txs":           channel.TransferTxs,
			"transfer_txs_value":     channel.TransferTxsValue,
			"update_at":              time.Now().Unix(),
		},
	}
	return repo.coll().UpdateOne(context.Background(), query, update)
}
