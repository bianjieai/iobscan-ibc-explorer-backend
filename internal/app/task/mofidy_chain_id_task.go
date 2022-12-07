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
	chainIdNameMap, err := repository.GetChainIdNameMap()
	if err != nil {
		logrus.Errorf("task %s getChainIdNameMap err, %v", t.Name(), err)
		return -1
	}
	handlerOne := newModifyChainIdHandlerOne(chainIdNameMap)
	handlerOne.exec("")
	return 1
}

func (t *ModifyChainIdTask) RunWithParam(category string, coll string) int {
	chainIdNameMap, err := repository.GetChainIdNameMap()
	if err != nil {
		logrus.Errorf("task %s getChainIdNameMap err, %v", t.Name(), err)
		return -1
	}

	switch category {
	case "one", "1":
		handlerOne := newModifyChainIdHandlerOne(chainIdNameMap)
		handlerOne.exec(coll)

	case "two", "2":
		handlerTwo := newModifyChainIdHandlerTwo(coll, chainIdNameMap)
		handlerTwo.Run()

	case "three", "3":
		handlerThree := newModifyChainIdHandlerThree(chainIdNameMap)
		handlerThree.exec()
	}
	return 1
}

// ============================================================================
// ============================================================================
// ============================================================================
// bson

func filterBson(k, v string) bson.M {
	return bson.M{
		k: v,
	}
}

func filterSegmentBson(k, v string, segKey string, seg *segment) bson.M {
	return bson.M{
		k: v,
		segKey: bson.M{
			"$gte": seg.StartTime,
			"$lte": seg.EndTime,
		},
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
	return "MCHOne"
}

func (h *ModifyChainIdHandlerOne) defaultColls() []string {
	return []string{
		entity.ChainRegistry{}.CollectionName(),
		entity.AuthDenom{}.CollectionName(),
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
	case entity.AuthDenom{}.CollectionName():
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

// ============================================================================
// ============================================================================
// ============================================================================

// ModifyChainIdHandlerThree 处理第一类问题
type ModifyChainIdHandlerThree struct {
	chainIdNameMap map[string]string
}

func newModifyChainIdHandlerThree(chainIdNameMap map[string]string) *ModifyChainIdHandlerThree {
	return &ModifyChainIdHandlerThree{
		chainIdNameMap: chainIdNameMap,
	}
}

func (h *ModifyChainIdHandlerThree) name() string {
	return "MCHThree"
}

func (h *ModifyChainIdHandlerThree) getTxSegments(isTargetHistory bool) ([]*segment, error) {
	var startTime int64

	if isTargetHistory {
		first, err := ibcTxRepo.FirstHistory()
		if err != nil {
			return nil, err
		}
		startTime = first.CreateAt
	} else {
		first, err := ibcTxRepo.First()
		if err != nil {
			return nil, err
		}
		startTime = first.CreateAt
	}

	segs := segmentTool(3600*24*2, startTime, time.Now().Unix())
	return segs, nil
}

func (h *ModifyChainIdHandlerThree) exec() {
	// ex_ibc_tx
	segments, err := h.getTxSegments(true)
	if err != nil {
		logrus.Errorf("%s get tx segments err, %v", h.name(), err)
		return
	}

	doHandleSegments(h.name(), 16, segments, true, h.updateColl)

	// ex_ibc_tx_latest
	latestSegments, err := h.getTxSegments(false)
	if err != nil {
		logrus.Errorf("%s get latest tx segments err, %v", h.name(), err)
		return
	}

	doHandleSegments(h.name(), 8, latestSegments, false, h.updateColl)
}

func (h *ModifyChainIdHandlerThree) updateColl(seg *segment, isTargetHistory bool) {
	collName := entity.CollectionNameExIbcTxLatest
	if isTargetHistory {
		collName = entity.CollectionNameExIbcTx
	}

	const (
		scChainIdField        = "sc_chain_id"
		dcChainIdField        = "dc_chain_id"
		scChainField          = "sc_chain"
		dcChainField          = "dc_chain"
		baseDenomChainIdField = "base_denom_chain_id"
		baseDenomChainField   = "base_denom_chain"
		createAtField         = "create_at"
		refundedTxInfoField   = "refunded_tx_info"
		ackTimeoutTxInfoField = "ack_timeout_tx_info"
	)

	singleScChainIdFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterSegmentBson(scChainIdField, chainId, createAtField, seg)
			sb := setBson(scChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}
		return nil
	}

	singleDcChainIdFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterSegmentBson(dcChainIdField, chainId, createAtField, seg)
			sb := setBson(dcChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}

		fb := filterSegmentBson(dcChainIdField, "", createAtField, seg)
		sb := setBson(dcChainField, "")
		if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
			return err
		}
		return nil
	}

	singleBaseDenomChainIdFunc := func() error {
		for chainId, chainName := range h.chainIdNameMap {
			if chainName == "" {
				logrus.Errorf("%s chain(%s) name is blank", h.name(), chainId)
				continue
			}
			fb := filterSegmentBson(baseDenomChainIdField, chainId, createAtField, seg)
			sb := setBson(baseDenomChainField, chainName)
			if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
				return err
			}
		}

		fb := filterSegmentBson(baseDenomChainIdField, "", createAtField, seg)
		sb := setBson(baseDenomChainField, "")
		if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
			return err
		}
		return nil
	}

	unsetFieldBatchFunc := func() error {
		fb := bson.M{
			createAtField: bson.M{
				"$gte": seg.StartTime,
				"$lte": seg.EndTime,
			},
		}
		sb := bson.M{
			"$unset": bson.M{
				baseDenomChainIdField: nil,
				scChainIdField:        nil,
				dcChainIdField:        nil,
			},
		}
		if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
			return err
		}
		return nil
	}

	// rename refunded_tx_info -> ack_timeout_tx_info
	renameFunc := func() error {
		fb := bson.M{
			createAtField: bson.M{
				"$gte": seg.StartTime,
				"$lte": seg.EndTime,
			},
		}
		sb := bson.M{
			"$rename": bson.M{
				refundedTxInfoField: ackTimeoutTxInfoField,
			},
		}
		if err := repository.CustomerUpdateAll(collName, context.Background(), fb, sb); err != nil {
			return err
		}
		return nil
	}

	logFunc := func(fn string, err error) {
		logrus.Errorf("%s import err, %s-%d-%d, %v", h.name(), fn, seg.StartTime, seg.EndTime, err)
	}

	if err := singleScChainIdFunc(); err != nil {
		logFunc("singleScChainIdFunc", err)
		return
	}
	if err := singleDcChainIdFunc(); err != nil {
		logFunc("singleDcChainIdFunc", err)
		return
	}
	if err := singleBaseDenomChainIdFunc(); err != nil {
		logFunc("singleBaseDenomChainIdFunc", err)
		return
	}
	if err := unsetFieldBatchFunc(); err != nil {
		logFunc("unsetFieldBatchFunc", err)
		return
	}
	if err := renameFunc(); err != nil {
		logFunc("renameFunc", err)
		return
	}
}
