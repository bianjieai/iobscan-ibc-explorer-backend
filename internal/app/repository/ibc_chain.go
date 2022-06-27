package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ChainFieldChainId = "chain_id"
)

type IChainRepo interface {
	InserOrUpdate(chain entity.IBCChain) error
	UpdateIbcTokenValue(chainId string, tokens int64, tokenValue float64) error
	UpdateTransferTxs(chainId string, txs int64, txsValue string) error
	UpdateRelayers(chainId string, relayers int64) error
	FindAll() ([]*entity.IBCChain, error)
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

func (repo *IbcChainRepo) UpdateRelayers(chainId string, relayers int64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: chainId},
		bson.M{
			"$set": bson.M{
				"relayers": relayers,
			},
		})
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

func (repo *IbcChainRepo) UpdateIbcTokenValue(chainId string, tokens int64, tokenValue float64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: chainId},
		bson.M{
			"$set": bson.M{
				"ibc_tokens":       tokens,
				"ibc_tokens_value": tokenValue,
				"update_at":        time.Now().Unix(),
			},
		})
}

func (repo *IbcChainRepo) UpdateTransferTxs(chainId string, txs int64, txsValue string) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: chainId},
		bson.M{
			"$set": bson.M{
				"transfer_txs":       txs,
				"transfer_txs_value": txsValue,
				"update_at":          time.Now().Unix(),
			},
		})
}

func (repo *IbcChainRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-chain_id"},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})

	ensureIndexes(entity.IBCChain{}.CollectionName(), indexes)
}
