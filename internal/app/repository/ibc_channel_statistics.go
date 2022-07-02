package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IChannelStatisticsRepo interface {
	BatchSwap(channelId string, batch []*entity.IBCChannelStatistics) error
}

var _ IChannelStatisticsRepo = new(ChannelStatisticsRepo)

type ChannelStatisticsRepo struct {
}

//func (repo *ChannelStatisticsRepo) EnsureIndexes() {
//	var indexes []options.IndexModel
//	indexes = append(indexes, options.IndexModel{
//		Key:          []string{"channel_id, transfer_base_denom"},
//		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
//	})
//
//	ensureIndexes(entity.IBCChannelStatistics{}.CollectionName(), indexes)
//}

func (repo *ChannelStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannelStatistics{}.CollectionName())
}

func (repo *ChannelStatisticsRepo) BatchSwap(channelId string, batch []*entity.IBCChannelStatistics) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"channel_id": channelId,
		}
		if _, err := repo.coll().RemoveAll(sessCtx, query); err != nil {
			return nil, err
		}

		if len(batch) == 0 {
			return nil, nil
		}

		if _, err := repo.coll().InsertMany(sessCtx, batch); err != nil {
			return nil, err
		}

		return nil, nil
	}
	_, err := mgo.DoTransaction(context.Background(), callback)
	return err
}
