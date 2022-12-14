package service

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"testing"
)

func TestAddressService_TokenList(t *testing.T) {
	data, err := new(AddressService).TokenList("irisnet", "iaa1z2sdef0ypat9lq7wsxrt7ue3uzdnzcsd34wsl4")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(utils.MarshalJsonIgnoreErr(data)))
}
