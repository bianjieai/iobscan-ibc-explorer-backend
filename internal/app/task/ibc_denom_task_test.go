package task

import "testing"

func Test_DenomCalculate(t *testing.T) {
	new(IbcDenomCalculateTask).Run()
}

func Test_IbcHash(t *testing.T) {
	hash := new(IbcDenomCalculateTask).IbcHash("transfer/channel-44", "uiris")
	t.Log(hash)
}

func Test_DenomUpdateTask(t *testing.T) {
	new(IbcDenomUpdateTask).Run()
}
