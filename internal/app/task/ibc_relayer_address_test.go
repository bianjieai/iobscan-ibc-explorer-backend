package task

import "testing"

func Test_IbcRelayerAddressInitTask(t *testing.T) {
	relayerAddressInitTask.Run()
}

func Test_RelayerAddressGatherTask(t *testing.T) {
	relayerAddressGatherTask.Run()
}
