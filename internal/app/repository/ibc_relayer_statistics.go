package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"
)

// TODO remove this

type IRelayerStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InserOrUpdate(data entity.IBCRelayerStatistics) error
	CountRelayerBaseDenomAmt() ([]*dto.CountRelayerBaseDenomAmtDTO, error)
	Insert(relayerStatistics []entity.IBCRelayerStatistics) error
	InsertToNew(relayerStatistics []entity.IBCRelayerStatistics) error
	AggregateRelayerTxs() ([]*dto.AggRelayerTxsDTO, error)
	CreateStatisticId(scChain, dcChain, scChannel, dcChannel string) (string, string)
}

var _ IRelayerStatisticsRepo = new(RelayerStatisticsRepo)

type RelayerStatisticsRepo struct {
}

func (repo *RelayerStatisticsRepo) CreateStatisticId(scChain, dcChain, scChannel, dcChannel string) (string, string) {
	return fmt.Sprintf("%s|%s|%s|%s", scChain, scChannel, dcChain, dcChannel), fmt.Sprintf("%s|%s|%s|%s", dcChain, dcChannel, scChain, scChannel)
}

//func (repo *RelayerStatisticsRepo) EnsureIndexes() {
//	var indexes []options.IndexModel
//	indexes = append(indexes, options.IndexModel{
//		Key:          []string{"-transfer_base_denom", "-relayer_id", "-chain_id", "-channel"},
//		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
//	})
//	indexes = append(indexes, options.IndexModel{
//		Key: []string{"-relayer_id", "-chain_id"},
//	})
//
//	ensureIndexes(entity.IBCRelayerStatistics{}.CollectionName(), indexes)
//}

func (repo *RelayerStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerStatisticsCollName)
}

func (repo *RelayerStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerStatisticsNewCollName)
}

func (repo *RelayerStatisticsRepo) CreateNew() error {
	indexOpts := officialOpts.Index().SetUnique(true).SetName("relayer_statistics_unique")
	key := []string{"transfer_base_denom", "address", "statistic_id", "-segment_start_time", "-segment_end_time"}
	return repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts})
}

func (repo *RelayerStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCRelayerStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCRelayerStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}

func (repo *RelayerStatisticsRepo) Insert(relayerStatistics []entity.IBCRelayerStatistics) error {
	if _, err := repo.coll().InsertMany(context.Background(), relayerStatistics); err != nil {
		return err
	}
	return nil
}

func (repo *RelayerStatisticsRepo) InsertToNew(relayerStatistics []entity.IBCRelayerStatistics) error {
	if _, err := repo.collNew().InsertMany(context.Background(), relayerStatistics); err != nil {
		return err
	}
	return nil
}

func (repo *RelayerStatisticsRepo) InserOrUpdate(data entity.IBCRelayerStatistics) error {
	var res *entity.IBCRelayerStatistics
	filter := bson.M{
		"transfer_base_denom": data.TransferBaseDenom,
		"statistic_id":        data.StatisticId,
		"address":             data.Address,
		"segment_start_time":  data.SegmentStartTime,
		"segment_end_time":    data.SegmentEndTime,
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
				"transfer_amount":   data.TransferAmount,
				"success_total_txs": data.SuccessTotalTxs,
				"total_txs":         data.TotalTxs,
				"update_at":         time.Now().Unix(),
			},
		})
}

func (repo *RelayerStatisticsRepo) CountRelayerBaseDenomAmt() ([]*dto.CountRelayerBaseDenomAmtDTO, error) {
	match := bson.M{
		"$match": bson.M{},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"address":             "$address",
				"statistic_id":        "$statistic_id",
				"base_denom":          "$transfer_base_denom",
				"base_denom_chain_id": "$base_denom_chain_id",
			},
			"amount": bson.M{
				"$sum": bson.M{"$toDouble": "$transfer_amount"},
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":                 0,
			"address":             "$_id.address",
			"statistic_id":        "$_id.statistic_id",
			"base_denom":          "$_id.base_denom",
			"base_denom_chain_id": "$_id.base_denom_chain_id",
			"amount":              "$amount",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.CountRelayerBaseDenomAmtDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *RelayerStatisticsRepo) AggregateRelayerTxs() ([]*dto.AggRelayerTxsDTO, error) {
	match := bson.M{
		"$match": bson.M{},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"address":      "$address",
				"statistic_id": "$statistic_id",
			},
			"total_txs": bson.M{
				"$sum": "$total_txs",
			},
			"success_total_txs": bson.M{
				"$sum": "$success_total_txs",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":               0,
			"address":           "$_id.address",
			"statistic_id":      "$_id.statistic_id",
			"total_txs":         "$total_txs",
			"success_total_txs": "$success_total_txs",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.AggRelayerTxsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
