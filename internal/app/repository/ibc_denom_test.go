package repository

import (
	"testing"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
)

func TestDenomRepo_InsertBatch(t *testing.T) {
	denoms := []*entity.IBCDenom{
		{
			Symbol:           "IRIS",
			ChainId:          "bigbang",
			Denom:            "uiris",
			PrevDenom:        "",
			PrevChainId:      "",
			BaseDenom:        "uiris",
			BaseDenomChainId: "bigbang",
			DenomPath:        "",
			IsSourceChain:    false,
			IsBaseDenom:      true,
			CreateAt:         time.Now().Unix(),
			UpdateAt:         time.Now().Unix(),
		},
		{
			Symbol:           "IRIS",
			ChainId:          "bigbang",
			Denom:            "test0909",
			PrevDenom:        "",
			PrevChainId:      "",
			BaseDenom:        "uiris",
			BaseDenomChainId: "bigbang",
			DenomPath:        "",
			IsSourceChain:    false,
			IsBaseDenom:      true,
			CreateAt:         time.Now().Unix(),
			UpdateAt:         time.Now().Unix(),
		},
	}
	err := new(DenomRepo).InsertBatch(denoms)
	if err != nil {
		t.Fatal(err)
	}
}
