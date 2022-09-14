package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ITxRepo interface {
	GetRelayerScChainAddr(packetId, chainId string) (string, error)
	GetTimePeriodByUpdateClient(chainId, address, clientId string, startTime int64) (int64, int64, error)
	GetLatestRecvPacketTime(chainId, address, channelId string, startTime int64) (int64, error)
	GetChannelOpenConfirmTime(chainId, channelId string) (int64, error)
	GetTransferTx(chainId string, height, limit int64) ([]*entity.Tx, error)
	FindByTypeAndHeight(chainId, txType string, height int64) ([]*entity.Tx, error)
	GetTxByHash(chainId string, hash string) (entity.Tx, error)
	GetTxByHashes(chainId string, hashs []string) ([]*entity.Tx, error)
	GetAcknowledgeTxs(chainId, packetId string) (entity.Tx, error)
	GetRecvPacketTxs(chainId, packetId string) ([]*entity.Tx, error)
	FindByPacketIds(chainId, txType string, packetIds []string, status *entity.TxStatus) ([]*entity.Tx, error)
	FindAllAckTxs(chainId string, height int64) ([]*entity.Tx, error)
	FindHeight(chainId string, min bool) (entity.Tx, error)
	UpdateAckPacketId(chainId string, height int64, txHash, packetId string) error
}

var _ ITxRepo = new(TxRepo)

type TxRepo struct {
}

func (repo *TxRepo) coll(chainId string) *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.Tx{}.CollectionName(chainId))
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
//2: time_period
//3: error
func (repo *TxRepo) GetTimePeriodByUpdateClient(chainId, address, clientId string, startTime int64) (int64, int64, error) {
	var res []*entity.Tx
	query := bson.M{
		"msgs.type":          constant.MsgTypeUpdateClient,
		"msgs.msg.signer":    address,
		"msgs.msg.client_id": clientId,
		"time": bson.M{
			"$gte": startTime,
		},
	}
	err := repo.coll(chainId).Find(context.Background(), query).
		Select(bson.M{"time": 1, "msgs.type": 1}).Sort("-time").Hint("msgs.msg.signer_1_msgs.type_1_time_1").Limit(2).All(&res)
	if err != nil {
		return 0, 0, err
	}

	if len(res) == 2 {
		return res[0].Time, res[0].Time - res[1].Time, nil
	}
	if len(res) == 1 {
		return res[0].Time, -1, nil
	}
	return 0, -1, nil
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

func (repo *TxRepo) GetAcknowledgeTxs(chainId, packetId string) (entity.Tx, error) {
	var res entity.Tx
	query := bson.M{
		"msgs.msg.packet_id": packetId,
		"msgs.type":          constant.MsgTypeAcknowledgement,
		"status":             entity.TxStatusSuccess,
	}
	err := repo.coll(chainId).Find(context.Background(), query).One(&res)

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
		"types":              constant.MsgTypeAcknowledgement,
		"height":             bson.M{"$gt": height, "$lte": height + constant.IncreHeight},
		"msgs.msg.packet_id": bson.M{"$exists": false},
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
		"types":              constant.MsgTypeAcknowledgement,
		"msgs.msg.packet_id": bson.M{"$exists": false},
	}).Sort(sorts).Limit(1).One(&tx)
	return tx, err
}

func (repo *TxRepo) UpdateAckPacketId(chainId string, height int64, txHash, packetId string) error {
	filter, update := bson.M{
		"height": height, "tx_hash": txHash,
		"msgs.type": constant.MsgTypeAcknowledgement,
	}, bson.M{
		"$set": bson.M{
			"msgs.$.msg.packet_id": packetId,
		},
	}
	err := repo.coll(chainId).UpdateOne(context.Background(), filter, update)
	return err
}
