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
	GetFirstTx(chainId string) (*entity.Tx, error)
	GetRelayerScChainAddr(packetId, chainId string) (string, error)
	GetUpdateTimeByUpdateClient(chainId, address, clientId string, startTime int64) (int64, error)
	GetLatestRecvPacketTime(chainId, address, channelId string, startTime int64) (int64, error)
	GetChannelOpenConfirmTime(chainId, channelId string) (int64, error)
	GetTransferTx(chainId string, height, limit int64) ([]*entity.Tx, error)
	FindByTypeAndHeight(chainId, txType string, height int64) ([]*entity.Tx, error)
	GetTxByHash(chainId string, hash string) (entity.Tx, error)
	GetTxByHashes(chainId string, hashs []string) ([]*entity.Tx, error)
	GetAcknowledgeTxs(chainId, packetId string) ([]*entity.Tx, error)
	GetRecvPacketTxs(chainId, packetId string) ([]*entity.Tx, error)
	FindByPacketIds(chainId, txType string, packetIds []string, status *entity.TxStatus) ([]*entity.Tx, error)
	FindAllAckTxs(chainId string, height int64) ([]*entity.Tx, error)
	FindHeight(chainId string, min bool) (entity.Tx, error)
	UpdateAckPacketId(chainId string, height int64, txHash string, msgs []interface{}) error
	RelayerDenomStatistics(chainId string, startTime, endTime int64) ([]*dto.RelayerDenomStatisticsDTO, error)
	RelayerFeeStatistics(chainId string, startTime, endTime int64) ([]*dto.RelayerFeeStatisticsDTO, error)
}

var _ ITxRepo = new(TxRepo)

type TxRepo struct {
}

func (repo *TxRepo) coll(chainId string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.Tx{}.CollectionName(chainId))
}

func (repo *TxRepo) GetFirstTx(chainId string) (*entity.Tx, error) {
	var res entity.Tx
	err := repo.coll(chainId).Find(context.Background(), bson.M{}).Sort("time").One(&res)
	return &res, err
}

func (repo *TxRepo) GetRelayerScChainAddr(packetId, chainId string) (string, error) {
	var res entity.Tx
	//get relayer address by packet_id and acknowledge_packet or timeout_packet
	err := repo.coll(chainId).Find(context.Background(), bson.M{
		"msgs.msg.packet_id": packetId,
		"msgs.type": bson.M{ //filter ibc transfer
			"$in": []string{constant.MsgTypeAcknowledgement, constant.MsgTypeTimeoutPacket},
		},
	}).Sort("-height").Limit(1).One(&res)
	if len(res.DocTxMsgs) > 0 {
		for _, msg := range res.DocTxMsgs {
			cmsg := msg.CommonMsg()
			if cmsg.PacketId == packetId {
				return cmsg.Signer, nil
			}
		}
	}
	return "", err
}

// return value description
//1: latest update_client tx_time
//2: error
func (repo *TxRepo) GetUpdateTimeByUpdateClient(chainId, address, clientId string, startTime int64) (int64, error) {
	var res *entity.Tx
	query := bson.M{
		"msgs.type":          constant.MsgTypeUpdateClient,
		"msgs.msg.signer":    address,
		"msgs.msg.client_id": clientId,
		"time": bson.M{
			"$gte": startTime,
		},
	}
	err := repo.coll(chainId).Find(context.Background(), query).
		Select(bson.M{"time": 1, "msgs.type": 1}).Sort("-time").Hint("msgs.msg.signer_1_msgs.type_1_time_1").Limit(1).One(&res)
	if err != nil {
		return 0, err
	}

	if res != nil {
		return res.Time, nil
	}
	return 0, nil
}

func (repo *TxRepo) GetLatestRecvPacketTime(chainId, address, channelId string, startTime int64) (int64, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.type":                      constant.MsgTypeRecvPacket,
		"msgs.msg.signer":                address,
		"msgs.msg.packet.source_channel": channelId,
		"time": bson.M{
			"$gte": startTime,
		},
	}
	err := repo.coll(chainId).Find(context.Background(), query).
		Select(bson.M{"time": 1}).Sort("-time").Hint("msgs.msg.signer_1_msgs.type_1_time_1").Limit(1).All(&res)
	if err != nil {
		return 0, err
	}

	if len(res) == 1 {
		return res[0].Time, nil
	}
	return 0, nil
}

func (repo *TxRepo) GetChannelOpenConfirmTime(chainId, channelId string) (int64, error) {
	var res entity.Tx
	query := bson.M{
		"msgs.type":           constant.MsgTypeChannelOpenConfirm,
		"msgs.msg.channel_id": channelId,
	}
	err := repo.coll(chainId).Find(context.Background(), query).
		Select(bson.M{"time": 1}).Sort("-time").Limit(1).One(&res)

	if err != nil {
		return 0, err
	}
	return res.Time, nil
}

func (repo *TxRepo) GetTransferTx(chainId string, height, limit int64) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"types": constant.MsgTypeTransfer,
		"height": bson.M{
			"$gt": height,
		},
	}

	err := repo.coll(chainId).Find(context.Background(), query).Sort("height").Limit(limit).All(&res)
	return res, err
}

func (repo *TxRepo) FindByTypeAndHeight(chainId, txType string, height int64) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"types":  txType,
		"height": height,
	}

	err := repo.coll(chainId).Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *TxRepo) FindByPacketIds(chainId, txType string, packetIds []string, status *entity.TxStatus) ([]*entity.Tx, error) {
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

	err := repo.coll(chainId).Find(context.Background(), query).All(&res)
	return res, err
}

func (repo *TxRepo) GetTxByHash(chainId string, hash string) (entity.Tx, error) {
	var res entity.Tx
	err := repo.coll(chainId).Find(context.Background(), bson.M{"tx_hash": hash}).Sort("-height").One(&res)
	return res, err
}
func (repo *TxRepo) GetTxByHashes(chainId string, hashs []string) ([]*entity.Tx, error) {
	var res []*entity.Tx
	err := repo.coll(chainId).Find(context.Background(), bson.M{"tx_hash": bson.M{
		"$in": hashs,
	}}).Sort("height").All(&res)
	return res, err
}

func (repo *TxRepo) GetAcknowledgeTxs(chainId, packetId string) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.msg.packet_id": packetId,
		"msgs.type":          constant.MsgTypeAcknowledgement,
		"status":             entity.TxStatusSuccess,
	}
	//取"成功"状态最新的acknowledge_tx交易
	err := repo.coll(chainId).Find(context.Background(), query).Sort("-height").All(&res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (repo *TxRepo) GetRecvPacketTxs(chainId, packetId string) ([]*entity.Tx, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.type":          constant.MsgTypeRecvPacket,
		"msgs.msg.packet_id": packetId,
	}
	err := repo.coll(chainId).Find(context.Background(), query).Sort("-height").All(&res)

	if err != nil {
		return res, err
	}
	return res, nil
}

func (repo *TxRepo) FindAllAckTxs(chainId string, height int64) ([]*entity.Tx, error) {
	var txs []*entity.Tx
	err := repo.coll(chainId).Find(context.Background(), bson.M{
		"types":  constant.MsgTypeAcknowledgement,
		"height": bson.M{"$gt": height, "$lte": height + constant.IncreHeight},
	}).Sort("height").All(&txs)
	return txs, err
}

func (repo *TxRepo) FindHeight(chainId string, min bool) (entity.Tx, error) {
	var tx entity.Tx
	sorts := "-height"
	if min {
		sorts = "+height"
	}
	err := repo.coll(chainId).Find(context.Background(), bson.M{
		"types": constant.MsgTypeAcknowledgement,
	}).Sort(sorts).Limit(1).One(&tx)
	return tx, err
}

func (repo *TxRepo) UpdateAckPacketId(chainId string, height int64, txHash string, msgs []interface{}) error {
	filter, update := bson.M{
		"height": height, "tx_hash": txHash,
	}, bson.M{
		"$set": bson.M{
			"msgs": msgs,
		},
	}
	err := repo.coll(chainId).UpdateOne(context.Background(), filter, update)
	return err
}

func (repo *TxRepo) RelayerDenomStatistics(chainId string, startTime, endTime int64) ([]*dto.RelayerDenomStatisticsDTO, error) {
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
	err := repo.coll(chainId).Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *TxRepo) RelayerFeeStatistics(chainId string, startTime, endTime int64) ([]*dto.RelayerFeeStatisticsDTO, error) {
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
	err := repo.coll(chainId).Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
