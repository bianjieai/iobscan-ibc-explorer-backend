package repository

import (
	"context"
	"fmt"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"
)

type IRelayerFeeStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InsertMany(batch []*entity.IBCRelayerFeeStatistics) error
	InsertManyToNew(batch []*entity.IBCRelayerFeeStatistics) error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerFeeStatistics) error
}

var _ IRelayerFeeStatisticsRepo = new(RelayerFeeStatisticsRepo)

type RelayerFeeStatisticsRepo struct {
}

func (repo *RelayerFeeStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerFeeStatisticsCollName)
}

func (repo *RelayerFeeStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerFeeStatisticsNewCollName)
}

func (repo *RelayerFeeStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true).SetName("relayer_statistics_unique")
	key := []string{"relayer_address", "tx_type", "tx_status", "fee_denom", "segment_start_time", "segment_end_time"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})
}

func (repo *RelayerFeeStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCRelayerFeeStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCRelayerFeeStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}

func (repo *RelayerFeeStatisticsRepo) InsertMany(batch []*entity.IBCRelayerFeeStatistics) error {
	if _, err := repo.coll().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *RelayerFeeStatisticsRepo) InsertManyToNew(batch []*entity.IBCRelayerFeeStatistics) error {
	if _, err := repo.collNew().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *RelayerFeeStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerFeeStatistics) error {
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
