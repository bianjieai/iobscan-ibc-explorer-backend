package repository

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	officialOpts "go.mongodb.org/mongo-driver/mongo/options"
)

type IAddressStatisticsRepo interface {
	CreateNew() error
	SwitchColl() error
	InsertMany(batch []*entity.IBCAddressStatistics) error
	InsertManyToNew(batch []*entity.IBCAddressStatistics) error
	BatchSwap(chain string, segmentStartTime, segmentEndTime int64, batch []*entity.IBCAddressStatistics) error
	AddressStatistics(chain string, startTime, endTime int64) (*dto.ChainActiveAddressStatisticsDTO, error)
}

var _ IAddressStatisticsRepo = new(AddressStatisticsRepo)

type AddressStatisticsRepo struct {
}

func (repo *AddressStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCAddressStatisticsCollName)
}

func (repo *AddressStatisticsRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCAddressStatisticsNewCollName)
}

func (repo *AddressStatisticsRepo) CreateNew() error {
	ukOpts := officialOpts.Index().SetUnique(true).SetName("statistics_unique")
	uk := []string{"segment_start_time", "segment_end_time", "chain_name"}
	if err := repo.collNew().CreateOneIndex(context.Background(), opts.IndexModel{Key: uk, IndexOptions: ukOpts}); err != nil {
		return err
	}

	return nil
}

func (repo *AddressStatisticsRepo) SwitchColl() error {
	command := bson.D{{"renameCollection", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCAddressStatisticsNewCollName)},
		{"to", fmt.Sprintf("%s.%s", ibcDatabase, entity.IBCAddressStatisticsCollName)},
		{"dropTarget", true}}
	return mgo.Database(adminDatabase).RunCommand(context.Background(), command).Err()
}

func (repo *AddressStatisticsRepo) InsertMany(batch []*entity.IBCAddressStatistics) error {
	if _, err := repo.coll().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *AddressStatisticsRepo) InsertManyToNew(batch []*entity.IBCAddressStatistics) error {
	if _, err := repo.collNew().InsertMany(context.Background(), batch); err != nil {
		return err
	}
	return nil
}

func (repo *AddressStatisticsRepo) BatchSwap(chain string, segmentStartTime, segmentEndTime int64, batch []*entity.IBCAddressStatistics) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		query := bson.M{
			"segment_start_time": segmentStartTime,
			"segment_end_time":   segmentEndTime,
			"chain_name":         chain,
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

func (repo *AddressStatisticsRepo) AddressStatistics(chain string, startTime, endTime int64) (*dto.ChainActiveAddressStatisticsDTO, error) {
	match := bson.M{}
	if chain == "" {
		match = bson.M{
			operator.Match: bson.M{
				"segment_start_time": bson.M{
					operator.Gte: startTime,
				},
				"segment_end_time": bson.M{
					operator.Lte: endTime,
				},
			},
		}
	} else {
		match = bson.M{
			operator.Match: bson.M{
				"chain_name": chain,
				"segment_start_time": bson.M{
					operator.Gte: startTime,
				},
				"segment_end_time": bson.M{
					operator.Lte: endTime,
				},
			},
		}
	}

	group := bson.M{
		operator.Group: bson.M{
			"_id": "",
			"amount": bson.M{
				operator.Sum: "$active_address_num",
			},
		},
	}

	project := bson.M{
		operator.Project: bson.M{
			"address_amount": "$amount",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res dto.ChainActiveAddressStatisticsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).One(&res)
	return &res, err
}
