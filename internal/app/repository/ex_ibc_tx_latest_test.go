package repository

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestExIbcTxRepo_AggrIBCChainInflow(t *testing.T) {
	data, err := new(ExIbcTxRepo).AggrIBCChainInflow(1672386690, 1672386691, false)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}

func TestExIbcTxRepo_FindAndUpdate(t *testing.T) {
	repo := new(ExIbcTxRepo)
	data, err := repo.FindProcessingTxs("bigbangname", 1)
	if err != nil {
		t.Fatal(err)
	}

	tx := data[0]
	t.Log(fmt.Sprintf("_id %s", tx.Id))
	set := bson.M{
		"$set": bson.M{
			"update_at":    time.Now().Unix(),
			"process_info": "Processing",
		},
	}
	if err = repo.UpdateOne(tx.Id, false, set); err != nil {
		t.Fatal(err)
	}
}
