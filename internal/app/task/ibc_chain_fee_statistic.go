package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type IBCChainFeeStatisticTask struct {
}

func (t *IBCChainFeeStatisticTask) Name() string {
	return "ibc_chain_fee_statistics_task"
}

func (t *IBCChainFeeStatisticTask) Cron() string {
	if taskConf.IBCChainFeeStatisticTask != "" {
		return taskConf.IBCChainFeeStatisticTask
	}
	return "0 */5 * * * ?"
}

func (t *IBCChainFeeStatisticTask) BeforeHook() error {
	return nil
}

var chainStatisticsCoordinator *stringQueueCoordinator

func (t *IBCChainFeeStatisticTask) Run() {
	err := t.yesterdayStatistics()
	if err != nil {
		TaskMetricMap.Store(t.Name(), float64(-1))
		return
	}
	err = t.todayStatistics()
	if err != nil {
		TaskMetricMap.Store(t.Name(), float64(-1))
		return
	}
	TaskMetricMap.Store(t.Name(), float64(1))

	return
}

// Run 全量更新
func (t *IBCChainFeeStatisticTask) RunAllChain() {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s Run getAllChainMap err, %v", t.Name(), err)
		return
	}
	// init coordinator
	chainQueue := new(utils.QueueString)
	for _, v := range chainMap {
		chainQueue.Push(v.ChainName)
	}
	chainStatisticsCoordinator = &stringQueueCoordinator{
		stringQueue: chainQueue,
	}

	if err = chainFeeStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s chainFeeStatisticsRepo.CreateNew err, %v", t.Name(), err)
		return
	}

	workerNum := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	for i := 1; i <= workerNum; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			defer waitGroup.Done()
			newChainStatisticsWorker(t.Name(), wn, chainMap).exec()
		}(workName)
	}
	waitGroup.Wait()

	if err = chainFeeStatisticsRepo.SwitchColl(); err != nil {
		logrus.Errorf("task %s chainFeeStatisticsRepo.SwitchColl() err, %v", t.Name(), err)
		return
	}

	return
}
func (t *IBCChainFeeStatisticTask) todayStatistics() error {
	logrus.Infof("task %s exec today statistics", t.Name())
	startTime, endTime := utils.TodayUnix()

	if err := t.RunIncrement(&segment{
		StartTime: startTime,
		EndTime:   endTime,
	}); err != nil {
		logrus.Errorf("task %s todayStatistics error, %v", t.Name(), err)
		return err
	}

	return nil
}

func (t *IBCChainFeeStatisticTask) yesterdayStatistics() error {
	startTime, endTime := utils.YesterdayUnix()

	logrus.Infof("task %s check yeaterday statistics", t.Name())
	if err := t.RunIncrement(&segment{
		StartTime: startTime,
		EndTime:   endTime,
	}); err != nil {
		logrus.Errorf("task %s yesterdayStatistics error, %v", t.Name(), err)
		return err
	}

	return nil
}

// RunIncrement 增量统计
func (t *IBCChainFeeStatisticTask) RunIncrement(seg *segment) error {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s IncrementRun getAllChainMap err, %v", t.Name(), err)
		return err
	}

	segs := []*segment{seg}
	worker := newChainStatisticsWorker(t.Name(), "increment", chainMap)
	for chain, _ := range chainMap {
		worker.statistics(chain, segs, opUpdate)
	}

	return nil
}

// RunWithParam 自定义统计
func (t *IBCChainFeeStatisticTask) RunWithParam(chain string, startTime, endTime int64) int {
	segments := segmentTool(segmentStep, startTime, endTime)
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s RunWithParam getAllChainMap err, %v", t.Name(), err)
		return -1
	}
	if _, ok := chainMap[chain]; !ok {
		logrus.Warnf("this chain[%s] no found in chain_config", chain)
		return 1
	}

	workerName := fmt.Sprintf("%s-%s", "cus", chain)
	if len(workerName) > 7 {
		workerName = workerName[:7]
	}
	worker := newChainStatisticsWorker(t.Name(), workerName, chainMap)
	worker.statistics(chain, segments, opUpdate)
	return 1
}

// =========================================================================
// =========================================================================
// worker

func newChainStatisticsWorker(taskName, workerName string, chainMap map[string]*entity.ChainConfig) *chainStatisticsWorker {
	return &chainStatisticsWorker{
		taskName:   taskName,
		workerName: workerName,
		chainMap:   chainMap,
	}
}

type chainStatisticsWorker struct {
	taskName   string
	workerName string
	chainMap   map[string]*entity.ChainConfig
}

func (w *chainStatisticsWorker) getChain() (string, error) {
	return chainStatisticsCoordinator.getOne()
}

// exec 全量统计
func (w *chainStatisticsWorker) exec() {
	logrus.Infof("task %s worker %s start", w.taskName, w.workerName)
	for {
		chain, err := w.getChain()
		if err != nil {
			logrus.Infof("task %s worker %s exit", w.taskName, w.workerName)
			break
		}

		if _, ok := w.chainMap[chain]; !ok {
			logrus.Warnf("task %s worker %s chain %s is no found in chainMap", w.taskName, w.workerName, chain)
			continue
		}

		logrus.Infof("task %s worker %s get chain: %v", w.taskName, w.workerName, chain)
		firstTx, err := txRepo.GetFirstTx(chain)
		if err != nil {
			logrus.Errorf("task %s worker %s chain %s GetFirstTx err, %v", w.taskName, w.workerName, chain, err)
			continue
		}

		segments := segmentTool(segmentStep, firstTx.Time, time.Now().Unix())
		w.statistics(chain, segments, opInsert)
	}

}

func (w *chainStatisticsWorker) statistics(chain string, segments []*segment, op int) {
	startTime := time.Now().Unix()
	logrus.Infof("task %s worker %s statistics chain: %s, total segments: %d", w.taskName, w.workerName, chain, len(segments))
	doHandleSegments(w.workerName, 5, segments, func(seg *segment) {
		// fee statistics
		var userFeeStats, totalFeeStats []*dto.ChainFeeStatisticsDTO
		gw := sync.WaitGroup{}
		gw.Add(2)
		go func() {
			defer gw.Done()
			var err error
			userFeeStats, err = txRepo.ChainUserFeeStatistics(chain, seg.StartTime, seg.EndTime)
			if err != nil {
				logrus.Errorf("task %s worker %s ChainUserFeeStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
			}
		}()
		go func() {
			defer gw.Done()
			var err error
			totalFeeStats, err = txRepo.ChainFeeStatistics(chain, seg.StartTime, seg.EndTime)
			if err != nil {
				logrus.Errorf("task %s worker %s ChainFeeStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
			}
		}()
		gw.Wait()
		feeStatList := make([]*entity.IBCChainFeeStatistics, 0, len(userFeeStats))
		userPayeMap := make(map[string]float64, len(userFeeStats))
		if len(userFeeStats) > 0 {
			for _, v := range userFeeStats {
				item := entity.IBCChainFeeStatistics{
					ChainName:        chain,
					TxStatus:         entity.TxStatus(v.Status),
					PayerType:        entity.UserPay,
					FeeDenom:         v.Denom,
					FeeAmount:        v.DenomAmount,
					SegmentStartTime: seg.StartTime,
					SegmentEndTime:   seg.EndTime,
					CreateAt:         time.Now().Unix(),
					UpdateAt:         time.Now().Unix(),
				}
				key := strconv.FormatInt(v.Status, 10) + v.Denom + strconv.FormatInt(seg.StartTime, 10) + strconv.FormatInt(seg.EndTime, 10)
				userPayeMap[key] = v.DenomAmount
				feeStatList = append(feeStatList, &item)
			}
		}

		if len(totalFeeStats) > 0 {
			for _, v := range totalFeeStats {
				key := strconv.FormatInt(v.Status, 10) + v.Denom + strconv.FormatInt(seg.StartTime, 10) + strconv.FormatInt(seg.EndTime, 10)
				item := entity.IBCChainFeeStatistics{
					ChainName:        chain,
					TxStatus:         entity.TxStatus(v.Status),
					PayerType:        entity.RelayerPay,
					FeeDenom:         v.Denom,
					FeeAmount:        v.DenomAmount - userPayeMap[key],
					SegmentStartTime: seg.StartTime,
					SegmentEndTime:   seg.EndTime,
					CreateAt:         time.Now().Unix(),
					UpdateAt:         time.Now().Unix(),
				}
				feeStatList = append(feeStatList, &item)
			}

		}
		if len(feeStatList) == 0 {
			return
		}
		err := w.saveFeeStat(chain, feeStatList, seg, op)
		if err != nil {
			logrus.Errorf("task %s worker %s saveFeeStat err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
		}

	})
	logrus.Infof("task %s worker %s statistics chain %s end,time use: %d(s)", w.taskName, w.workerName, chain, time.Now().Unix()-startTime)
}

func (w *chainStatisticsWorker) saveFeeStat(chain string, feeStatList []*entity.IBCChainFeeStatistics, segment *segment, op int) error {

	var err error
	if op == opInsert {
		if err = chainFeeStatisticsRepo.InsertManyToNew(feeStatList); err != nil {
			logrus.Errorf("task %s chainFeeStatisticsRepo.InsertManyToNew chain: %s err, %v", w.taskName, chain, err)
		}
	} else {
		if err = chainFeeStatisticsRepo.BatchSwap(chain, segment.StartTime, segment.EndTime, feeStatList); err != nil {
			logrus.Errorf("task %s chainFeeStatisticsRepo.BatchSwap chain: %s err, %v", w.taskName, chain, err)
		}
	}

	return err
}
