package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	ChainFieldChainId = "chain_id"
)

type IChainRepo interface {
	InserOrUpdate(chain entity.IBCChain) error
}

var _ IChainRepo = new(IbcChainRepo)

type IbcChainRepo struct {
}

func (repo *IbcChainRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChain{}.CollectionName())
}

func (repo *IbcChainRepo) FindAll() ([]*entity.IBCChain, error) {
	var res []*entity.IBCChain
	err := repo.coll().Find(context.Background(), bson.M{}).All(&res)
	return res, err
}

func (repo *IbcChainRepo) InserOrUpdate(chain entity.IBCChain) error {
	var res *entity.IBCChain
	err := repo.coll().Find(context.Background(), bson.M{ChainFieldChainId: chain.ChainId}).
		Select(bson.M{ChainFieldChainId: 1}).One(&res)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			if _, err := repo.coll().InsertOne(context.Background(), chain); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: res.ChainId},
		bson.M{
			"$set": bson.M{
				"channels":         chain.Channels,
				"connected_chains": chain.ConnectedChains,
				"update_at":        time.Now().Unix(),
			},
		})
}
