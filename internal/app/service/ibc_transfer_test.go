package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestTransferService_TraceSource(t *testing.T) {
	data, err := new(TransferService).TraceSource("87DD9D44F64EC8E509508B99AD48554F9FCD3A79D775A400FE900CCA030290BE",
		&vo.TraceSourceReq{
			ChainId: "laozi_mainnet",
			MsgType: constant.MsgTypeTransfer,
		})
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
}

func TestGetTxDataFromChain(t *testing.T) {
	data, err := GetTxDataFromChain("https://mainnet.crescent.network:1317",
		"0E000429F0CCB543D0FE0CDA57DF3A470E8DE54498FF071E755736CDBECE1C72")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
}
