package repository

import (
	"fmt"
	"github.com/qiniu/qmgo"
	"testing"
)

func TestIbcTokenTraceRepo_UpdateBaseDenomAndChain(t *testing.T) {
	err := new(TokenTraceRepo).UpdateBaseDenomAndChain("akashere", "ibc/2B406972A01294D27DE752EFD27DB404F0BAA984DD51349BD22D8B7B3A492254", "ibc/F7AA8487A9FACE44A46AC1CD9A668EA10F9F16389F8591555C674FCD8CF70D5F_test", "cosmoshub_test")
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			fmt.Println("no err")
		}
		t.Fatal(err.Error())
	}
}
