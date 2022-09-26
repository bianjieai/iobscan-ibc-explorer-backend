package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type IbcDenomCalculateTask struct {
}

var _ Task = new(IbcDenomCalculateTask)

func (t *IbcDenomCalculateTask) Name() string {
	return "ibc_denom_calculate_task"
}

func (t *IbcDenomCalculateTask) Cron() int {
	if taskConf.CronTimeDenomCalculateTask > 0 {
		return taskConf.CronTimeDenomCalculateTask
	}
	return EveryMinute
}

func (t *IbcDenomCalculateTask) Run() int {
	chainConfMap, err := t.getChainConfigMap()
	if err != nil {
		return -1
	}

	baseDenomList, err := baseDenomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s  baseDenomRepo.FindAll  error, %v", t.Name(), err)
		return -1
	}

	for _, v := range baseDenomList {
		chainConf, ok := chainConfMap[v.ChainId]
		if ok {
			_ = t.calculateDenom(v, chainConf)
		}
	}

	return 1
}

func (t *IbcDenomCalculateTask) getChainConfigMap() (map[string]*entity.ChainConfig, error) {
	confList, err := chainConfigRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s  chainConfigRepo.FindAll  error, %v", t.Name(), err)
		return nil, err
	}

	res := make(map[string]*entity.ChainConfig, len(confList))
	for _, v := range confList {
		res[v.ChainId] = v
	}

	return res, nil
}

// calculateDenom 计算一跳的ibc hash denom
func (t *IbcDenomCalculateTask) calculateDenom(baseDenom *entity.IBCBaseDenom, chainConf *entity.ChainConfig) error {
	if len(chainConf.IbcInfo) == 0 {
		return nil
	}

	hashCode := utils.Md5(utils.MustMarshalJsonToStr(chainConf.IbcInfo))
	if hashCode == baseDenom.IbcInfoHashCaculate {
		return nil
	}

	existedDenomMap, err := t.getExistedDenom(baseDenom.ChainId)
	if err != nil {
		return err
	}

	var newDenomList []*entity.IBCDenomCalculate
	newDenomMap := make(map[string]string)
	for _, ibcInfo := range chainConf.IbcInfo {
		for _, path := range ibcInfo.Paths {
			if path.ChainId == "" || path.Counterparty.ChannelId == "" || path.Counterparty.PortId == "" {
				continue
			}

			denomPath := fmt.Sprintf("%s/%s", path.Counterparty.PortId, path.Counterparty.ChannelId)
			existKey := fmt.Sprintf("%s/%s/%s", path.ChainId, denomPath, baseDenom.Denom)
			_, ok1 := existedDenomMap[existKey]
			_, ok2 := newDenomMap[existKey]
			if ok1 || ok2 {
				continue
			}

			ibcHash := t.IbcHash(denomPath, baseDenom.Denom)
			newDenomList = append(newDenomList, &entity.IBCDenomCalculate{
				Symbol:    baseDenom.Symbol,
				BaseDenom: baseDenom.Denom,
				Denom:     ibcHash,
				DenomPath: denomPath,
				ChainId:   path.ChainId,
				ScChainId: baseDenom.ChainId,
				CreateAt:  time.Now().Unix(),
				UpdateAt:  time.Now().Unix(),
			})
			newDenomMap[existKey] = ""
		}
	}

	if err = denomCalculateRepo.InsertTransaction(newDenomList, baseDenom.Denom, baseDenom.ChainId, hashCode); err != nil {
		logrus.Errorf("task %s InsertTransaction error, %v", t.Name(), err)
		return err
	}

	return nil
}

func (t *IbcDenomCalculateTask) getExistedDenom(scChainId string) (map[string]string, error) {
	denomList, err := denomCalculateRepo.FindByScChainId(scChainId)
	if err != nil {
		logrus.Errorf("task %s getExistedDenom error, %v", t.Name(), err)
		return nil, err
	}

	denomMap := make(map[string]string, len(denomList))
	for _, v := range denomList {
		denomMap[fmt.Sprintf("%s/%s/%s", v.ChainId, v.DenomPath, v.BaseDenom)] = ""
	}

	return denomMap, nil
}

func (t *IbcDenomCalculateTask) IbcHash(denomPath, baseDenom string) string {
	if denomPath == "" {
		return baseDenom
	}

	hash := utils.Sha256(fmt.Sprintf("%s/%s", denomPath, baseDenom))
	return fmt.Sprintf("%s/%s", constant.IBCTokenPrefix, strings.ToUpper(hash))
}
