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

type IRelayerDenomStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InsertMany(batch []*entity.IBCRelayerDenomStatistics) error
	InsertManyToNew(batch []*entity.IBCRelayerDenomStatistics) error
	BatchSwap(chain string, segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerDenomStatistics) error
	AggrRelayerBaseDenomAmtAndTxs(relayAddrs []string) ([]*dto.CountRelayerBaseDenomAmtDTO, error)
	AggrRelayerAmtAndTxsBySegment(relayAddrs []string, segmentStartTime, segmentEndTime int64) ([]*dto.CountRelayerBaseDenomAmtDTO, error)
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
	ukOpts := officialOpts.Index().SetUnique(true).SetName("statistics_unique")
	uk := []string{"relayer_address", "tx_type", "tx_status", "base_denom", "base_denom_chain_id", "segment_start_time", "segment_end_time"}
	if err := repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: uk, IndexOptions: ukOpts}); err != nil {
		return err
	}

	indexOpts := officialOpts.Index()
	key := []string{"statistics_chain", "segment_start_time", "segment_end_time"}
	if err := repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: key, IndexOptions: indexOpts}); err != nil {
		return err
	}

	return nil
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

func (repo *RelayerDenomStatisticsRepo) BatchSwap(chain string, segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerDenomStatistics) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"statistics_chain":   chain,
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

func (repo *RelayerDenomStatisticsRepo) AggrRelayerBaseDenomAmtAndTxs(relayAddrs []string) ([]*dto.CountRelayerBaseDenomAmtDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"relayer_address": bson.M{"$in": relayAddrs},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":          "$base_denom",
				"base_denom_chain_id": "$base_denom_chain_id",
				"tx_status":           "$tx_status",
			},
			"amount": bson.M{
				"$sum": bson.M{"$toDouble": "$relayed_amount"},
			},
			"relayed_txs": bson.M{
				"$sum": "$relayed_txs",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":                 0,
			"base_denom":          "$_id.base_denom",
			"base_denom_chain_id": "$_id.base_denom_chain_id",
			"tx_status":           "$_id.tx_status",
			"amount":              "$amount",
			"total_txs":           "$relayed_txs",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.CountRelayerBaseDenomAmtDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

/***
db.ibc_relayer_denom_statistics.createIndex({
    "relayer_address": 1,
    "segment_start_time": 1,
    "segment_end_time": 1,
}, {background: true});
*/
func (repo *RelayerDenomStatisticsRepo) AggrRelayerAmtAndTxsBySegment(relayAddrs []string, segmentStartTime, segmentEndTime int64) ([]*dto.CountRelayerBaseDenomAmtDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"relayer_address":    bson.M{"$in": relayAddrs},
			"segment_start_time": bson.M{"$gte": segmentStartTime},
			"segment_end_time":   bson.M{"$lte": segmentEndTime},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":          "$base_denom",
				"base_denom_chain_id": "$base_denom_chain_id",
			},
			"amount": bson.M{
				"$sum": bson.M{"$toDouble": "$relayed_amount"},
			},
			"relayed_txs": bson.M{
				"$sum": "$relayed_txs",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":                 0,
			"base_denom":          "$_id.base_denom",
			"base_denom_chain_id": "$_id.base_denom_chain_id",
			"amount":              "$amount",
			"total_txs":           "$relayed_txs",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.CountRelayerBaseDenomAmtDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
