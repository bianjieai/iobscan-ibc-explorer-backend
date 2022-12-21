package bech32

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	ctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	enccodec "github.com/tharsis/ethermint/encoding/codec"
)

func GetProtoCodec() *codec.ProtoCodec {
	//var cdc = codec.NewLegacyAmino()
	//cryptocodec.RegisterCrypto(cdc)
	interfaceRegistry := ctypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	enccodec.RegisterInterfaces(interfaceRegistry)
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
