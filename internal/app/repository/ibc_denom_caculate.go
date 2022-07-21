package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IDenomCalculateRepo interface {
	InsertTransaction(denomList []*entity.IBCDenomCalculate, baseDenom, scChainId, ibcInfoHashCalculate string) error
	FindByBaseDenom(baseDenom string) ([]*entity.IBCDenomCalculate, error)
	FindByScChainId(ScChainId string) ([]*entity.IBCDenomCalculate, error)
}

var _ IDenomCalculateRepo = new(DenomCalculateRepo)

type DenomCalculateRepo struct {
}

func (repo *DenomCalculateRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCDenomCalculate{}.CollectionName())
}

func (repo *DenomCalculateRepo) FindByBaseDenom(baseDenom string) ([]*entity.IBCDenomCalculate, error) {
	var res []*entity.IBCDenomCalculate
	query := bson.M{"base_denom": baseDenom}
	err := repo.coll().Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *DenomCalculateRepo) FindByScChainId(ScChainId string) ([]*entity.IBCDenomCalculate, error) {
	var res []*entity.IBCDenomCalculate
	query := bson.M{"sc_chain_id": ScChainId}
	err := repo.coll().Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *DenomCalculateRepo) InsertTransaction(denomList []*entity.IBCDenomCalculate, baseDenom, scChainId, ibcInfoHashCalculate string) error {
	callback := func(sessCtx context.Context) (interface{}, error) {
		if len(denomList) > 0 {
			if _, err := repo.coll().InsertMany(sessCtx, denomList); err != nil {
				return nil, err
			}
		}

		return nil, new(BaseDenomRepo).UpdateIbcInfoHashCalculate(baseDenom, scChainId, ibcInfoHashCalculate)
	}
	_, err := mgo.DoTransaction(context.Background(), callback)
	return err
}
