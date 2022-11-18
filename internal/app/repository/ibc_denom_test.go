package repository

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

func TestDenomRepo_InsertBatch(t *testing.T) {
	denoms := []*entity.IBCDenom{
		{
			Symbol:         "IRIS",
			Chain:          "bigbang",
			Denom:          "uiris",
			PrevDenom:      "",
			PrevChain:      "",
			BaseDenom:      "uiris",
			BaseDenomChain: "bigbang",
			DenomPath:      "",
			IsBaseDenom:    true,
			CreateAt:       time.Now().Unix(),
			UpdateAt:       time.Now().Unix(),
		},
		{
			Symbol:         "IRIS",
			Chain:          "bigbang",
			Denom:          "test0909",
			PrevDenom:      "",
			PrevChain:      "",
			BaseDenom:      "uiris",
			BaseDenomChain: "bigbang",
			DenomPath:      "",
			IsBaseDenom:    true,
			CreateAt:       time.Now().Unix(),
			UpdateAt:       time.Now().Unix(),
		},
	}
	err := new(DenomRepo).InsertBatch(denoms)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDenomRepo_GetBaseDenomNoSymbol(t *testing.T) {
	res, err := new(DenomRepo).GetBaseDenomNoSymbol()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(res)))
}
