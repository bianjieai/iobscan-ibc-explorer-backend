package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IDenomRepo interface {
	FindBaseDenom() (entity.IBCDenomList, error)
	FindByBaseDenom(baseDenom string) (entity.IBCDenomList, error)
	GetDenomGroupByBaseDenom() ([]*dto.GetDenomGroupByBaseDenomDTO, error)
	FindTokenOthers() (entity.IBCDenomList, error)
}

var _ IDenomRepo = new(DenomRepo)

type DenomRepo struct {
}

func (repo *DenomRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenom{}.CollectionName())
}

func (repo *DenomRepo) FindBaseDenom() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"is_base_denom": true, "is_source_chain": true}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindByBaseDenom(baseDenom string) (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"base_denom": baseDenom}).All(&res)
	return res, err
}

func (repo *DenomRepo) GetDenomGroupByBaseDenom() ([]*dto.GetDenomGroupByBaseDenomDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": "$base_denom",
			"denom": bson.M{
				"$addToSet": "$denom",
			},
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group)
	var res []*dto.GetDenomGroupByBaseDenomDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *DenomRepo) FindTokenOthers() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"symbol": "", "is_base_denom": true}).All(&res)
	return res, err
}
