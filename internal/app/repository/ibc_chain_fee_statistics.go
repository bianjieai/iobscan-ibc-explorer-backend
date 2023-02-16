package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainFeeStatisticsRepo interface {
	ChainFeeStatistics(chain string, startTime, endTime int64) ([]*dto.DenomAmountStatisticsDTO, error)
}

var _ IChainFeeStatisticsRepo = new(ChainFeeStatisticsRepo)

type ChainFeeStatisticsRepo struct {
}

func (repo *ChainFeeStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainFeeStatisticsCollName)
}

func (repo *ChainFeeStatisticsRepo) ChainFeeStatistics(chain string, startTime, endTime int64) ([]*dto.DenomAmountStatisticsDTO, error) {
	match := bson.M{
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

	group := bson.M{
		operator.Group: bson.M{
			"_id": bson.M{
				"fee_denom": "$fee_denom",
			},
			"fee_amount": bson.M{
				operator.Sum: "$fee_amount",
			},
		},
	}

	project := bson.M{
		operator.Project: bson.M{
			"fee_denom":  "$_id.fee_denom",
			"fee_amount": "$fee_amount",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.DenomAmountStatisticsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
