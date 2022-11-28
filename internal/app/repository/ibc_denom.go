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
	FindAll() (entity.IBCDenomList, error)
	FindBaseDenom() (entity.IBCDenomList, error)
	FindByBaseDenom(baseDenom, baseDenomChain string) (entity.IBCDenomList, error)
	FindByChain(chain string) (entity.IBCDenomList, error)
	FindByDenom(denom string) (entity.IBCDenomList, error)
	FindByDenomChain(denom, chain string) (*entity.IBCDenom, error)
	GetDenomGroupByChain() ([]*dto.GetDenomGroupByChainDTO, error)
	FindNoSymbolDenoms() (entity.IBCDenomList, error)
	FindSymbolDenoms() (entity.IBCDenomList, error)
	GetBaseDenomNoSymbol() ([]*dto.GetBaseDenomFromIbcDenomDTO, error)
	Count() (int64, error)
	BasedDenomCount() (int64, error)
	LatestCreateAt() (int64, error)
	UpdateSymbol(chain, denom, symbol string) error
	Insert(denom *entity.IBCDenom) error
	InsertBatch(denoms entity.IBCDenomList) error
	InsertBatchToNew(denoms entity.IBCDenomList) error
	UpdateDenom(denom *entity.IBCDenom) error
}

var _ IDenomRepo = new(DenomRepo)

type DenomRepo struct {
}

func (repo *DenomRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenom{}.CollectionName(false))
}

func (repo *DenomRepo) collNew() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenom{}.CollectionName(true))
}

func (repo *DenomRepo) FindAll() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindBaseDenom() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"is_base_denom": true}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindByBaseDenom(baseDenom, baseDenomChain string) (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"base_denom": baseDenom, "base_denom_chain": baseDenomChain}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindByChain(chain string) (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"chain": chain}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindByDenom(denom string) (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"denom": denom}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindByDenomChain(denom, chain string) (*entity.IBCDenom, error) {
	var res *entity.IBCDenom
	err := repo.coll().Find(context.Background(), bson.M{"denom": denom, "chain": chain}).One(&res)
	return res, err
}

func (repo *DenomRepo) GetDenomGroupByChain() ([]*dto.GetDenomGroupByChainDTO, error) {
	group := bson.M{
		"$group": bson.M{
			"_id": "$chain",
			"denom": bson.M{
				"$addToSet": "$denom",
			},
		},
	}

	var pipe []bson.M
	pipe = append(pipe, group)
	var res []*dto.GetDenomGroupByChainDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *DenomRepo) FindNoSymbolDenoms() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"symbol": ""}).All(&res)
	return res, err
}

func (repo *DenomRepo) FindSymbolDenoms() (entity.IBCDenomList, error) {
	var res entity.IBCDenomList
	err := repo.coll().Find(context.Background(), bson.M{"symbol": bson.M{"$ne": ""}}).All(&res)
	return res, err
}

func (repo *DenomRepo) UpdateSymbol(chain, denom, symbol string) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{"chain": chain, "denom": denom}, bson.M{
		"$set": bson.M{
			"symbol": symbol,
		}})
}

func (repo *DenomRepo) Insert(denom *entity.IBCDenom) error {
	_, err := repo.coll().InsertOne(context.Background(), denom)
	return err
}

func (repo *DenomRepo) InsertBatch(denoms entity.IBCDenomList) error {
	_, err := repo.coll().InsertMany(context.Background(), denoms, insertIgnoreErrOpt)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}

	return err
}

func (repo *DenomRepo) InsertBatchToNew(denoms entity.IBCDenomList) error {
	_, err := repo.collNew().InsertMany(context.Background(), denoms, insertIgnoreErrOpt)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	}

	return err
}

func (repo *DenomRepo) UpdateDenom(denom *entity.IBCDenom) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{"chain": denom.Chain, "denom": denom.Denom}, bson.M{
		"$set": bson.M{
			"base_denom":       denom.BaseDenom,
			"base_denom_chain": denom.BaseDenomChain,
			"prev_denom":       denom.PrevDenom,
			"prev_chain":       denom.PrevChain,
			"is_base_denom":    denom.IsBaseDenom,
		},
	})
}

func (repo *DenomRepo) LatestCreateAt() (int64, error) {
	var res entity.IBCDenom
	err := repo.coll().Find(context.Background(), bson.M{}).Sort("-create_at").One(&res)
	if err != nil {
		return 0, err
	}
	return res.CreateAt, nil
}

func (repo *DenomRepo) Count() (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{}).Count()
}

func (repo *DenomRepo) BasedDenomCount() (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{"is_base_denom": true}).Count()
}

func (repo *DenomRepo) GetBaseDenomNoSymbol() ([]*dto.GetBaseDenomFromIbcDenomDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"symbol": "",
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": "$base_denom",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group)
	var res []*dto.GetBaseDenomFromIbcDenomDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
