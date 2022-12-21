package bech32

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGetEscrowAddress(t *testing.T) {
	bz := getEscrowAddress("transfer", "channel-3")
	andEncode, err := ConvertAndEncode("iaa", bz)
	if err != nil {
		panic(err)
	}

	t.Log(andEncode)
}

func getEscrowAddress(portID, channelID string) []byte {
	contents := fmt.Sprintf("%s/%s", portID, channelID)
	// ADR 028 AddressHash construction
	preImage := []byte("ics20-1")
	preImage = append(preImage, 0)
	preImage = append(preImage, contents...)
	hash := sha256.Sum256(preImage)
	return hash[:20]
}

func TestConvert(t *testing.T) {
	dst := "fva"
	bech32Str := "faa17cjdg63thy2vfqvvgj5lfv5dp339t0lr99wc8p"
	wanted := "fva17cjdg63thy2vfqvvgj5lfv5dp339t0lrs5yh6x"
	res := Convert(dst, bech32Str)

	if res != wanted {
		t.Fatal("No Pass")
	}
	t.Log(res)
}

func TestBech32Encode(t *testing.T) {
	val, _ := Bech32Encode("rphr", []byte("rph1xak8zdr8v4e8qemcxy6nwd3sx5mnwdfjg3ahpu"))
	t.Log(val)
	_, data, _ := Bech32Decode(val)
	t.Log(string(data))
}

func TestBech32Decode(t *testing.T) {
	str := "fap1addwnpepqth8487w2wewvnfudrlgcm838a4zu4jwxnumavt0pk4yz78deajekecdzgq"
	if hrp, data, err := Bech32Decode(str); err != nil {
		t.Fatal(err)
	} else {
		t.Log(hrp)
		t.Log(data)
		t.Log(hex.EncodeToString(data))
	}
}

func TestPubKeyToProposerAddrHash(t *testing.T) {
	pubKey := "A4AhOmtQSLxYqsbT/pVT9YoAhoyZZnrvr0EAbsYlVgNY"
	actual, err := PubKeyToProposerAddrHash(pubKey)
	if err != nil {
		t.Fatal(err)
	}

	expected := "940F3F224C42A435327E7057D726E19F8F731252"
	if actual != expected {
		t.Fatalf("not equal, actual:%s, expected:%s", actual, expected)
	}
}
