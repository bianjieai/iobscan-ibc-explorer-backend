package oneoff

type (
	TxNew struct {
		Chain  string `bson:"chain"`
		Height int64  `bson:"height"`
		TxHash string `bson:"tx_hash"`
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
