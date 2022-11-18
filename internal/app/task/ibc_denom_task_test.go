package task

import "testing"

// TODO remove this

func Test_DenomCalculate(t *testing.T) {
	new(IbcDenomCalculateTask).Run()
}

func Test_DenomUpdateTask(t *testing.T) {
	new(IbcDenomUpdateTask).Run()
}
