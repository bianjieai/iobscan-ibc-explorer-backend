package repository

import (
	"context"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	RelayerFieldelayerId      = "relayer_id"
	RelayerFieldTotalTxs      = "relayed_total_txs"
	RelayerFieldSuccessTxs    = "relayed_success_txs"
	RelayerFieldTotalTxsValue = "relayed_total_txs_value"
	RelayerFieldTotalFeeValue = "total_fee_value"
	RelayerFieldServedChains  = "served_chains"
	RelayerFieldeRelayerName  = "relayer_name"
	RelayerFieldUpdateTime    = "update_time"
	RelayerFieldChainA        = "channel_pair_info.chain_a"
	RelayerFieldChainB        = "channel_pair_info.chain_b"
	RelayerFieldChannelA      = "channel_pair_info.channel_a"
	RelayerFieldChannelB      = "channel_pair_info.channel_b"
	RelayerFieldChainAAddress = "channel_pair_info.chain_a_address"
	RelayerFieldChainBAddress = "channel_pair_info.chain_b_address"
	RelayerFieldUpdateAt      = "update_at"

	RelayerAllType      = 0
	RelayerRegisterType = 1
	RelayerUnknowType   = 2
)

type IRelayerRepo interface {
	InsertOne(relayer *entity.IBCRelayerNew) error
	InsertBatch(relayer []entity.IBCRelayerNew) error
	UpdateRelayerTime(relayerId string, updateTime int64) error
	UpdateTxsInfo(relayerId string, txs, txsSuccess int64, totalValue, totalFeeValue string) error
	FindAll(skip, limit int64, relayType int) ([]*entity.IBCRelayerNew, error)
	FindAllBycond(relayerName, relayerAddr string, skip, limit int64, useCount bool) ([]*entity.IBCRelayerNew, int64, error)
	CountBycond(relayerName, relayerAddr string) (int64, error)
	CountChainRelayers(chainId string) (int64, error)
	CountChannelRelayers(chainA, channelA, chainB, channelB string) (int64, error)
	//FindRelayer(chainId, relayerAddr, channel string) ([]*entity.IBCRelayerNew, error)
	FindOneByRelayerId(relayerId string) (*entity.IBCRelayerNew, error)
	FindOneByRelayerName(name string) (*entity.IBCRelayerNew, error)
	FindEmptyAddrAll(skip, limit int64) ([]*entity.IBCRelayerNew, error)
	//UpdateSrcAddress(relayerId string, addrs []string) error
	UpdateChannelPairInfo(relayerId string, infos []entity.ChannelPairInfo) error
	Update(relayer *entity.IBCRelayerNew) error
	RemoveDumpData(ids []string) error
	FindUnknownByAddrPair(addrA, addrB string) ([]*entity.IBCRelayerNew, error)
}

var _ IRelayerRepo = new(IbcRelayerRepo)

type IbcRelayerRepo struct {
}

func (repo *IbcRelayerRepo) Update(relayer *entity.IBCRelayerNew) error {
	updateData := bson.M{
		RelayerFieldUpdateAt: time.Now().Unix(),
	}
	if len(relayer.ChannelPairInfo) > 0 {
		updateData["channel_pair_info"] = relayer.ChannelPairInfo
	}
	if relayer.RelayedTotalTxsValue != "" {
		updateData[RelayerFieldTotalTxsValue] = relayer.RelayedTotalTxsValue
	}
	if relayer.TotalFeeValue != "" {
		updateData[RelayerFieldTotalFeeValue] = relayer.TotalFeeValue
	}
	if relayer.RelayedTotalTxs > 0 {
		updateData[RelayerFieldTotalTxs] = relayer.RelayedTotalTxs
	}

	if relayer.RelayedSuccessTxs > 0 {
		updateData[RelayerFieldSuccessTxs] = relayer.RelayedSuccessTxs
	}

	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayer.RelayerId}, bson.M{
		"$set": updateData})
}

func (repo *IbcRelayerRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.IBCRelayerNew{}.CollectionName())
}

func (repo *IbcRelayerRepo) FindAll(skip, limit int64, relayType int) ([]*entity.IBCRelayerNew, error) {
	var res []*entity.IBCRelayerNew
	cond := bson.M{}
	switch relayType {
	case RelayerRegisterType:
		cond = bson.M{RelayerFieldeRelayerName: bson.M{"$ne": ""}}
	case RelayerUnknowType:
		cond = bson.M{RelayerFieldeRelayerName: ""}
	}
	err := repo.coll().Find(context.Background(), cond).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *IbcRelayerRepo) RemoveDumpData(ids []string) error {
	_, err := repo.coll().RemoveAll(context.Background(), bson.M{RelayerFieldelayerId: bson.M{
		"$in": ids,
	}})
	return err
}

func (repo *IbcRelayerRepo) FindChannelPairInfos() ([]*entity.IBCRelayerNew, error) {
	var res []*entity.IBCRelayerNew
	err := repo.coll().Find(context.Background(), bson.M{}).Select(bson.M{RelayerFieldelayerId: 1, "channel_pair_info": 1}).All(&res)
	return res, err
}

func (repo *IbcRelayerRepo) FindEmptyAddrAll(skip, limit int64) ([]*entity.IBCRelayerNew, error) {
	var res []*entity.IBCRelayerNew
	err := repo.coll().Find(context.Background(), bson.M{RelayerFieldChainAAddress: ""}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *IbcRelayerRepo) analyzeCond(relayerName, relayerAddr string) bson.M {
	filter := bson.M{}
	if relayerName != "" {
		filter[RelayerFieldeRelayerName] = bson.M{
			"$regex": relayerName + "*",
		}
	}
	if relayerAddr != "" {
		filter["$or"] = []bson.M{
			{RelayerFieldChainAAddress: relayerAddr},
			{RelayerFieldChainBAddress: relayerAddr},
		}
	}
	return filter
}

func (repo *IbcRelayerRepo) FindAllBycond(relayerName, relayerAddr string, skip, limit int64, useCount bool) ([]*entity.IBCRelayerNew, int64, error) {
	var (
		res   []*entity.IBCRelayerNew
		total int64
	)
	filter := repo.analyzeCond(relayerName, relayerAddr)
	err := repo.coll().Find(context.Background(), filter).Skip(skip).Limit(limit).Sort("-"+RelayerFieldeRelayerName, "-"+RelayerFieldTotalTxs).All(&res)
	if useCount {
		total, err = repo.coll().Find(context.Background(), filter).Count()
	}
	return res, total, err
}

func (repo *IbcRelayerRepo) CountBycond(relayerName, relayerAddr string) (int64, error) {
	filter := repo.analyzeCond(relayerName, relayerAddr)
	return repo.coll().Find(context.Background(), filter).Count()
}

func (repo *IbcRelayerRepo) InsertOne(relayer *entity.IBCRelayerNew) error {
	if _, err := repo.coll().InsertOne(context.Background(), relayer); err != nil {
		return err
	}
	return nil
}

func (repo *IbcRelayerRepo) InsertBatch(relayer []entity.IBCRelayerNew) error {
	if _, err := repo.coll().InsertMany(context.Background(), relayer, insertIgnoreErrOpt); err != nil && !qmgo.IsDup(err) {
		return err
	}
	return nil
}
func (repo *IbcRelayerRepo) UpdateChannelPairInfo(relayerId string, infos []entity.ChannelPairInfo) error {
	updateData := bson.M{
		RelayerFieldUpdateAt: time.Now().Unix(),
	}
	if len(infos) > 0 {
		updateData["channel_pair_info"] = infos
	}

	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": updateData})
}

func (repo *IbcRelayerRepo) UpdateTxsInfo(relayerId string, txs, txsSuccess int64, totalValue, totalFeeValue string) error {
	updateData := bson.M{
		RelayerFieldUpdateAt: time.Now().Unix(),
	}
	if totalValue != "" {
		updateData[RelayerFieldTotalTxsValue] = totalValue
	}
	if totalFeeValue != "" {
		updateData[RelayerFieldTotalFeeValue] = totalFeeValue
	}
	if txs > 0 {
		updateData[RelayerFieldTotalTxs] = txs
	}

	if txsSuccess > 0 {
		updateData[RelayerFieldSuccessTxs] = txsSuccess
	}
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": updateData})
}
func (repo *IbcRelayerRepo) UpdateRelayerTime(relayerId string, updateTime int64) error {
	return repo.coll().UpdateOne(context.Background(), bson.M{RelayerFieldelayerId: relayerId}, bson.M{
		"$set": bson.M{
			RelayerFieldUpdateAt:   time.Now().Unix(),
			RelayerFieldUpdateTime: updateTime,
		}})
}

//func (repo *IbcRelayerRepo) UpdateSrcAddress(relayerId string, addrs []string) error {
//
//	if len(addrs) == 0 {
//		return nil
//	}
//	update := bson.M{
//		RelayerFieldUpdateAt:      time.Now().Unix(),
//		RelayerFieldChainAAddress: addrs[0],
//	}
//	return repo.coll().UpdateOne(context.Background(), bson.M{
//		RelayerFieldelayerId:      relayerId,
//		RelayerFieldChainAAddress: "",
//	}, bson.M{
//		"$set": update})
//}

func (repo *IbcRelayerRepo) CountChainRelayers(chainId string) (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{
		"$or": []bson.M{
			{RelayerFieldChainA: chainId},
			{RelayerFieldChainB: chainId},
		},
	}).Count()
}

//func (repo *IbcRelayerRepo) FindRelayer(chainId, relayerAddr, channel string) ([]*entity.IBCRelayerNew, error) {
//	var res []*entity.IBCRelayerNew
//	err := repo.coll().Find(context.Background(), bson.M{
//		"$or": []bson.M{
//			{RelayerFieldChainA: chainId, RelayerFieldChannelA: channel, RelayerFieldChainAAddress: relayerAddr},
//			{RelayerFieldChainB: chainId, RelayerFieldChannelB: channel, RelayerFieldChainBAddress: relayerAddr},
//		},
//	}).All(&res)
//	return res, err
//}

func (repo *IbcRelayerRepo) CountChannelRelayers(chainA, channelA, chainB, channelB string) (int64, error) {
	return repo.coll().Find(context.Background(), bson.M{
		RelayerFieldChainA: chainA, RelayerFieldChannelA: channelA,
		RelayerFieldChainB: chainB, RelayerFieldChannelB: channelB,
	}).Count()
}

func (repo *IbcRelayerRepo) FindOneByRelayerId(relayerId string) (*entity.IBCRelayerNew, error) {
	var res *entity.IBCRelayerNew
	err := repo.coll().Find(context.Background(), bson.M{RelayerFieldelayerId: relayerId}).One(&res)
	return res, err
}

func (repo *IbcRelayerRepo) FindOneByRelayerName(name string) (*entity.IBCRelayerNew, error) {
	var res *entity.IBCRelayerNew
	err := repo.coll().Find(context.Background(), bson.M{RelayerFieldeRelayerName: name}).One(&res)
	return res, err
}

func (repo *IbcRelayerRepo) FindUnknownByAddrPair(addrA, addrB string) ([]*entity.IBCRelayerNew, error) {
	var res []*entity.IBCRelayerNew
	err := repo.coll().Find(context.Background(), bson.M{RelayerFieldChainAAddress: addrA, RelayerFieldChainBAddress: addrB, RelayerFieldeRelayerName: ""}).All(&res)
	return res, err
}
