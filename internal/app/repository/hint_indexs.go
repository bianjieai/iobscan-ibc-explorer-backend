package repository

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/conf"

const (
	IndexNameMsgsMsgSignerMsgsTypeTime = "msgs.msg.signer_1_msgs.type_1_time_1"
	IndexNameTimeMsgsType              = "time_-1_msgs.type_-1"
	IndexNameCreateAt                  = "create_at_-1"
)

var (
	indexNameConf conf.HintIndexName
)

func LoadIndexNameConf(indexNameCfg conf.HintIndexName) {
	indexNameConf = indexNameCfg
}

func MsgsMsgSignerMsgsTypeTimeIndexName() string {
	if indexNameConf.MsgsMsgSignerMsgsTypeTimeIndexName != "" {
		return indexNameConf.MsgsMsgSignerMsgsTypeTimeIndexName
	}
	return IndexNameMsgsMsgSignerMsgsTypeTime
}

func TimeMsgsTypeIndexName() string {
	if indexNameConf.TimeMsgsTypeIndexName != "" {
		return indexNameConf.TimeMsgsTypeIndexName
	}
	return IndexNameTimeMsgsType
}

func CreateAtIndexName() string {
	if indexNameConf.CreateAtIndexName != "" {
		return indexNameConf.CreateAtIndexName
	}
	return IndexNameCreateAt
}
