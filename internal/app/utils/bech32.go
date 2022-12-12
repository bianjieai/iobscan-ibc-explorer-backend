package utils

/***
code from page detail: https://pkg.go.dev/github.com/btcsuite/btcutil/bech32
*/

import (
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/codec"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("%v encoding bech32 failed", err.Error())
	}
	return bech32.Encode(hrp, converted)

}

//DecodeAndConvert decodes a bech32 encoded string and converts to base64 encoded bytes
func DecodeAndConvert(bech string) (string, []byte, error) {
	hrp, data, err := bech32.Decode(bech)
	if err != nil {
		return "", nil, fmt.Errorf("%v decoding bech32 failed", err.Error())
	}
	converted, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, fmt.Errorf("%v decoding bech32 failed", err.Error())
	}
	return hrp, converted, nil
}

func Bech32Encode(hrp string, data []byte) (string, error) {
	if regrouped, err := bech32.ConvertBits(data, 8, 5, true); err != nil {
		return "", err
	} else {
		return bech32.Encode(hrp, regrouped)
	}
}

func Bech32Decode(bechStr string) (string, []byte, error) {
	if hrp, regrouped, err := bech32.Decode(bechStr); err != nil {
		return "", nil, err
	} else {
		if data, err := bech32.ConvertBits(regrouped, 5, 8, false); err != nil {
			return hrp, nil, err
		} else {
			return hrp, data, nil
		}
	}
}

func GetProtoCodec() *codec.ProtoCodec {
	//var cdc = codec.NewLegacyAmino()
	//cryptocodec.RegisterCrypto(cdc)
	interfaceRegistry := ctypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	//keyring.RegisterLegacyAminoCodec(cdc)
	//cryptocodec.RegisterCrypto(cdc)
	return marshaler
}
func GetAddressFromPubkey(addPrefix string, jsonPubKeyData string) (string, error) {
	var acc types.BaseAccount
	protoC := GetProtoCodec()
	err := protoC.UnmarshalJSON([]byte(jsonPubKeyData), &acc)
	if err != nil {
		return "", err
	}
	var pubKey cryptotypes.PubKey
	if acc.PubKey == nil {
		return "", fmt.Errorf("acc.Pubkey is nil")
	}
	if err := protoC.UnpackAny(acc.PubKey, &pubKey); err != nil {
		return "", err
	}
	addr, err := ConvertAndEncode(addPrefix, pubKey.Address())
	if err != nil {
		return "", err
	}
	return addr, nil
}
