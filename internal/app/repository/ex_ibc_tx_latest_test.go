package repository

import (
	"encoding/json"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"testing"
)

func TestExIbcTxRepo_FindAllByStatus(t *testing.T) {
	status := []entity.IbcTxStatus{1, 2}
	data, err := new(ExIbcTxRepo).FindAllByStatus(status, 0, 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret, _ := json.Marshal(data)
	t.Log(string(ret))
}
