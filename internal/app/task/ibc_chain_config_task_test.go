package task

import "testing"

func Test_ibcChainConfigTask(t *testing.T) {
	_ibcChainConfigTask.Run()
}

func TestEmptyTxs(t *testing.T) {
	cfgs, err := chainConfigRepo.FindAll()
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, val := range cfgs {
		_, err := txRepo.FindHeight(val.CurrentChainId, true)
		if err != nil {
			t.Errorf(err.Error(), val.CurrentChainId)
		}
	}
}
