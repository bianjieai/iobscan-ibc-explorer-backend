package repository

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	ChainFieldChainId          = "chain_id"
	ChainFieldIbcTokens        = "ibc_tokens"
	ChainFieldRelayers         = "relayers"
	ChainFieldChannels         = "channels"
	ChainFieldConnectedChains  = "connected_chains"
	ChainFieldIbcTokensValue   = "ibc_tokens_value"
	ChainFieldTransferTxs      = "transfer_txs"
	ChainFieldTransferTxsValue = "transfer_txs_value"
	ChainFieldUpdateAt         = "update_at"
)

type IChainRepo interface {
	InserOrUpdate(chain entity.IBCChain) error
	UpdateIbcTokenValue(chainId string, tokens int64, tokenValue float64) error
	UpdateTransferTxs(chainId string, txs int64, txsValue string) error
	UpdateRelayers(chainId string, relayers int64) error
	FindAll(skip, limit int64) ([]*entity.IBCChain, error)
}

var _ IChainRepo = new(IbcChainRepo)

type IbcChainRepo struct {
}

func (repo *IbcChainRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCChain{}.CollectionName())
}

func (repo *IbcChainRepo) FindAll(skip, limit int64) ([]*entity.IBCChain, error) {
	var res []*entity.IBCChain
	err := repo.coll().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).Sort("-" + ChainFieldIbcTokens).All(&res)
	return res, err
}

func (repo *IbcChainRepo) UpdateRelayers(chainId string, relayers int64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: chainId},
		bson.M{
			"$set": bson.M{
				ChainFieldRelayers: relayers,
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
				ChainFieldChannels:        chain.Channels,
				ChainFieldConnectedChains: chain.ConnectedChains,
				ChainFieldUpdateAt:        time.Now().Unix(),
			},
		})
}

func (repo *IbcChainRepo) UpdateIbcTokenValue(chainId string, tokens int64, tokenValue float64) error {
	updateData := bson.M{
		ChainFieldIbcTokens:      tokens,
		ChainFieldUpdateAt:       time.Now().Unix(),
		ChainFieldIbcTokensValue: "",
	}
	if tokenValue > 0 {
		updateData[ChainFieldIbcTokensValue] = fmt.Sprint(tokenValue)
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: chainId},
		bson.M{
			"$set": updateData,
		})
}

func (repo *IbcChainRepo) UpdateTransferTxs(chainId string, txs int64, txsValue string) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChainId: chainId},
		bson.M{
			"$set": bson.M{
				ChainFieldTransferTxs:      txs,
				ChainFieldTransferTxsValue: txsValue,
				ChainFieldUpdateAt:         time.Now().Unix(),
			},
		})
}

func (repo *IbcChainRepo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:          []string{"-" + ChainFieldChainId},
		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
	})

	ensureIndexes(entity.IBCChain{}.CollectionName(), indexes)
}
