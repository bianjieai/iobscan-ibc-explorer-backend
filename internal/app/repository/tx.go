package repository

import (
	"context"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITxRepo interface {
	GetFirstTx(chain string) (*entity.Tx, error)
	GetUpdateTimeByUpdateClient(chain, address, clientId string, startTime int64) (int64, error)
	GetLatestRecvPacketTime(chain, address, channelId string, startTime int64) (int64, error)
	GetChannelOpenConfirmTime(chain, channelId string) (int64, error)
	GetTransferTx(chain string, height, limit int64) ([]*entity.Tx, error)
	FindByTypeAndHeight(chain, txType string, height int64) ([]*entity.Tx, error)
	GetTxByHash(chain string, hash string) (entity.Tx, error)
	GetTxByHashes(chain string, hashs []string) ([]*entity.Tx, error)
	GetAcknowledgeTxs(chain, packetId string) ([]*entity.Tx, error)
	GetRecvPacketTxs(chain, packetId string) ([]*entity.Tx, error)
	FindByPacketIds(chain, txType string, packetIds []string, status *entity.TxStatus) ([]*entity.Tx, error)
	RelayerDenomStatistics(chain string, startTime, endTime int64) ([]*dto.RelayerDenomStatisticsDTO, error)
	RelayerFeeStatistics(chain string, startTime, endTime int64) ([]*dto.RelayerFeeStatisticsDTO, error)
	GetRelayerTxs(chain string, relayerAddrs []string, txTypes []string,
		txTimeStart, txTimeEnd, skip, limit int64) ([]*entity.Tx, error)
	CountRelayerTxs(chain string, relayerAddrs []string, txTypes []string,
		txTimeStart, txTimeEnd int64) (int64, error)
	GetAddressTxs(chain, address string, skip, limit int64) ([]*entity.Tx, error)
	CountAddressTxs(chain, address string) (int64, error)
	GetAddressLatestTx(chain, address string) (*entity.Tx, error)
}

var _ ITxRepo = new(TxRepo)

type TxRepo struct {
}

func (repo *TxRepo) coll(chain string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.Tx{}.CollectionName(chain))
}

func (repo *TxRepo) GetFirstTx(chain string) (*entity.Tx, error) {
	var res entity.Tx
	err := repo.coll(chain).Find(context.Background(), bson.M{}).Sort("time").One(&res)
	return &res, err
}

// return value description
//1: latest update_client tx_time
//2: error
func (repo *TxRepo) GetUpdateTimeByUpdateClient(chain, address, clientId string, startTime int64) (int64, error) {
	var res *entity.Tx
	query := bson.M{
		"msgs.type":          constant.MsgTypeUpdateClient,
		"msgs.msg.signer":    address,
		"msgs.msg.client_id": clientId,
		"time": bson.M{
			"$gte": startTime,
		},
	}
	err := repo.coll(chain).Find(context.Background(), query).
		Select(bson.M{"time": 1, "msgs.type": 1}).Sort("-time").Hint(GetRelayerUpdateTimeHintIndexName()).Limit(1).One(&res)
	if err != nil {
		return 0, err
	}

	if res != nil {
		return res.Time, nil
	}
	return 0, nil
}

func (repo *TxRepo) GetLatestRecvPacketTime(chain, address, channelId string, startTime int64) (int64, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.type":                      constant.MsgTypeRecvPacket,
		"msgs.msg.signer":                address,
		"msgs.msg.packet.source_channel": channelId,
		"time": bson.M{
			"$gte": startTime,
		},
	}
	err := repo.coll(chain).Find(context.Background(), query).
		Select(bson.M{"time": 1}).Sort("-time").Hint(GetLatestRecvPacketTimeHintIndexName()).Limit(1).All(&res)
	if err != nil {
		return 0, err
	}

	if len(res) == 1 {
		return res[0].Time, nil
	}
	return 0, nil
}

func (repo *TxRepo) GetChannelOpenConfirmTime(chain, channelId string) (int64, error) {
	var res entity.Tx
	query := bson.M{
		"types":               constant.MsgTypeChannelOpenConfirm,
		"msgs.msg.channel_id": channelId,
	}
	err := repo.coll(chain).Find(context.Background(), query).
		Select(bson.M{"time": 1}).Sort("-time").Limit(1).One(&res)

	if err != nil {
		return 0, err
	}
	return res.Time, nil
}

func (repo *TxRepo) GetTransferTx(chain string, height, limit int64) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"types": constant.MsgTypeTransfer,
		"height": bson.M{
			"$gt": height,
		},
	}

	err := repo.coll(chain).Find(context.Background(), query).Sort("height").Limit(limit).All(&res)
	return res, err
}

func (repo *TxRepo) FindByTypeAndHeight(chain, txType string, height int64) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"types":  txType,
		"height": height,
	}

	err := repo.coll(chain).Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *TxRepo) FindByPacketIds(chain, txType string, packetIds []string, status *entity.TxStatus) ([]*entity.Tx, error) {
	if len(packetIds) == 0 {
		return nil, nil
	}

	var res []*entity.Tx
	query := bson.M{
		"msgs.type": txType,
		"msgs.msg.packet_id": bson.M{
			"$in": packetIds,
		},
	}
	if status != nil {
		query["status"] = status
	}

	err := repo.coll(chain).Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *TxRepo) GetTxByHash(chain string, hash string) (entity.Tx, error) {
	var res entity.Tx
	err := repo.coll(chain).Find(context.Background(), bson.M{"tx_hash": hash}).Sort("-height").One(&res)
	return res, err
}
func (repo *TxRepo) GetTxByHashes(chain string, hashs []string) ([]*entity.Tx, error) {
	var res []*entity.Tx
	err := repo.coll(chain).Find(context.Background(), bson.M{"tx_hash": bson.M{
		"$in": hashs,
	}}).Sort("height").All(&res)
	return res, err
}

func (repo *TxRepo) GetAcknowledgeTxs(chain, packetId string) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.msg.packet_id": packetId,
		"msgs.type":          constant.MsgTypeAcknowledgement,
		"status":             entity.TxStatusSuccess,
	}
	//取"成功"状态最新的acknowledge_tx交易
	err := repo.coll(chain).Find(context.Background(), query).Sort("-height").All(&res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (repo *TxRepo) GetRecvPacketTxs(chain, packetId string) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.type":          constant.MsgTypeRecvPacket,
		"msgs.msg.packet_id": packetId,
	}
	err := repo.coll(chain).Find(context.Background(), query).Sort("-height").All(&res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (repo *TxRepo) RelayerDenomStatistics(chain string, startTime, endTime int64) ([]*dto.RelayerDenomStatisticsDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"time": bson.M{
				"$lte": endTime,
				"$gte": startTime,
			},
			"msgs.type": bson.M{
				"$in": []entity.TxType{entity.TxTypeRecvPacket, entity.TxTypeAckPacket, entity.TxTypeTimeoutPacket},
			},
		},
	}

	unwind := bson.M{
		"$unwind": "$msgs",
	}

	match2 := bson.M{
		"$match": bson.M{
			"msgs.type": bson.M{
				"$in": []entity.TxType{entity.TxTypeRecvPacket, entity.TxTypeAckPacket, entity.TxTypeTimeoutPacket},
			},
		},
	}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"signer":     "$msgs.msg.signer",
				"status":     "$status",
				"tx_type":    "$msgs.type",
				"denom":      "$msgs.msg.packet.data.denom",
				"sc_channel": "$msgs.msg.packet.source_channel",
				"dc_channel": "$msgs.msg.packet.destination_channel",
			},
			"denom_amount": bson.M{
				"$sum": bson.M{
					"$toDouble": "$msgs.msg.packet.data.amount",
				},
			},
			"txs_count": bson.M{
				"$sum": 1,
			},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"_id":          0,
			"signer":       "$_id.signer",
			"status":       "$_id.status",
			"tx_type":      "$_id.tx_type",
			"denom":        "$_id.denom",
			"sc_channel":   "$_id.sc_channel",
			"dc_channel":   "$_id.dc_channel",
			"denom_amount": "$denom_amount",
			"txs_count":    "$txs_count",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, unwind, match2, group, project)
	var res []*dto.RelayerDenomStatisticsDTO
	err := repo.coll(chain).Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *TxRepo) RelayerFeeStatistics(chain string, startTime, endTime int64) ([]*dto.RelayerFeeStatisticsDTO, error) {
	match := bson.M{
		"$match": bson.M{
			"time": bson.M{
				"$lte": endTime,
				"$gte": startTime,
			},
			"msgs.type": bson.M{
				"$in": []entity.TxType{entity.TxTypeRecvPacket, entity.TxTypeAckPacket, entity.TxTypeTimeoutPacket},
			},
		},
	}

	unwind := bson.M{
		"$unwind": "$msgs",
	}

	match2 := bson.M{
		"$match": bson.M{
			"msgs.type": bson.M{
				"$in": []entity.TxType{entity.TxTypeRecvPacket, entity.TxTypeAckPacket, entity.TxTypeTimeoutPacket},
			},
		},
	}

	unwind2 := bson.M{
		"$unwind": "$fee.amount",
	}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"signer":  "$msgs.msg.signer",
				"status":  "$status",
				"tx_type": "$msgs.type",
				"denom":   "$fee.amount.denom",
			},
			"denom_amount": bson.M{
				"$sum": bson.M{
					"$toDouble": "$fee.amount.amount",
				},
			},
			"txs_count": bson.M{
				"$sum": 1,
			},
		},
	}

	project := bson.M{
		"$project": bson.M{
			"_id":          0,
			"signer":       "$_id.signer",
			"status":       "$_id.status",
			"tx_type":      "$_id.tx_type",
			"denom":        "$_id.denom",
			"denom_amount": "$denom_amount",
			"txs_count":    "$txs_count",
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, unwind, match2, unwind2, group, project)
	var res []*dto.RelayerFeeStatisticsDTO
	err := repo.coll(chain).Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
func createQueryRelayerTxs(relayerAddrs []string, txTypes []string, txTimeStart, txTimeEnd int64) bson.M {
	query := bson.M{}
	if len(relayerAddrs) > 0 {
		query["msgs.msg.signer"] = bson.M{
			"$in": relayerAddrs,
		}
	}
	if len(txTypes) > 0 {
		query["msgs.type"] = bson.M{
			"$in": txTypes,
		}
	}

	//time
	if txTimeStart > 0 && txTimeEnd > 0 {
		query["time"] = bson.M{
			"$gte": txTimeStart,
			"$lte": txTimeEnd,
		}
	} else if txTimeStart > 0 {
		query["time"] = bson.M{
			"$gte": txTimeStart,
		}
	} else if txTimeEnd > 0 {
		query["time"] = bson.M{
			"$lte": txTimeEnd,
		}
	}
	return query
}
func (repo *TxRepo) GetRelayerTxs(chain string, relayerAddrs []string, txTypes []string,
	txTimeStart, txTimeEnd, skip, limit int64) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := createQueryRelayerTxs(relayerAddrs, txTypes, txTimeStart, txTimeEnd)
	err := repo.coll(chain).Find(context.Background(), query).Sort("-time").Hint(GetRelayerTxsHintIndexName()).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *TxRepo) CountRelayerTxs(chain string, relayerAddrs []string, txTypes []string, txTimeStart, txTimeEnd int64) (int64, error) {
	query := createQueryRelayerTxs(relayerAddrs, txTypes, txTimeStart, txTimeEnd)
	return repo.coll(chain).Find(context.Background(), query).Hint(CountRelayerTxsHintIndexName()).Count()
}

func (repo *TxRepo) GetAddressTxs(chain, address string, skip, limit int64) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"addrs": address,
		"msgs.type": bson.M{
			"$in": entity.ICS20TransferTxTypes,
		},
	}
	err := repo.coll(chain).Find(context.Background(), query).Sort("-time").Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *TxRepo) CountAddressTxs(chain, address string) (int64, error) {
	query := bson.M{
		"addrs": address,
		"msgs.type": bson.M{
			"$in": entity.ICS20TransferTxTypes,
		},
	}
	return repo.coll(chain).Find(context.Background(), query).Count()
}

func (repo *TxRepo) GetAddressLatestTx(chain, address string) (*entity.Tx, error) {
	var res entity.Tx
	query := bson.M{
		"addrs": address,
		"msgs.type": bson.M{
			"$in": entity.ICS20AllTxTypes,
		},
	}
	err := repo.coll(chain).Find(context.Background(), query).Sort("-time").One(&res)
	return &res, err
}
