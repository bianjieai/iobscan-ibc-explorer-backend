package task

import (
	"strings"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
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

	supportDenomList, err := t.getNeedHandleIbcDenoms(denomSymbolMap)
	if err != nil {
		return -1
	}

	baseDenomList, chainDenomsMap := t.collectChainDenomsMap(supportDenomList)
	t.updateBaseDenom(baseDenomList, denomSymbolMap)
	t.updateIbcDenom(chainDenomsMap)

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
		denomSymbolMap[v.Denom] = v.Symbol
	}
	return denomSymbolMap, nil
}

func (t *IbcDenomUpdateTask) getNeedHandleIbcDenoms(denomSymbolMap map[string]string) (entity.IBCDenomList, error) {
	denomList, err := denomRepo.FindNoSymbolDenoms()
	if err != nil {
		logrus.Errorf("task %s denomRepo.FindNoSymbolDenoms error, %v", t.Name(), err)
		return nil, err
	}

	var supportDenomList entity.IBCDenomList
	for _, v := range denomList {
		split := strings.Split(v.DenomPath, "/")
		if len(split) > 2 { // 当前calculate只支持1跳
			continue
		}

		if _, ok := denomSymbolMap[v.BaseDenom]; ok {
			supportDenomList = append(supportDenomList, v)
		}
	}
	return supportDenomList, nil
}

func (t *IbcDenomUpdateTask) collectChainDenomsMap(supportDenomList entity.IBCDenomList) (entity.IBCDenomList, map[string][]string) {
	var baseDenomList entity.IBCDenomList
	chainDenomsMap := make(map[string][]string)
	for _, v := range supportDenomList {
		if !strings.HasPrefix(v.Denom, constant.IBCTokenPreFix) { // base denom
			baseDenomList = append(baseDenomList, v)
		} else {
			chainDenomsMap[v.ChainId] = append(chainDenomsMap[v.ChainId], v.Denom)
		}
	}
	return baseDenomList, chainDenomsMap
}

func (t *IbcDenomUpdateTask) updateBaseDenom(denomList entity.IBCDenomList, denomSymbolMap map[string]string) {
	for _, v := range denomList {
		if err := denomRepo.UpdateSymbol(v.ChainId, v.Denom, denomSymbolMap[v.BaseDenom]); err != nil {
			logrus.Errorf("task %s denomRepo.UpdateSymbol error, %v", t.Name(), err)
		}
	}
}

func (t *IbcDenomUpdateTask) updateIbcDenom(chainDenomsMap map[string][]string) {
	for chainId, denoms := range chainDenomsMap {
		byDenoms, err := denomCalculateRepo.FindByDenoms(chainId, denoms)
		if err != nil {
			logrus.Errorf("task %s denomCalculateRepo.FindByDenoms error, %v", t.Name(), err)
			continue
		}

		for _, v := range byDenoms {
			if err = denomRepo.UpdateSymbol(v.ChainId, v.Denom, v.Symbol); err != nil {
				logrus.Errorf("task %s denomRepo.UpdateSymbol error, %v", t.Name(), err)
			}
		}
	}
}
