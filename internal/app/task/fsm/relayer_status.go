package fsm

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/looplab/fsm"
)

type (
	ibcRelayerStatsFSMAction struct {
		relayer repository.IbcRelayerRepo
	}
)

const (
	IbcRelayerEventUnknown = "Unknown"
	IbcRelayerEventRunning = "Running"
)

func NewIbcRelayerFSM(initial string) *fsm.FSM {
	events := []fsm.EventDesc{
		{
			Name: IbcRelayerEventUnknown,
			Src:  []string{entity.RelayerStopStr},
			Dst:  entity.RelayerRunningStr,
		},
		{
			Name: IbcRelayerEventRunning,
			Src:  []string{entity.RelayerRunningStr},
			Dst:  entity.RelayerStopStr,
		},
	}

	action := ibcRelayerStatsFSMAction{}
	callbacks := fsm.Callbacks{
		IbcRelayerEventUnknown: action.changeState,
		IbcRelayerEventRunning: action.changeState,
	}
	f := fsm.NewFSM(initial, events, callbacks)

	return f
}

func (action *ibcRelayerStatsFSMAction) changeState(e *fsm.Event) {
	args := e.Args
	if len(args) != 1 {
		e.Err = fmt.Errorf("num of args must be 1")
		return
	}
	data, ok := args[0].(entity.IBCRelayer)
	if !ok {
		e.Err = fmt.Errorf("convert arg[0] to db ibc_relayer fail: %s", utils.MarshalJsonIgnoreErr(args))
		return
	}
	if err := action.relayer.UpdateStatusAndTime(data.RelayerId, int(data.Status), data.UpdateTime, data.TimePeriod); err != nil {
		e.Err = fmt.Errorf("update ibc_relayer status failed: %s", err.Error())
		return
	}
}
