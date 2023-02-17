package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type IChainInflowStatisticsRepo interface {
	InflowStatistics(chain string, startTime, endTime int64) ([]*dto.FlowStatisticsDTO, error)
}

var _ IChainInflowStatisticsRepo = new(ChainInflowStatisticsRepo)

type ChainInflowStatisticsRepo struct {
}

func (repo *ChainInflowStatisticsRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChainInflowStatisticsCollName)
}

func (repo *ChainInflowStatisticsRepo) InflowStatistics(chain string, startTime, endTime int64) ([]*dto.FlowStatisticsDTO, error) {
	match := bson.M{
		operator.Match: bson.M{
			"status": entity.IbcTxStatusSuccess,
			"chain":  chain,
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
				"base_denom":       "$base_denom",
				"base_denom_chain": "$base_denom_chain",
			},
			"denom_amount": bson.M{
				operator.Sum: "$denom_amount",
			},
			"txs_count": bson.M{
				operator.Sum: "$txs_number",
			},
		},
	}

	project := bson.M{
		operator.Project: bson.M{
			"_id":              0,
			"base_denom":       "$_id.base_denom",
			"base_denom_chain": "$_id.base_denom_chain",
			"denom_amount":     "$denom_amount",
			"txs_count":        "$txs_count",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.FlowStatisticsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
