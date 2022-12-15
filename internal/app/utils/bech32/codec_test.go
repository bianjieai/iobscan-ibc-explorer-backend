package bech32

import (
	"fmt"
	"testing"
)

func TestAddrConvert(t *testing.T) {
	bech32Str := "cosmos1gu65xf8eluc0g7v84qudhulalam5cm7wzv3ecw"

	_, data, err := DecodeAndConvert(bech32Str)
	if err != nil {
		t.Fatal(err.Error())
	}

	prefixs := []string{
		"akash",
		"axelar",
		"band",
		"bcna",
		"bitsong",
		"bostrom",
		"cerberus",
		"certik",
		"cheqd",
		"chihuahua",
		"comdex",
		"cosmos",
		"crc",
		"cre",
		"cro",
		"cudos",
		"darc",
		"decentr",
		"desmos",
		"dig",
		"emoney",
		"evmos",
		"ex",
		"fetch",
		"gravity",
		"iaa",
		"inj",
		"ixo",
		"juno",
		"kava",
		"ki",
		"kujira",
		"like",
		"lum",
		"mantle",
		"micro",
		"omniflix",
		"osmo",
		"panacea",
		"pasg",
		"pb",
		"persistence",
		"regen",
		"rizon",
		"secret",
		"sent",
		"sif",
		"somm",
		"stafi",
		"star",
		"stars",
		"stride",
		"terra",
		"tgrade",
		"umee",
		"vdl",
	}
	for _, addrPrefix := range prefixs {
		dstAddr, err := ConvertAndEncode(addrPrefix, data)
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Println(dstAddr)
	}

}

func TestGetAddressFromPubkey(t *testing.T) {
	jsonPubKeyData := "{\"pub_key\":{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"AoCSeYV5N3anty1GZv7gLAr0aQ+yaIz62n+bIAgBQvoH\"}}"
	addr, err := GetAddressFromPubkey("kava", jsonPubKeyData)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(addr)
}

func TestGetAddressFromEthPubkey(t *testing.T) {
	jsonPubKeyData := "{\"pub_key\":{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"Anv7L1WQFtJDEVSNrEgjUvfonk2AsXFbVlr2E+pbVqmW\"}}"
	addr, err := GetAddressFromPubkey("evmos", jsonPubKeyData)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(addr)
}
