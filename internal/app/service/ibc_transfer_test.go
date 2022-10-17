package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestTransferService_TraceSource(t *testing.T) {
	data, err := new(TransferService).TraceSource("84CFEBF67B278BE41120F95519E6C96BC41765A5FF5C3C0B272E01CA89B4C4DF",
		&vo.TraceSourceReq{
			ChainId: "irishub_1",
			MsgType: constant.MsgTypeRecvPacket,
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
