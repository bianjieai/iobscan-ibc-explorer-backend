package repository

import (
	"context"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	RelayerFieldelayerId              = "relayer_id"
	RelayerFieldLatestTxTime          = "latest_tx_time"
	RelayerFieldTransferTotalTxs      = "transfer_total_txs"
	RelayerFieldTransferSuccessTxs    = "transfer_success_txs"
	RelayerFieldTransferTotalTxsValue = "transfer_total_txs_value"
	RelayerFieldUpdateTime            = "update_time"
	RelayerFieldTimePeriod            = "time_period"
	RelayerFieldStatus                = "status"
	RelayerFieldChainA                = "chain_a"
	RelayerFieldChainB                = "chain_b"
	RelayerFieldChannelA              = "channel_a"
	RelayerFieldChannelB              = "channel_b"
	RelayerFieldChainAAddress         = "chain_a_address"
	RelayerFieldChainAALLAddress      = "chain_a_all_address"
	RelayerFieldChainBAddress         = "chain_b_address"
	RelayerFieldUpdateAt              = "update_at"
)

type IRelayerRepo interface {
	FindLatestOne() (*entity.IBCRelayer, error)
	Insert(relayer []entity.IBCRelayer) error
	UpdateStatusAndTime(relayerId string, status int, updateTime, timePeriod int64) error
	UpdateTxsInfo(relayerId string, txs, txsSuccess int64, totalValue string) error
	FindAll(skip, limit int64) ([]*entity.IBCRelayer, error)
	FindAllBycond(chainId string, status int, skip, limit int64, useCount bool) ([]*entity.IBCRelayer, int64, error)
	CountBycond(chainId string, status int) (int64, error)
	CountChainRelayers(chainId string) (int64, error)
	CountChannelRelayers() ([]*dto.CountChannelRelayersDTO, error)
	FindRelayer(chainId, relayerAddr, channel string) (*entity.IBCRelayer, error)
	FindEmptyAddrAll(skip, limit int64) ([]*entity.IBCRelayer, error)
	UpdateSrcAddress(relayerId string, addrs []string) error
}

var _ IRelayerRepo = new(IbcRelayerRepo)

type IbcRelayerRepo struct {
}

//func (repo *IbcRelayerRepo) EnsureIndexes() {
//	var indexes []options.IndexModel
//	indexes = append(indexes, options.IndexModel{
//		Key:          []string{"-" + RelayerFieldChainA, "-" + RelayerFieldChannelA, "-" + RelayerFieldChainAAddress},
//		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
//	})
//	indexes = append(indexes, options.IndexModel{
//		Key:          []string{"-" + RelayerFieldChainB, "-" + RelayerFieldChannelB, "-" + RelayerFieldChainBAddress},
//		IndexOptions: new(moptions.IndexOptions).SetUnique(true),
//	})
//	indexes = append(indexes, options.IndexModel{
//		Key: []string{"-" + RelayerFieldChainBAddress, "-" + RelayerFieldChainB},
//	})
//	indexes = append(indexes, options.IndexModel{
//		Key: []string{"-" + RelayerFieldChainAAddress, "-" + RelayerFieldChainA},
//	})
//
//	ensureIndexes(entity.IBCRelayer{}.CollectionName(), indexes)
//}

func (repo *IbcRelayerRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayer{}.CollectionName())
}

func (repo *IbcRelayerRepo) FindAll(skip, limit int64) ([]*entity.IBCRelayer, error) {
	var res []*entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *IbcRelayerRepo) FindEmptyAddrAll(skip, limit int64) ([]*entity.IBCRelayer, error) {
	var res []*entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{RelayerFieldChainAAddress: ""}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *IbcRelayerRepo) analyzeCond(chainId string, status int) bson.M {
	filter := bson.M{}
	if chainId != "" {
		chains := strings.Split(chainId, ",")
		if length := len(chains); length <= 2 {
			switch length {
			case 1:
				if !strings.Contains(chainId, constant.AllChain) {
					filter["$or"] = []bson.M{
						{RelayerFieldChainA: chains[0]},
						{RelayerFieldChainB: chains[0]},
					}
				}
				break
			case 2:
				if strings.Contains(chainId, constant.AllChain) {
					if chains[0] == chains[1] && chains[0] == constant.AllChain {
						//nothing to do
					} else {
						index := strings.Index(chainId, constant.AllChain)
						if index > 0 {
							filter["$or"] = []bson.M{
								{RelayerFieldChainA: chains[0]},
								{RelayerFieldChainB: chains[0]},
							}
						} else {
							filter["$or"] = []bson.M{
								{RelayerFieldChainA: chains[1]},
								{RelayerFieldChainB: chains[1]},
							}
						}
					}
				} else {
					filter["$or"] = []bson.M{
						{RelayerFieldChainA: chains[0], RelayerFieldChainB: chains[1]},
						{RelayerFieldChainA: chains[1], RelayerFieldChainB: chains[0]},
					}
				}
				break
			}
		}
	}
	if status > 0 {
		filter[RelayerFieldStatus] = status
	}
	return filter
}

func (repo *IbcRelayerRepo) FindAllBycond(chainId string, status int, skip, limit int64, useCount bool) ([]*entity.IBCRelayer, int64, error) {
	var (
		res   []*entity.IBCRelayer
		total int64
	)
	filter := repo.analyzeCond(chainId, status)
	err := repo.coll().Find(context.Background(), filter).Skip(skip).Limit(limit).Sort("-" + RelayerFieldTransferTotalTxs).All(&res)
	if useCount {
		total, err = repo.coll().Find(context.Background(), filter).Count()
	}
	return res, total, err
}

func (repo *IbcRelayerRepo) CountBycond(chainId string, status int) (int64, error) {
	filter := repo.analyzeCond(chainId, status)
	return repo.coll().Find(context.Background(), filter).Count()
}

func (repo *IbcRelayerRepo) Insert(relayer []entity.IBCRelayer) error {
	if _, err := repo.coll().InsertMany(context.Background(), relayer); err != nil {
		return err
	}
	return nil
}

func (repo *IbcRelayerRepo) UpdateTxsInfo(relayerId string, txs, txsSuccess int64, totalValue string) error {
	updateData := bson.M{
		RelayerFieldUpdateAt: time.Now().Unix(),
	}
	if totalValue != "" {
		updateData[RelayerFieldTransferTotalTxsValue] = totalValue
	}
	if txs > 0 {
		updateData[RelayerFieldTransferTotalTxs] = txs
	}

	if txsSuccess > 0 {
		updateData[RelayerFieldTransferSuccessTxs] = txsSuccess
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": updateData})
}
func (repo *IbcRelayerRepo) UpdateStatusAndTime(relayerId string, status int, updateTime, timePeriod int64) error {
	update := bson.M{
		RelayerFieldUpdateAt: time.Now().Unix(),
	}
	if status > 0 {
		update[RelayerFieldStatus] = status
	}
	if updateTime > 0 {
		update[RelayerFieldUpdateTime] = updateTime
	}
	if timePeriod > 0 {
		update[RelayerFieldTimePeriod] = timePeriod
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": update})
}

func (repo *IbcRelayerRepo) UpdateSrcAddress(relayerId string, addrs []string) error {

	if len(addrs) == 0 {
		return nil
	}
	update := bson.M{
		RelayerFieldUpdateAt:         time.Now().Unix(),
		RelayerFieldChainAAddress:    addrs[0],
		RelayerFieldChainAALLAddress: addrs,
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{
		RelayerFieldelayerId:      relayerId,
		RelayerFieldChainAAddress: "",
	}, bson.M{
		"$set": update})
}

func (repo *IbcRelayerRepo) FindLatestOne() (*entity.IBCRelayer, error) {
	var res *entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{}).Sort("-" + RelayerFieldLatestTxTime).One(&res)
	return res, err
}

func (repo *IbcRelayerRepo) CountChainRelayers(chainId string) (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{
		RelayerFieldStatus: entity.RelayerRunning,
		"$or": []bson.M{
			{RelayerFieldChainA: chainId},
			{RelayerFieldChainB: chainId},
		},
	}).Count()
}

func (repo *IbcRelayerRepo) FindRelayer(chainId, relayerAddr, channel string) (*entity.IBCRelayer, error) {
	var res *entity.IBCRelayer
	err := repo.coll().Find(context.Background(), bson.M{
		"$or": []bson.M{
			{RelayerFieldChainA: chainId, RelayerFieldChannelA: channel, RelayerFieldChainAAddress: relayerAddr},
			{RelayerFieldChainB: chainId, RelayerFieldChannelB: channel, RelayerFieldChainBAddress: relayerAddr},
		},
	}).One(&res)
	return res, err
}

func (repo *IbcRelayerRepo) CountChannelRelayers() ([]*dto.CountChannelRelayersDTO, error) {
	match := bson.M{
		"$match": bson.M{
			RelayerFieldStatus: entity.RelayerRunning,
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"chain_a":   "$chain_a",
				"chain_b":   "$chain_b",
				"channel_a": "$channel_a",
				"channel_b": "$channel_b",
			},
			"count": bson.M{
				"$sum": 1,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":       0,
			"chain_a":   "$_id.chain_a",
			"chain_b":   "$_id.chain_b",
			"channel_a": "$_id.channel_a",
			"channel_b": "$_id.channel_b",
			"count":     "$count",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	var res []*dto.CountChannelRelayersDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
