package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type IRelayerStatisticsRepo interface {
	InserOrUpdate(data entity.IBCRelayerStatistics) error
	CountRelayerTotalValue() ([]*dto.CountRelayerTotalValueDTO, error)
}

var _ IRelayerStatisticsRepo = new(RelayerStatisticsRepo)

type RelayerStatisticsRepo struct {
}

func (repo *RelayerStatisticsRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-transfer_base_denom", "-relayer_id", "-chain_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})
	indexes = append(indexes, options.IndexModel{
		Key: []string{"-relayer_id", "-chain_id"},
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

func (repo *RelayerStatisticsRepo) InserOrUpdate(data entity.IBCRelayerStatistics) error {
	var res *entity.IBCRelayerStatistics
	filter := bson.M{
		"transfer_base_denom": data.TransferBaseDenom,
		"relayer_id":          data.RelayerId,
		"chain_id":            data.ChainId,
	}
	err := repo.coll().Find(context.Background(), filter).One(&res)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			if _, err := repo.coll().InsertOne(context.Background(), data); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return repo.coll().UpdateOne(context.Background(), filter,
		bson.M{
			"$set": bson.M{
				"transfer_amount":      data.TransferAmount,
				"transfer_total_value": data.TransferTotalValue,
				"update_at":            time.Now().Unix(),
			},
		})
}

func (repo *RelayerStatisticsRepo) CountRelayerTotalValue() ([]*dto.CountRelayerTotalValueDTO, error) {
	match := bson.M{
		"$match": bson.M{},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"relayer_id": "$relayer_id",
				"chain_id":   "$chain_id",
			},
			"amount": bson.M{
				"$sum": bson.M{"$toDouble": "$transfer_total_value"},
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":        0,
			"relayer_id": "$_id.relayer_id",
			"chain_id":   "$_id.chain_id",
			"amount":     "$amount",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.CountRelayerTotalValueDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
