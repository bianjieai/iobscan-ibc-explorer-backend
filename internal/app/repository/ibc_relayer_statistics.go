package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

type IRelayerStatisticsRepo interface {
}

var _ IRelayerStatisticsRepo = new(RelayerStatisticsRepo)

type RelayerStatisticsRepo struct {
}

func (repo *RelayerStatisticsRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-transfer_base_denom", "-relayer_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})
	indexes = append(indexes, options.IndexModel{
		Key: []string{"-relayer_id"},
	})

	ensureIndexes(entity.IBCRelayerStatistics{}.CollectionName(), indexes)
}

func (repo *RelayerStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerStatistics{}.CollectionName())
}

func (repo *RelayerStatisticsRepo) Insert(relayerStatistics []entity.IBCRelayerStatistics) error {
	if _, err := repo.coll().InsertMany(context.Background(), relayerStatistics); err != nil {
		return err
	}
	return nil
}
