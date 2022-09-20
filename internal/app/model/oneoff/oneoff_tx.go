package oneoff

import (
	"go.mongodb.org/mongo-driver/bson"
)

type (
	TxNew struct {
		ChainId string `bson:"chain_id"`
		Height  int64  `bson:"height"`
		TxHash  string `bson:"tx_hash"`
	}
)

/***
db.sync_tx_new.createIndex({
"height": -1,
"chain_id": -1,
"tx_hash": -1
}, {background: true, unique: true});
*/
func (i TxNew) CollectionName() string {
	return "sync_tx_new"
}

func (i TxNew) PkKvPair() map[string]interface{} {
	return bson.M{"chain_id": i.ChainId, "tx_hash": i.TxHash}
}
