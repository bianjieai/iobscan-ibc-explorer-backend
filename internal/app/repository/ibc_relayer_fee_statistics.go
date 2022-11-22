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

type IRelayerFeeStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InsertMany(batch []*entity.IBCRelayerFeeStatistics) error
	InsertManyToNew(batch []*entity.IBCRelayerFeeStatistics) error
	BatchSwap(chain string, segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerFeeStatistics) error
	AggrRelayerFeeDenomAmt(relayAddrs []string) ([]*dto.AggrRelayerTxsAmtDTo, error)
	AggrChainAddressPair() ([]*dto.AggrChainAddrDTO, error)
	UpdateChainAddressComb(chain, address, chainAddressComb string) error
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
	ukOpts := officialOpts.Index().SetUnique(true).SetName("statistics_unique")
	uk := []string{"relayer_address", "tx_type", "tx_status", "fee_denom", "segment_start_time", "segment_end_time"}
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

func (repo *RelayerFeeStatisticsRepo) BatchSwap(chain string, segmentStartTime, segmentEndTime int64, batch []*entity.IBCRelayerFeeStatistics) error {
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

func (repo *RelayerFeeStatisticsRepo) AggrRelayerFeeDenomAmt(relayAddrs []string) ([]*dto.AggrRelayerTxsAmtDTo, error) {
	match := bson.M{
		"$match": bson.M{
			"relayer_address": bson.M{"$in": relayAddrs},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"fee_denom":        "$fee_denom",
				"statistics_chain": "$statistics_chain",
			},
			"amount": bson.M{
				"$sum": "$fee_amount",
			},
			"total_txs": bson.M{
				"$sum": "$relayed_txs",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":       0,
			"fee_denom": "$_id.fee_denom",
			"chain_id":  "$_id.statistics_chain",
			"amount":    "$amount",
			"total_txs": "$total_txs",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.AggrRelayerTxsAmtDTo
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *RelayerFeeStatisticsRepo) AggrChainAddressPair() ([]*dto.AggrChainAddrDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"chain":   "$statistics_chain",
				"address": "$relayer_address",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":     0,
			"chain":   "$_id.chain",
			"address": "$_id.address",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, group, project)
	var res []*dto.AggrChainAddrDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *RelayerFeeStatisticsRepo) UpdateChainAddressComb(chain, address, chainAddressComb string) error {
	query := bson.M{
		"statistics_chain": chain, "relayer_address": address,
	}
	set := bson.M{
		"$set": bson.M{
			"chain_address_comb": chainAddressComb,
		},
	}
	_, err := repo.coll().UpdateAll(context.Background(), query, set)
	return err
}
