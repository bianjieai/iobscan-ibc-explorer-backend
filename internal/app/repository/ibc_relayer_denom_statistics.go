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

type IRelayerDenomStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InsertMany(batch []*entity.IBCRelayerDenomStatistics) error
	InsertManyToNew(batch []*entity.IBCRelayerDenomStatistics) error
	BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerDenomStatistics) error
}

var _ IRelayerDenomStatisticsRepo = new(RelayerDenomStatisticsRepo)

type RelayerDenomStatisticsRepo struct {
}

func (repo *RelayerDenomStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerDenomStatisticsCollName)
}

func (repo *RelayerDenomStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerDenomStatisticsNewCollName)
}

func (repo *RelayerDenomStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true).SetName("relayer_statistics_unique")
	key := []string{"relayer_address", "tx_type", "tx_status", "base_denom", "base_denom_chain_id", "segment_start_time", "segment_end_time"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})

}
func (repo *RelayerDenomStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCRelayerDenomStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCRelayerDenomStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}
func (repo *RelayerDenomStatisticsRepo) InsertMany(batch []*entity.IBCRelayerDenomStatistics) error {
	if _, err := repo.coll().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *RelayerDenomStatisticsRepo) InsertManyToNew(batch []*entity.IBCRelayerDenomStatistics) error {
	if _, err := repo.collNew().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *RelayerDenomStatisticsRepo) BatchSwap(segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerDenomStatistics) error {
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