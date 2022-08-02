package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDenomRepo interface {
	FindBaseDenom() (entity.IBCDenomList, error)
	FindByBaseDenom(baseDenom string) (entity.IBCDenomList, error)
	FindByChainId(chainId string) (entity.IBCDenomList, error)
	FindByDenom(denom string) (entity.IBCDenomList, error)
	GetDenomGroupByChainId() ([]*dto.GetDenomGroupByChainIdDTO, error)
	FindNoSymbolDenoms() (entity.IBCDenomList, error)
	UpdateSymbol(chainId, denom, symbol string) error
	InsertBatch(denoms entity.IBCDenomList) error
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

func (repo *DenomRepo) FindByChainId(chainId string) (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"chain_id": chainId}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindByDenom(denom string) (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"denom": denom}).All(&res)
	return res, err
}

func (repo *DenomRepo) GetDenomGroupByChainId() ([]*dto.GetDenomGroupByChainIdDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": "$chain_id",
			"denom": bson.M{
				"$addToSet": "$denom",
			},
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group)
	var res []*dto.GetDenomGroupByChainIdDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *DenomRepo) FindNoSymbolDenoms() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"symbol": ""}).All(&res)
	return res, err
}

func (repo *DenomRepo) UpdateSymbol(chainId, denom, symbol string) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{"chain_id": chainId, "denom": denom}, bson.M{
		"$set": bson.M{
			"symbol": symbol,
		}})
}

func (repo *DenomRepo) InsertBatch(denoms entity.IBCDenomList) error {
	_, err := repo.coll().InsertMany(context.Background(), denoms, insertIgnoreErrOpt)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}

	return err
}
