package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ChainFieldChain            = "chain"
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
	UpdateIbcTokenValue(chain string, tokens int64, tokenValue string) error
	UpdateTransferTxs(chain string, txs int64, txsValue string) error
	UpdateRelayers(chain string, relayers int64) error
	FindAll(skip, limit int64) ([]*entity.IBCChain, error)
	Count() (int64, error)
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

func (repo *IbcChainRepo) UpdateRelayers(chain string, relayers int64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChain: chain},
		bson.M{
			"$set": bson.M{
				ChainFieldRelayers: relayers,
			},
		})
}

func (repo *IbcChainRepo) InserOrUpdate(chain entity.IBCChain) error {
	var res *entity.IBCChain
	err := repo.coll().Find(context.Background(), bson.M{ChainFieldChain: chain.Chain}).
		Select(bson.M{ChainFieldChain: 1}).One(&res)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			if _, err := repo.coll().InsertOne(context.Background(), chain); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChain: res.Chain},
		bson.M{
			"$set": bson.M{
				ChainFieldChannels:        chain.Channels,
				ChainFieldConnectedChains: chain.ConnectedChains,
				ChainFieldUpdateAt:        time.Now().Unix(),
			},
		})
}

func (repo *IbcChainRepo) UpdateIbcTokenValue(chain string, tokens int64, tokenValue string) error {
	updateData := bson.M{
		ChainFieldIbcTokens:      tokens,
		ChainFieldUpdateAt:       time.Now().Unix(),
		ChainFieldIbcTokensValue: tokenValue,
	}

	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChain: chain},
		bson.M{
			"$set": updateData,
		})
}

func (repo *IbcChainRepo) UpdateTransferTxs(chain string, txs int64, txsValue string) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{ChainFieldChain: chain},
		bson.M{
			"$set": bson.M{
				ChainFieldTransferTxs:      txs,
				ChainFieldTransferTxsValue: txsValue,
				ChainFieldUpdateAt:         time.Now().Unix(),
			},
		})
}

func (repo *IbcChainRepo) Count() (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{}).Count()
}
