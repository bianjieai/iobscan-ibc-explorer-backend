package entity

import (
	"fmt"
)

type SyncBlock struct {
	Height   int64  `bson:"height"`
	Hash     string `bson:"hash"`
	Txn      int64  `bson:"txn"`
	Time     int64  `bson:"time"`
	Proposer string `bson:"proposer"`
}

func (s SyncBlock) CollectionName(chain string) string {
	return fmt.Sprintf("sync_%s_block", chain)
}
