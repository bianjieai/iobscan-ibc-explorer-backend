package task

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type IbcDenomUpdateTask struct {
}

var _ Task = new(IbcDenomUpdateTask)

func (t *IbcDenomUpdateTask) Name() string {
	return "ibc_denom_update_task"
}

func (t *IbcDenomUpdateTask) Cron() int {
	if taskConf.CronTimeDenomUpdateTask > 0 {
		return taskConf.CronTimeDenomUpdateTask
	}
	return ThreeMinute
}

func (t *IbcDenomUpdateTask) Run() int {
	denomSymbolMap, err := t.getBaseDenomSysbolMap()
	if err != nil {
		return -1
	}

	if len(denomSymbolMap) == 0 {
		return 1
	}

	if err = t.handleIbcDenoms(denomSymbolMap); err != nil {
		return -1
	}

	return 1
}

func (t *IbcDenomUpdateTask) getBaseDenomSysbolMap() (map[string]string, error) {
	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s baseDenomRepo.FindAll error, %v", t.Name(), err)
		return nil, err
	}

	denomSymbolMap := make(map[string]string, len(baseDenomList))
	for _, v := range baseDenomList {
		denomSymbolMap[fmt.Sprintf("%s%s", v.Chain, v.Denom)] = v.Symbol
	}
	return denomSymbolMap, nil
}

func (t *IbcDenomUpdateTask) handleIbcDenoms(denomSymbolMap map[string]string) error {
	denomList, err := denomRepo.FindNoSymbolDenoms()
	if err != nil {
		logrus.Errorf("task %s denomRepo.FindNoSymbolDenoms error, %v", t.Name(), err)
		return err
	}

	for _, v := range denomList {
		symbol, ok := denomSymbolMap[fmt.Sprintf("%s%s", v.BaseDenomChain, v.BaseDenom)]
		if ok {
			if err = denomRepo.UpdateSymbol(v.Chain, v.Denom, symbol); err != nil {
				logrus.Errorf("task %s UpdateSymbol error, %v", t.Name(), err)
			}
		}
	}

	return nil
}
