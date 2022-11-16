package task

import (
	"context"
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

type ModifyChainIdTask struct {
}

func (t *ModifyChainIdTask) Name() string {
	return "modify_chain_id_task"
}

func (t *ModifyChainIdTask) Switch() bool {
	return false
}

func (t *ModifyChainIdTask) Run() int {
	chainIdNameMap, err := getChainIdNameMap()
	if err != nil {
		logrus.Errorf("task %s getChainIdNameMap err, %v", t.Name(), err)
		return -1
	}
	handlerOne := newModifyChainIdHandlerOne(chainIdNameMap)
	handlerOne.exec("")
	return 1
}

func (t *ModifyChainIdTask) RunWithParam(category string, coll string) int {
	chainIdNameMap, err := getChainIdNameMap()
	if err != nil {
		logrus.Errorf("task %s getChainIdNameMap err, %v", t.Name(), err)
		return -1
	}

	switch category {
	case "one", "1":
		handlerOne := newModifyChainIdHandlerOne(chainIdNameMap)
		handlerOne.exec(coll)

	case "threee", "3":

	}

	return 1
}

// ============================================================================
// ============================================================================
// ============================================================================

// ModifyChainIdHandlerOne 处理第一类问题
type ModifyChainIdHandlerOne struct {
	chainIdNameMap map[string]string
}

func newModifyChainIdHandlerOne(chainIdNameMap map[string]string) *ModifyChainIdHandlerOne {
	return &ModifyChainIdHandlerOne{
		chainIdNameMap: chainIdNameMap,
	}
}

func (h *ModifyChainIdHandlerOne) name() string {
	return "MCH-1"
}

func (h *ModifyChainIdHandlerOne) defaultColls() []string {
	return []string{
		entity.ChainRegistry{}.CollectionName(),
		entity.IBCBaseDenom{}.CollectionName(),
		entity.IBCChain{}.CollectionName(),
		entity.IBCChannelConfig{}.CollectionName(),
		entity.IBCDenom{}.CollectionName(false),
		entity.IBCRelayerAddressChannelCollName,
		entity.IBCRelayerDenomStatisticsCollName,
		entity.IBCRelayerFeeStatisticsCollName,
		entity.IBCToken{}.CollectionName(),
		entity.IBCTokenStatisticsCollName,
		entity.IBCTokenTrace{}.CollectionName(),
		entity.IBCTokenTraceStatisticsCollName,
	}
}

func (h *ModifyChainIdHandlerOne) exec(coll string) {
	if coll == "" {
		h.handleDefault()
	} else {
		h.updateColl(coll)
	}
}

func (h *ModifyChainIdHandlerOne) handleDefault() {
	colls := h.defaultColls()
	collQueue := new(utils.QueueString)
	for _, v := range colls {
		collQueue.Push(v)
	}
	var collsCoordinator *stringQueueCoordinator
	collsCoordinator = &stringQueueCoordinator{
		stringQueue: collQueue,
	}

	const workerNum = 5
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	st := time.Now().Unix()
	for i := 1; i <= workerNum; i++ {
		go func(workName string) {
			defer waitGroup.Done()
			for {
				collName, err := collsCoordinator.getOne()
				if err != nil {
					logrus.Infof("%s %s end", h.name(), workName)
					return
				}

				h.updateColl(collName)
			}
		}(fmt.Sprintf("worker-%d", i))
	}
	waitGroup.Wait()
	logrus.Infof("%s handleDefault end, time use %d(s)", h.name(), time.Now().Unix()-st)
}

func filterBson(k, v string) bson.M {
	return bson.M{
		k: v,
	}
}

func setBson(k, v string) bson.M {
	return bson.M{
		"$set": bson.M{
			k: v,
		},
	}
}

func unsetBson(k string) bson.M {
	return bson.M{
		"$unset": bson.M{
			k: nil,
		},
	}
}

func (h *ModifyChainIdHandlerOne) updateColl(collName string) {
	const (
		chainIdField          = "chain_id"
		chainField            = "chain"
		chainAField           = "chain_a"
		chainBField           = "chain_b"
		baseDenomChainIdField = "base_denom_chain_id"
		baseDenomChainField   = "base_denom_chain"
		prevChainIdField      = "prev_chain_id"
		prevChainField        = "prev_chain"
		statisticsChainField  = "statistics_chain"
	)

	singleChainIdFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(chainIdField, chainId)
			sb := setBson(chainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}

		ub := unsetBson(chainIdField)
		if err := repository.CustomerUpdateAll(collName, context.Background(), bson.M{}, ub); err != nil {
			return err
		}
		return nil
	}

	channelConfigFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(chainAField, chainId)
			sb := setBson(chainAField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}

			fb2 := filterBson(chainBField, chainId)
			sb2 := setBson(chainBField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb2, sb2); err != nil {
				return err
			}
		}
		return nil
	}

	singleBaseDenomChainIdFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(baseDenomChainIdField, chainId)
			sb := setBson(baseDenomChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}

		ub := unsetBson(baseDenomChainIdField)
		if err := repository.CustomerUpdateAll(collName, context.Background(), bson.M{}, ub); err != nil {
			return err
		}
		return nil
	}

	singlePrevChainIdFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(prevChainIdField, chainId)
			sb := setBson(prevChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}

		ub := unsetBson(prevChainIdField)
		if err := repository.CustomerUpdateAll(collName, context.Background(), bson.M{}, ub); err != nil {
			return err
		}
		return nil
	}

	statisticsChainFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(statisticsChainField, chainId)
			sb := setBson(statisticsChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}
		return nil
	}

	chainFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(chainField, chainId)
			sb := setBson(chainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}
		return nil
	}

	baseDenomChainFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterBson(baseDenomChainField, chainId)
			sb := setBson(baseDenomChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}
		return nil
	}

	var err error
	st := time.Now().Unix()
	logrus.Infof("%s update coll %s start", h.name(), collName)
	switch collName {
	case entity.ChainRegistry{}.CollectionName():
		err = singleChainIdFunc()
	case entity.IBCBaseDenom{}.CollectionName():
		err = singleChainIdFunc()
	case entity.IBCChain{}.CollectionName():
		err = singleChainIdFunc()
	case entity.IBCChannelConfig{}.CollectionName():
		err = channelConfigFunc()
	case entity.IBCDenom{}.CollectionName(false):
		if err = singleChainIdFunc(); err != nil {
			break
		}
		if err = singleBaseDenomChainIdFunc(); err != nil {
			break
		}
		if err = singlePrevChainIdFunc(); err != nil {
			break
		}
	case entity.IBCRelayerAddressChannelCollName:
		err = chainFunc()
	case entity.IBCRelayerDenomStatisticsCollName:
		if err = statisticsChainFunc(); err != nil {
			break
		}
		if err = baseDenomChainFunc(); err != nil {
			break
		}

	case entity.IBCRelayerFeeStatisticsCollName:
		err = statisticsChainFunc()
	case entity.IBCToken{}.CollectionName():
		err = singleChainIdFunc()
	case entity.IBCTokenStatisticsCollName:
		err = singleBaseDenomChainIdFunc()
	case entity.IBCTokenTrace{}.CollectionName():
		if err = singleChainIdFunc(); err != nil {
			break
		}
		if err = singleBaseDenomChainIdFunc(); err != nil {
			break
		}
	case entity.IBCTokenTraceStatisticsCollName:
		err = singleChainIdFunc()
	}

	if err != nil {
		logrus.Errorf("%s update coll %s err, %v", h.name(), collName, err)
	} else {
		logrus.Infof("%s update coll %s end, time use %d(s)", h.name(), collName, time.Now().Unix()-st)
	}
}
