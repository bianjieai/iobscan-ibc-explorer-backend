package repository

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"

const (
	IndexNameMsgsMsgSignerMsgsTypeTime = "msgs.msg.signer_1_msgs.type_1_time_1"
	IndexNameTimeMsgsType              = "time_-1_msgs.type_-1"
)

var (
	indexNameConf conf.HintIndexName
)

func LoadIndexNameConf(indexNameCfg conf.HintIndexName) {
	indexNameConf = indexNameCfg
}

func GetLatestRecvPacketTimeHintIndexName() string {
	if indexNameConf.GetLatestRecvPacketTimeHintIndex != "" {
		return indexNameConf.GetLatestRecvPacketTimeHintIndex
	}
	return IndexNameMsgsMsgSignerMsgsTypeTime
}

func GetRelayerUpdateTimeHintIndexName() string {
	if indexNameConf.GetRelayerUpdateTimeHintIndex != "" {
		return indexNameConf.GetLatestRecvPacketTimeHintIndex
	}
	return IndexNameMsgsMsgSignerMsgsTypeTime
}

func CountRelayerTxsHintIndexName() string {
	if indexNameConf.CountRelayerTxsHintIndex != "" {
		return indexNameConf.CountRelayerTxsHintIndex
	}
	return IndexNameMsgsMsgSignerMsgsTypeTime
}

func GetRelayerTxsHintIndexName() string {
	if indexNameConf.GetRelayerTxsHintIndex != "" {
		return indexNameConf.GetRelayerTxsHintIndex
	}
	return IndexNameTimeMsgsType
}
