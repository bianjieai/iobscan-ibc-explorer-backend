package repository

import (
	"context"
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"
)

type IChannelStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChannelStatistics) error
	BatchInsert(batch []*entity.IBCChannelStatistics) error
	BatchInsertToNew(batch []*entity.IBCChannelStatistics) error
	Aggr() ([]*dto.ChannelStatisticsAggrDTO, error)
}

var _ IChannelStatisticsRepo = new(ChannelStatisticsRepo)

type ChannelStatisticsRepo struct {
}

func (repo *ChannelStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannelStatisticsCollName)
}

func (repo *ChannelStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChannelStatisticsNewCollName)
}

func (repo *ChannelStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true).SetName("channel_statistics_unique")
	key := []string{"channel_id", "base_denom", "base_denom_chain", "status", "-segment_start_time", "-segment_end_time"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})
}

func (repo *ChannelStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCChannelStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCChannelStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}

func (repo *ChannelStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCChannelStatistics) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"segment_start_time": segmentStartTime,
			"segment_end_time":   segmentEndTime,
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

func (repo *ChannelStatisticsRepo) BatchInsert(batch []*entity.IBCChannelStatistics) error {
	if len(batch) == 0 {
		return nil
	}

	_, err := repo.coll().InsertMany(context.Background(), batch)
	return err
}

func (repo *ChannelStatisticsRepo) BatchInsertToNew(batch []*entity.IBCChannelStatistics) error {
	if len(batch) == 0 {
		return nil
	}

	_, err := repo.collNew().InsertMany(context.Background(), batch)
	return err
}

func (repo *ChannelStatisticsRepo) Aggr() ([]*dto.ChannelStatisticsAggrDTO, error) {
	ibcTxUseStatus := []entity.IbcTxStatus{entity.IbcTxStatusSuccess, entity.IbcTxStatusProcessing, entity.IbcTxStatusRefunded}
	match := bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": ibcTxUseStatus,
			},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"channel_id":       "$channel_id",
				"base_denom":       "$base_denom",
				"base_denom_chain": "$base_denom_chain",
			},
			"count": bson.M{
				"$sum": "$transfer_txs",
			},
			"amount": bson.M{
				"$sum": bson.M{
					"$toDouble": "$transfer_amount",
				},
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":              0,
			"channel_id":       "$_id.channel_id",
			"base_denom":       "$_id.base_denom",
			"base_denom_chain": "$_id.base_denom_chain",
			"count":            "$count",
			"amount":           "$amount",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.ChannelStatisticsAggrDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
