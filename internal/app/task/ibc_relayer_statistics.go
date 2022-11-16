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
	for chainId, _ := range chainMap {
		worker.statistics(chainId, segs, opUpdate)
	}
	t.flushCache()
	return nil
}

// RunWithParam 自定义统计
func (t *RelayerStatisticsTask) RunWithParam(chainId string, startTime, endTime int64) int {
	segments := segmentTool(segmentStepLatest, startTime, endTime)
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s RunWithParam getAllChainMap err, %v", t.Name(), err)
		return -1
	}

	workerName := fmt.Sprintf("%s-%s", "cus", chainId)
	if len(workerName) > 7 {
		workerName = workerName[:7]
	}
	worker := newRelayerStatisticsWorker(t.Name(), workerName, chainMap)
	worker.statistics(chainId, segments, opUpdate)
	t.flushCache()
	return 1
}

// flushCache 清除relayer相关缓存
func (t *RelayerStatisticsTask) flushCache() {
	_ = relayerDataCache.DelTotalFeeCost()
	_ = relayerDataCache.DelTotalRelayedValue()
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

func (w *relayerStatisticsWorker) statistics(chainId string, segments []*segment, op int) {
	startTime := time.Now().Unix()
	logrus.Infof("task %s worker %s statistics chain: %s, total segments: %d", w.taskName, w.workerName, chainId, len(segments))

	for _, v := range segments {
		// denom statistics
		denomStats, err := txRepo.RelayerDenomStatistics(chainId, v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s worker %s RelayerDenomStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chainId, v.StartTime, v.EndTime, err)
		} else {
			denomStatMap, addrChannelMap := w.aggrDenomStat(chainId, v, denomStats)
			_ = w.saveDenomStat(chainId, denomStatMap, v, op)
			_ = w.saveAddrChannel(addrChannelMap)
		}

		// fee statistics
		feeStats, err := txRepo.RelayerFeeStatistics(chainId, v.StartTime, v.EndTime)
		if err != nil {
			logrus.Errorf("task %s worker %s RelayerFeeStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chainId, v.StartTime, v.EndTime, err)
		} else {
			_ = w.saveFeeStat(chainId, feeStats, v, op)
		}
	}

	logrus.Infof("task %s worker %s statistics chain %s end,time use: %d(s)", w.taskName, w.workerName, chainId, time.Now().Unix()-startTime)
}

func (w *relayerStatisticsWorker) aggrDenomStat(chainId string, segment *segment, stats []*dto.RelayerDenomStatisticsDTO) (map[string]*entity.IBCRelayerDenomStatistics, map[string]*entity.IBCRelayerAddressChannel) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("task %s aggrDenomStatistics err, %s-%d-%d, %v", w.taskName, chainId, segment.StartTime, segment.EndTime, r)
		}
	}()

	// aggr denom statistics
	denomStatMap := make(map[string]*entity.IBCRelayerDenomStatistics)
	for _, v := range stats {
		var denomChain string
		if v.TxType == string(entity.TxTypeAckPacket) || v.TxType == string(entity.TxTypeTimeoutPacket) {
			denomChain = chainId
		} else {
			denomChain = w.chainMap[chainId].GetDcChainId(v.DcChannel, v.ScChannel)
		}

		denomEntity := traceDenom(v.Denom, denomChain, w.chainMap)
		dsmk := fmt.Sprintf("%s%s%d%s%s", v.Signer, v.TxType, v.Status, denomEntity.BaseDenom, denomEntity.BaseDenomChain)
		if _, ok := denomStatMap[dsmk]; ok {
			denomStatMap[dsmk].RelayedAmount += v.DenomAmount
			denomStatMap[dsmk].RelayedTxs += v.TxsCount
		} else {
			denomStatMap[dsmk] = &entity.IBCRelayerDenomStatistics{
				StatisticChain:   chainId,
				RelayerAddress:   v.Signer,
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
			Chain:               chainId,
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

func (w *relayerStatisticsWorker) saveDenomStat(chainId string, denomStatMap map[string]*entity.IBCRelayerDenomStatistics, segment *segment, op int) error {
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
			logrus.Errorf("task %s relayerDenomStatisticsRepo.InsertManyToNew chain: %s err, %v", w.taskName, chainId, err)
		}
	} else {
		if err = relayerDenomStatisticsRepo.BatchSwap(chainId, segment.StartTime, segment.EndTime, denomStats); err != nil {
			logrus.Errorf("task %s relayerDenomStatisticsRepo.BatchSwap chain: %s, err, %v", w.taskName, chainId, err)
		}
	}

	return err
}

func (w *relayerStatisticsWorker) saveFeeStat(chainId string, feeStats []*dto.RelayerFeeStatisticsDTO, segment *segment, op int) error {
	if len(feeStats) == 0 {
		return nil
	}
	feeStatList := make([]*entity.IBCRelayerFeeStatistics, 0, len(feeStats))
	for _, v := range feeStats {
		feeStatList = append(feeStatList, &entity.IBCRelayerFeeStatistics{
			StatisticChain:   chainId,
			RelayerAddress:   v.Signer,
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
			logrus.Errorf("task %s relayerFeeStatisticsRepo.InsertManyToNew chain: %s err, %v", w.taskName, chainId, err)
		}
	} else {
		if err = relayerFeeStatisticsRepo.BatchSwap(chainId, segment.StartTime, segment.EndTime, feeStatList); err != nil {
			logrus.Errorf("task %s relayerFeeStatisticsRepo.BatchSwap chain: %s err, %v", w.taskName, chainId, err)
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
