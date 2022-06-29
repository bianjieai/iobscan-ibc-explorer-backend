package repository

import (
	"context"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type IExIbcTxRepo interface {
	FindAll(skip, limit int64) ([]*entity.ExIbcTx, error)
	FindAllHistory(skip, limit int64) ([]*entity.ExIbcTx, error)
	CountBaseDenomTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error)
	CountBaseDenomHistoryTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error)
	CountIBCTokenRecvTxs(baseDenom, chainId string) ([]*dto.CountIBCTokenRecvTxsDTO, error)
	CountIBCTokenHistoryRecvTxs(baseDenom, chainId string) ([]*dto.CountIBCTokenRecvTxsDTO, error)
	GetRelayerInfo(latestTxTime int64) ([]*dto.GetRelayerInfoDTO, error)
	GetHistoryRelayerInfo(latestTxTime int64) ([]*dto.GetRelayerInfoDTO, error)
	GetLatestTxTime() (int64, error)
	GetOneRelayerScTxPacketId(dto *dto.GetRelayerInfoDTO) (entity.ExIbcTx, error)
	GetHistoryOneRelayerScTxPacketId(dto *dto.GetRelayerInfoDTO) (entity.ExIbcTx, error)
	CountHistoryRelayerSuccessPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error)
	CountRelayerSuccessPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error)
	CountHistoryRelayerPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error)
	CountRelayerPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error)
	CountHistoryRelayerPacketAmount() ([]*dto.CountRelayerPacketAmountDTO, error)
	CountRelayerPacketAmount() ([]*dto.CountRelayerPacketAmountDTO, error)
	AggrIBCChannelTxs() ([]*dto.AggrIBCChannelTxsDTO, error)
	AggrIBCChannelHistoryTxs() ([]*dto.AggrIBCChannelTxsDTO, error)
}

var _ IExIbcTxRepo = new(ExIbcTxRepo)

type ExIbcTxRepo struct {
}

func (repo *ExIbcTxRepo) coll() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExIbcTx{}.CollectionName(false))
}

func (repo *ExIbcTxRepo) collHistory() *qmgo.Collection {
	return mgo.Database(ibcDatabase).Collection(entity.ExIbcTx{}.CollectionName(true))
}

func (repo *ExIbcTxRepo) FindAll(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.coll().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) FindAllHistory(skip, limit int64) ([]*entity.ExIbcTx, error) {
	var res []*entity.ExIbcTx
	err := repo.collHistory().Find(context.Background(), bson.M{}).Skip(skip).Limit(limit).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) countBaseDenomTransferTxsPipe() []bson.M {
	match := bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": entity.IbcTxUsefulStatus,
			},
		},
	}

	group := bson.M{
		"$group": bson.M{
			"_id": "$base_denom",
			"count": bson.M{
				"$sum": 1,
			},
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group)
	return pipe
}

func (repo *ExIbcTxRepo) CountBaseDenomTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error) {
	pipe := repo.countBaseDenomTransferTxsPipe()
	var res []*dto.CountBaseDenomTransferAmountDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountBaseDenomHistoryTransferTxs() ([]*dto.CountBaseDenomTransferAmountDTO, error) {
	pipe := repo.countBaseDenomTransferTxsPipe()
	var res []*dto.CountBaseDenomTransferAmountDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) countIBCTokenRecvTxsPipe(baseDenom, chainId string) []bson.M {
	match := bson.M{
		"$match": bson.M{
			"denoms.sc_denom": baseDenom,
			"sc_chain_id":     chainId,
			"status":          entity.IbcTxStatusSuccess,
		},
	}

	group := bson.M{
		"$group": bson.M{
			"_id": "$denoms.dc_denom",
			"count": bson.M{
				"$sum": 1,
			},
		},
	}

	var pipe []bson.M
	pipe = append(pipe, match, group)
	return pipe
}

func (repo *ExIbcTxRepo) CountIBCTokenRecvTxs(baseDenom, chainId string) ([]*dto.CountIBCTokenRecvTxsDTO, error) {
	pipe := repo.countIBCTokenRecvTxsPipe(baseDenom, chainId)
	var res []*dto.CountIBCTokenRecvTxsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountIBCTokenHistoryRecvTxs(baseDenom, chainId string) ([]*dto.CountIBCTokenRecvTxsDTO, error) {
	pipe := repo.countIBCTokenRecvTxsPipe(baseDenom, chainId)
	var res []*dto.CountIBCTokenRecvTxsDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) GetRelayerInfo(latestTxTime int64) ([]*dto.GetRelayerInfoDTO, error) {
	pipe := repo.relayerInfoPipe(latestTxTime)
	var res []*dto.GetRelayerInfoDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) GetHistoryRelayerInfo(latestTxTime int64) ([]*dto.GetRelayerInfoDTO, error) {
	pipe := repo.relayerInfoPipe(latestTxTime)
	var res []*dto.GetRelayerInfoDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) relayerInfoPipe(txTime int64) []bson.M {
	match := bson.M{
		"$match": bson.M{
			"dc_tx_info.status": 1,
			"tx_time": bson.M{
				"$gte": txTime,
			},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"relayer":     "$dc_tx_info.msg.msg.signer",
				"sc_chain_id": "$sc_chain_id",
				"sc_channel":  "$sc_channel",
				"dc_chain_id": "$dc_chain_id",
				"dc_channel":  "$dc_channel",
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":              0,
			"dc_chain_address": "$_id.relayer",
			"sc_chain_id":      "$_id.sc_chain_id",
			"dc_chain_id":      "$_id.dc_chain_id",
			"sc_channel":       "$_id.sc_channel",
			"dc_channel":       "$_id.dc_channel",
		},
	}
	sort := bson.M{
		"$sort": bson.M{
			"tx_time": 1,
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project, sort)
	return pipe

}

func (repo *ExIbcTxRepo) GetLatestTxTime() (int64, error) {
	var res *entity.ExIbcTx
	err := repo.coll().Find(context.Background(), bson.M{}).Select(bson.M{"tx_time": 1}).Sort("-tx_time").One(&res)
	if err != nil {
		return 0, err
	}
	return res.TxTime, nil
}

func (repo *ExIbcTxRepo) oneRelayerPacketCond(relayer *dto.GetRelayerInfoDTO) bson.M {
	return bson.M{
		"dc_tx_info.msg.msg.signer": relayer.DcChainAddress,
		"sc_chain_id":               relayer.ScChainId,
		"dc_chain_id":               relayer.DcChainId,
		"sc_channel":                relayer.ScChannel,
		"dc_channel":                relayer.DcChannel,
	}
}

func (repo *ExIbcTxRepo) GetOneRelayerScTxPacketId(dto *dto.GetRelayerInfoDTO) (entity.ExIbcTx, error) {
	var res entity.ExIbcTx
	err := repo.coll().Find(context.Background(), repo.oneRelayerPacketCond(dto)).
		Select(bson.M{"sc_tx_info.msg.msg.packet_id": 1}).Sort("-tx_time").One(&res)
	return res, err
}

func (repo *ExIbcTxRepo) GetHistoryOneRelayerScTxPacketId(dto *dto.GetRelayerInfoDTO) (entity.ExIbcTx, error) {
	var res entity.ExIbcTx
	err := repo.collHistory().Find(context.Background(), repo.oneRelayerPacketCond(dto)).
		Select(bson.M{"sc_tx_info.msg.msg.packet_id": 1}).Sort("-tx_time").One(&res)
	return res, err
}

func (repo *ExIbcTxRepo) relayerSuccessPacketCond() []bson.M {
	match := bson.M{
		"$match": bson.M{
			"status": entity.IbcTxStatusSuccess,
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"dc_chain_id": "$dc_chain_id",
				"dc_channel":  "$dc_channel",
				"relayer":     "$dc_tx_info.msg.msg.signer",
			},
			"count": bson.M{
				"$sum": 1,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":              0,
			"dc_chain_address": "$_id.relayer",
			"dc_chain_id":      "$_id.dc_chain_id",
			"dc_channel":       "$_id.dc_channel",
			"count":            "$count",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	return pipe
}

func (repo *ExIbcTxRepo) relayerPacketCond() []bson.M {
	match := bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": entity.IbcTxUsefulStatus,
			},
			"sc_tx_info.status": entity.TxStatusSuccess,
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"dc_chain_id": "$dc_chain_id",
				"dc_channel":  "$dc_channel",
				"relayer":     "$dc_tx_info.msg.msg.signer",
			},
			"count": bson.M{
				"$sum": 1,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":              0,
			"dc_chain_address": "$_id.relayer",
			"dc_chain_id":      "$_id.dc_chain_id",
			"dc_channel":       "$_id.dc_channel",
			"count":            "$count",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	return pipe
}

func (repo *ExIbcTxRepo) relayerPacketAmountCond() []bson.M {
	match := bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": entity.IbcTxUsefulStatus,
			},
			"sc_tx_info.status": entity.TxStatusSuccess,
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"dc_chain_id": "$dc_chain_id",
				"dc_channel":  "$dc_channel",
				"relayer":     "$dc_tx_info.msg.msg.signer",
				"base_denom":  "$base_denom",
			},
			"amount": bson.M{
				"$sum": bson.M{"$toDouble": "$sc_tx_info.msg_amount.amount"},
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":              0,
			"dc_chain_address": "$_id.relayer",
			"dc_chain_id":      "$_id.dc_chain_id",
			"dc_channel":       "$_id.dc_channel",
			"base_denom":       "$_id.base_denom",
			"amount":           "$amount",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	return pipe
}

func (repo *ExIbcTxRepo) CountHistoryRelayerSuccessPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error) {
	pipe := repo.relayerSuccessPacketCond()
	var res []*dto.CountRelayerPacketTxsCntDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountRelayerSuccessPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error) {
	pipe := repo.relayerSuccessPacketCond()
	var res []*dto.CountRelayerPacketTxsCntDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountHistoryRelayerPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error) {
	pipe := repo.relayerPacketCond()
	var res []*dto.CountRelayerPacketTxsCntDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountRelayerPacketTxs() ([]*dto.CountRelayerPacketTxsCntDTO, error) {
	pipe := repo.relayerPacketCond()
	var res []*dto.CountRelayerPacketTxsCntDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountRelayerPacketAmount() ([]*dto.CountRelayerPacketAmountDTO, error) {
	pipe := repo.relayerPacketAmountCond()
	var res []*dto.CountRelayerPacketAmountDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) CountHistoryRelayerPacketAmount() ([]*dto.CountRelayerPacketAmountDTO, error) {
	pipe := repo.relayerPacketAmountCond()
	var res []*dto.CountRelayerPacketAmountDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) AggrIBCChannelTxsPipe() []bson.M {
	match := bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": entity.IbcTxUsefulStatus,
			},
		},
	}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"base_denom":  "$base_denom",
				"sc_chain_id": "$sc_chain_id",
				"dc_chain_id": "$dc_chain_id",
				"sc_channel":  "$sc_channel",
				"dc_channel":  "$dc_channel",
			},
			"count": bson.M{
				"$sum": 1,
			},
			"amount": bson.M{
				"$sum": bson.M{
					"$toDouble": "$sc_tx_info.msg_amount.amount",
				},
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":         0,
			"base_denom":  "$_id.base_denom",
			"sc_chain_id": "$_id.sc_chain_id",
			"dc_chain_id": "$_id.dc_chain_id",
			"sc_channel":  "$_id.sc_channel",
			"dc_channel":  "$_id.dc_channel",
			"count":       "$count",
			"amount":      "$amount",
		},
	}
	var pipe []bson.M
	pipe = append(pipe, match, group, project)
	return pipe
}

func (repo *ExIbcTxRepo) AggrIBCChannelTxs() ([]*dto.AggrIBCChannelTxsDTO, error) {
	pipe := repo.AggrIBCChannelTxsPipe()
	var res []*dto.AggrIBCChannelTxsDTO
	err := repo.coll().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}

func (repo *ExIbcTxRepo) AggrIBCChannelHistoryTxs() ([]*dto.AggrIBCChannelTxsDTO, error) {
	pipe := repo.AggrIBCChannelTxsPipe()
	var res []*dto.AggrIBCChannelTxsDTO
	err := repo.collHistory().Aggregate(context.Background(), pipe).All(&res)
	return res, err
}
