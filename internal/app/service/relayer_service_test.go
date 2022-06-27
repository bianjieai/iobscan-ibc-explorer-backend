package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"testing"
)

func TestRelayerService_List(t *testing.T) {
	resp, err := new(RelayerService).List(&vo.RelayerListReq{
		Page: vo.Page{
			PageNum:  1,
			PageSize: 10,
		},
		UseCount: false,
		Status:   2,
		Chain:    constant.AllChain,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(resp)
}
