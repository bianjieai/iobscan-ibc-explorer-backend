package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type RelayerStatisticsTask struct {
}

var relayerStatisticsCoordinator *stringQueueCoordinator

func (t *RelayerStatisticsTask) Name() string {
	return "ibc_relayer_statistics_task"
}

func (t *RelayerStatisticsTask) Switch() bool {
	return global.Config.Task.SwitchIbcRelayerStatisticsTask
}

// Run 全量更新
func (t *RelayerStatisticsTask) Run() int {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap err, %v", t.Name(), err)
		return -1
	}

	// init coordinator
	chainQueue := new(utils.QueueString)
	for _, v := range chainMap {
		chainQueue.Push(v.ChainName)
	}
	relayerStatisticsCoordinator = &stringQueueCoordinator{
		stringQueue: chainQueue,
	}

	if err = relayerDenomStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s relayerDenomStatisticsRepo.CreateNew err, %v", t.Name(), err)
		return -1
	}

	if err = relayerFeeStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s relayerFeeStatisticsRepo.CreateNew err, %v", t.Name(), err)
		return -1
	}

	workerNum := len(chainMap)
	if workerNum > relayerStatisticsWorkerNum {
		workerNum = relayerStatisticsWorkerNum
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	for i := 1; i <= workerNum; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			defer waitGroup.Done()
			newRelayerStatisticsWorker(t.Name(), wn, chainMap).exec()
		}(workName)
	}
	waitGroup.Wait()

	if err = relayerDenomStatisticsRepo.SwitchColl(); err != nil {
		logrus.Errorf("task %s relayerDenomStatisticsRepo.SwitchColl() err, %v", t.Name(), err)
		return -1
	}

	if err = relayerFeeStatisticsRepo.SwitchColl(); err != nil {
		logrus.Errorf("task %s relayerFeeStatisticsRepo.SwitchColl() err, %v", t.Name(), err)
		return -1
	}

	t.flushCache()
	return 1
}

// RunIncrement 增量统计
func (t *RelayerStatisticsTask) RunIncrement(seg *segment) error {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s IncrementRun getAllChainMap err, %v", t.Name(), err)
		return err
	}

	segs := []*segment{seg}
	worker := newRelayerStatisticsWorker(t.Name(), "increment", chainMap)
	for chain, _ := range chainMap {
		worker.statistics(chain, segs, opUpdate)
	}
	t.flushCache()
	return nil
}

// RunWithParam 自定义统计
func (t *RelayerStatisticsTask) RunWithParam(chain string, startTime, endTime int64) int {
	segments := segmentTool(segmentStepLatest, startTime, endTime)
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s RunWithParam getAllChainMap err, %v", t.Name(), err)
		return -1
	}

	workerName := fmt.Sprintf("%s-%s", "cus", chain)
	if len(workerName) > 7 {
		workerName = workerName[:7]
	}
	worker := newRelayerStatisticsWorker(t.Name(), workerName, chainMap)
	worker.statistics(chain, segments, opUpdate)
	t.flushCache()
	return 1
}

// flushCache 清除relayer相关缓存
func (t *RelayerStatisticsTask) flushCache() {
	//_ = relayerDataCache.DelTotalFeeCost()
	//_ = relayerDataCache.DelTotalRelayedValue()
	_ = relayerDataCache.DelTransferTypeTxs()
	_ = relayerDataCache.DelTotalTxs()
	_ = relayerDataCache.DelRelayedTrend()
}

// =========================================================================
// =========================================================================
// worker

func newRelayerStatisticsWorker(taskName, workerName string, chainMap map[string]*entity.ChainConfig) *relayerStatisticsWorker {
	return &relayerStatisticsWorker{
		taskName:   taskName,
		workerName: workerName,
		chainMap:   chainMap,
	}
}

type relayerStatisticsWorker struct {
	taskName   string
	workerName string
	chainMap   map[string]*entity.ChainConfig
}

func (w *relayerStatisticsWorker) getChain() (string, error) {
	return relayerStatisticsCoordinator.getOne()
}

// exec 全量统计
func (w *relayerStatisticsWorker) exec() {
	logrus.Infof("task %s worker %s start", w.taskName, w.workerName)
	for {
		chain, err := w.getChain()
		if err != nil {
			logrus.Infof("task %s worker %s exit", w.taskName, w.workerName)
			break
		}

		if cf, ok := w.chainMap[chain]; ok && cf.Status == entity.ChainStatusClosed {
			logrus.Infof("task %s worker %s chain %s is closed", w.taskName, w.workerName, chain)
			continue
		}

		logrus.Infof("task %s worker %s get chain: %v", w.taskName, w.workerName, chain)
		firstTx, err := txRepo.GetFirstTx(chain)
		if err != nil {
			logrus.Errorf("task %s worker %s chain %s GetFirstTx err, %v", w.taskName, w.workerName, chain, err)
			continue
		}

		segments := segmentTool(segmentStepLatest, firstTx.Time, time.Now().Unix())
		w.statistics(chain, segments, opInsert)
	}
}

func (w *relayerStatisticsWorker) statistics(chain string, segments []*segment, op int) {
	startTime := time.Now().Unix()
	logrus.Infof("task %s worker %s statistics chain: %s, total segments: %d", w.taskName, w.workerName, chain, len(segments))

	for _, v := range segments {
		// denom statistics
		denomStats, err := txRepo.RelayerDenomStatistics(chain, v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s worker %s RelayerDenomStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, v.StartTime, v.EndTime, err)
		} else {
			denomStatMap, addrChannelMap := w.aggrDenomStat(chain, v, denomStats)
			_ = w.saveDenomStat(chain, denomStatMap, v, op)
			_ = w.saveAddrChannel(addrChannelMap)
		}

		// fee statistics
		feeStats, err := txRepo.RelayerFeeStatistics(chain, v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s worker %s RelayerFeeStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, v.StartTime, v.EndTime, err)
		} else {
			_ = w.saveFeeStat(chain, feeStats, v, op)
		}
	}

	logrus.Infof("task %s worker %s statistics chain %s end,time use: %d(s)", w.taskName, w.workerName, chain, time.Now().Unix()-startTime)
}

func (w *relayerStatisticsWorker) aggrDenomStat(chain string, segment *segment, stats []*dto.RelayerDenomStatisticsDTO) (map[string]*entity.IBCRelayerDenomStatistics, map[string]*entity.IBCRelayerAddressChannel) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("task %s aggrDenomStatistics err, %s-%d-%d, %v", w.taskName, chain, segment.StartTime, segment.EndTime, r)
		}
	}()

	// aggr denom statistics
	denomStatMap := make(map[string]*entity.IBCRelayerDenomStatistics)
	for _, v := range stats {
		var denomChain string
		if v.TxType == string(entity.TxTypeAckPacket) || v.TxType == string(entity.TxTypeTimeoutPacket) {
			denomChain = chain
		} else {
			denomChain = w.chainMap[chain].GetDcChain(v.DcChannel, v.ScChannel)
		}

		denomEntity := traceDenom(v.Denom, denomChain, w.chainMap)
		dsmk := fmt.Sprintf("%s%s%d%s%s", v.Signer, v.TxType, v.Status, denomEntity.BaseDenom, denomEntity.BaseDenomChain)
		if _, ok := denomStatMap[dsmk]; ok {
			denomStatMap[dsmk].RelayedAmount += v.DenomAmount
			denomStatMap[dsmk].RelayedTxs += v.TxsCount
		} else {
			denomStatMap[dsmk] = &entity.IBCRelayerDenomStatistics{
				StatisticChain:   chain,
				RelayerAddress:   v.Signer,
				ChainAddressComb: entity.GenerateChainAddressComb(chain, v.Signer),
				TxStatus:         entity.TxStatus(v.Status),
				TxType:           entity.TxType(v.TxType),
				BaseDenom:        denomEntity.BaseDenom,
				BaseDenomChain:   denomEntity.BaseDenomChain,
				RelayedAmount:    v.DenomAmount,
				RelayedTxs:       v.TxsCount,
				SegmentStartTime: segment.StartTime,
				SegmentEndTime:   segment.EndTime,
				CreateAt:         time.Now().Unix(),
				UpdateAt:         time.Now().Unix(),
			}
		}
	}

	// aggr relayer address channel
	addrChannelMap := make(map[string]*entity.IBCRelayerAddressChannel)
	for _, v := range stats {
		addrChannel := entity.IBCRelayerAddressChannel{
			RelayerAddress:      v.Signer,
			Channel:             "",
			Chain:               chain,
			CounterPartyChannel: "",
			CreateAt:            time.Now().Unix(),
			UpdateAt:            time.Now().Unix(),
		}
		if v.TxType == string(entity.TxTypeAckPacket) || v.TxType == string(entity.TxTypeTimeoutPacket) {
			addrChannel.Channel = v.ScChannel
			addrChannel.CounterPartyChannel = v.DcChannel
		} else {
			addrChannel.Channel = v.DcChannel
			addrChannel.CounterPartyChannel = v.ScChannel
		}

		acmk := fmt.Sprintf("%s%s%s", addrChannel.RelayerAddress, addrChannel.Chain, addrChannel.Channel)
		if _, ok := addrChannelMap[acmk]; !ok {
			addrChannelMap[acmk] = &addrChannel
		}
	}

	return denomStatMap, addrChannelMap
}

func (w *relayerStatisticsWorker) saveDenomStat(chain string, denomStatMap map[string]*entity.IBCRelayerDenomStatistics, segment *segment, op int) error {
	if len(denomStatMap) == 0 {
		return nil
	}
	denomStats := make([]*entity.IBCRelayerDenomStatistics, 0, len(denomStatMap))
	for _, v := range denomStatMap {
		denomStats = append(denomStats, v)
	}

	var err error
	if op == opInsert {
		if err = relayerDenomStatisticsRepo.InsertManyToNew(denomStats); err != nil {
			logrus.Errorf("task %s relayerDenomStatisticsRepo.InsertManyToNew chain: %s err, %v", w.taskName, chain, err)
		}
	} else {
		if err = relayerDenomStatisticsRepo.BatchSwap(chain, segment.StartTime, segment.EndTime, denomStats); err != nil {
			logrus.Errorf("task %s relayerDenomStatisticsRepo.BatchSwap chain: %s, err, %v", w.taskName, chain, err)
		}
	}

	return err
}

func (w *relayerStatisticsWorker) saveFeeStat(chain string, feeStats []*dto.RelayerFeeStatisticsDTO, segment *segment, op int) error {
	if len(feeStats) == 0 {
		return nil
	}
	feeStatList := make([]*entity.IBCRelayerFeeStatistics, 0, len(feeStats))
	for _, v := range feeStats {
		feeStatList = append(feeStatList, &entity.IBCRelayerFeeStatistics{
			StatisticChain:   chain,
			RelayerAddress:   v.Signer,
			ChainAddressComb: entity.GenerateChainAddressComb(chain, v.Signer),
			TxStatus:         entity.TxStatus(v.Status),
			TxType:           entity.TxType(v.TxType),
			FeeDenom:         v.Denom,
			FeeAmount:        v.DenomAmount,
			RelayedTxs:       v.TxsCount,
			SegmentStartTime: segment.StartTime,
			SegmentEndTime:   segment.EndTime,
			CreateAt:         time.Now().Unix(),
			UpdateAt:         time.Now().Unix(),
		})
	}

	var err error
	if op == opInsert {
		if err = relayerFeeStatisticsRepo.InsertManyToNew(feeStatList); err != nil {
			logrus.Errorf("task %s relayerFeeStatisticsRepo.InsertManyToNew chain: %s err, %v", w.taskName, chain, err)
		}
	} else {
		if err = relayerFeeStatisticsRepo.BatchSwap(chain, segment.StartTime, segment.EndTime, feeStatList); err != nil {
			logrus.Errorf("task %s relayerFeeStatisticsRepo.BatchSwap chain: %s err, %v", w.taskName, chain, err)
		}
	}

	return err
}

func (w *relayerStatisticsWorker) saveAddrChannel(addrChannelMap map[string]*entity.IBCRelayerAddressChannel) error {
	if len(addrChannelMap) == 0 {
		return nil
	}

	addrChanList := make([]*entity.IBCRelayerAddressChannel, 0, len(addrChannelMap))
	for _, v := range addrChannelMap {
		addrChanList = append(addrChanList, v)
	}

	if err := relayerAddressChannelRepo.InsertMany(addrChanList); err != nil {
		logrus.Errorf("task %s relayerAddressChannelRepo.InsertMany err, %v", w.taskName, err)
		return err
	}

	return nil
}

// ===============================================================
// ===============================================================
// ===============================================================
// fix task

type FixRelayerStatisticsTask struct {
}

func (t *FixRelayerStatisticsTask) Name() string {
	return "fix_relayer_statistics_task"
}

func (t *FixRelayerStatisticsTask) Switch() bool {
	return false
}

// Run 全量更新
func (t *FixRelayerStatisticsTask) Run() int {
	st := time.Now().Unix()
	chainAddressList, err := relayerDenomStatisticsRepo.AggrChainAddressPair()
	if err != nil {
		logrus.Errorf("%s relayerDenomStatisticsRepo.AggrChainAddressPair err, %v", t.Name(), err)
		return -1
	}

	for _, v := range chainAddressList {
		chainAddressComb := entity.GenerateChainAddressComb(v.Chain, v.Address)
		if err = relayerDenomStatisticsRepo.UpdateChainAddressComb(v.Chain, v.Address, chainAddressComb); err != nil {
			logrus.Errorf("%s relayerDenomStatisticsRepo.UpdateChainAddressComb %s-%s-%s err, %v", t.Name(), v.Chain, v.Address, chainAddressComb, err)
		}
	}
	logrus.Infof("%s fix denom end, time use: %d[s]", t.Name(), time.Now().Unix()-st)

	st = time.Now().Unix()
	chainAddressList, err = relayerFeeStatisticsRepo.AggrChainAddressPair()
	if err != nil {
		logrus.Errorf("%s relayerFeeStatisticsRepo.AggrChainAddressPair err, %v", t.Name(), err)
		return -1
	}

	for _, v := range chainAddressList {
		chainAddressComb := entity.GenerateChainAddressComb(v.Chain, v.Address)
		if err = relayerFeeStatisticsRepo.UpdateChainAddressComb(v.Chain, v.Address, chainAddressComb); err != nil {
			logrus.Errorf("%s relayerFeeStatisticsRepo.UpdateChainAddressComb %s-%s-%s err, %v", t.Name(), v.Chain, v.Address, chainAddressComb, err)
		}
	}
	logrus.Infof("%s fix fee end, time use: %d[s]", t.Name(), time.Now().Unix()-st)

	return 1
}
