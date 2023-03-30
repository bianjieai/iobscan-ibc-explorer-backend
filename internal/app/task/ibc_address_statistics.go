package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/dto"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type IBCAddressStatisticTask struct {
}

func (t *IBCAddressStatisticTask) Name() string {
	return "ibc_address_statistics_task"
}

func (t *IBCAddressStatisticTask) Cron() string {
	if taskConf.IBCAddressStatisticTask != "" {
		return taskConf.IBCAddressStatisticTask
	}
	return "0 */5 * * * ?"
}

func (t *IBCAddressStatisticTask) BeforeHook() error {
	return nil
}

var addressStatisticsCoordinator *stringQueueCoordinator

func (t *IBCAddressStatisticTask) Run() {
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
func (t *IBCAddressStatisticTask) RunAllChain() {
	chainMap, err := GetAllChainMap()
	if err != nil {
		logrus.Errorf("task %s Run getAllChainMap err, %v", t.Name(), err)
		return
	}
	// init coordinator
	chainQueue := new(utils.QueueString)
	for _, v := range chainMap {
		chainQueue.Push(v.ChainName)
	}
	addressStatisticsCoordinator = &stringQueueCoordinator{
		stringQueue: chainQueue,
	}

	if err = addressStatisticsRepo.CreateNew(); err != nil {
		logrus.Errorf("task %s addressStatisticsRepo.CreateNew err, %v", t.Name(), err)
		return
	}

	workerNum := 3
	var waitGroup sync.WaitGroup
	waitGroup.Add(workerNum)
	for i := 1; i <= workerNum; i++ {
		workName := fmt.Sprintf("worker-%d", i)
		go func(wn string) {
			defer waitGroup.Done()
			newAddressStatisticsWorker(t.Name(), wn, chainMap).exec()
		}(workName)
	}
	waitGroup.Wait()

	if err = addressStatisticsRepo.SwitchColl(); err != nil {
		logrus.Errorf("task %s addressStatisticsRepo.SwitchColl() err, %v", t.Name(), err)
		return
	}

	return
}
func (t *IBCAddressStatisticTask) todayStatistics() error {
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

func (t *IBCAddressStatisticTask) yesterdayStatistics() error {
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
func (t *IBCAddressStatisticTask) RunIncrement(seg *segment) error {
	chainMap, err := GetAllChainMap()
	if err != nil {
		logrus.Errorf("task %s IncrementRun getAllChainMap err, %v", t.Name(), err)
		return err
	}

	segs := []*segment{seg}
	worker := newAddressStatisticsWorker(t.Name(), "increment", chainMap)
	for chain, _ := range chainMap {
		worker.statistics(chain, segs, opUpdate)
	}

	return nil
}

// RunWithParam 自定义统计
func (t *IBCAddressStatisticTask) RunWithParam(chain string, startTime, endTime int64) int {
	segments := segmentTool(segmentStep, startTime, endTime)
	chainMap, err := GetAllChainMap()
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
	worker := newAddressStatisticsWorker(t.Name(), workerName, chainMap)
	worker.statistics(chain, segments, opUpdate)
	return 1
}

// =========================================================================
// =========================================================================
// worker

func newAddressStatisticsWorker(taskName, workerName string, chainMap map[string]*entity.ChainConfig) *addressStatisticsWorker {
	return &addressStatisticsWorker{
		taskName:   taskName,
		workerName: workerName,
		chainMap:   chainMap,
	}
}

type addressStatisticsWorker struct {
	taskName   string
	workerName string
	chainMap   map[string]*entity.ChainConfig
}

func (w *addressStatisticsWorker) getChain() (string, error) {
	return addressStatisticsCoordinator.getOne()
}

// exec 全量统计
func (w *addressStatisticsWorker) exec() {
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

func (w *addressStatisticsWorker) statistics(chain string, segments []*segment, op int) {
	startTime := time.Now().Unix()
	logrus.Infof("task %s worker %s statistics chain: %s, total segments: %d", w.taskName, w.workerName, chain, len(segments))
	doHandleSegments(w.workerName, 5, segments, func(seg *segment) {
		var dcHistoryAddressStats, dcNewAddressStats, scHistoryAddressStats, scNewAddressStats []*dto.ChainActiveAddressesDTO
		gw := sync.WaitGroup{}
		gw.Add(4)
		go func() {
			defer gw.Done()
			var err error
			dcHistoryAddressStats, err = ibcTxRepo.AggrChainAddress(seg.StartTime, seg.EndTime, true, true)
			if err != nil {
				logrus.Errorf("task %s worker %s dcHistoryAddressStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
			}
		}()
		go func() {
			defer gw.Done()
			var err error
			dcNewAddressStats, err = ibcTxRepo.AggrChainAddress(seg.StartTime, seg.EndTime, false, true)
			if err != nil {
				logrus.Errorf("task %s worker %s dcNewAddressStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
			}
		}()
		go func() {
			defer gw.Done()
			var err error
			scHistoryAddressStats, err = ibcTxRepo.AggrChainAddress(seg.StartTime, seg.EndTime, true, false)
			if err != nil {
				logrus.Errorf("task %s worker %s scHistoryAddressStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
			}
		}()
		go func() {
			defer gw.Done()
			var err error
			scNewAddressStats, err = ibcTxRepo.AggrChainAddress(seg.StartTime, seg.EndTime, false, false)
			if err != nil {
				logrus.Errorf("task %s worker %s scNewAddressStatistics err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
			}
		}()
		gw.Wait()

		chainAddrsMap := make(map[string]utils.StringSet)
		addToSet(chainAddrsMap, dcHistoryAddressStats)
		addToSet(chainAddrsMap, dcNewAddressStats)
		addToSet(chainAddrsMap, scHistoryAddressStats)
		addToSet(chainAddrsMap, scNewAddressStats)

		var addressStatLlist []*entity.IBCAddressStatistics
		for chainName, addrs := range chainAddrsMap {
			if chain != "" {
				if chainName == chain {
					item := entity.IBCAddressStatistics{
						ChainName:        chain,
						ActiveAddressNum: int64(len(addrs)),
						SegmentStartTime: seg.StartTime,
						SegmentEndTime:   seg.EndTime,
						CreateAt:         time.Now().Unix(),
						UpdateAt:         time.Now().Unix(),
					}
					addressStatLlist = append(addressStatLlist, &item)
					break
				}
			} else {
				item := entity.IBCAddressStatistics{
					ChainName:        chain,
					ActiveAddressNum: int64(len(addrs)),
					SegmentStartTime: seg.StartTime,
					SegmentEndTime:   seg.EndTime,
					CreateAt:         time.Now().Unix(),
					UpdateAt:         time.Now().Unix(),
				}
				addressStatLlist = append(addressStatLlist, &item)
			}
		}

		err := w.saveAddressStat(chain, addressStatLlist, seg, op)
		if err != nil {
			fmt.Println(err.Error())
			logrus.Errorf("task %s worker %s saveAddressStat err, %s-%d-%d, %v", w.taskName, w.workerName, chain, seg.StartTime, seg.EndTime, err)
		}

	})
	logrus.Infof("task %s worker %s statistics chain %s end,time use: %d(s)", w.taskName, w.workerName, chain, time.Now().Unix()-startTime)
}

func addToSet(chainAddrsMap map[string]utils.StringSet, stats []*dto.ChainActiveAddressesDTO) {
	for _, v := range stats {
		if _, ok := chainAddrsMap[v.Chain]; ok {
			chainAddrsMap[v.Chain].AddAll(v.Addresses...)
		} else {
			set := utils.NewStringSetFromStr(v.Addresses...)
			chainAddrsMap[v.Chain] = set
		}
	}
}

func (w *addressStatisticsWorker) saveAddressStat(chain string, addressStatList []*entity.IBCAddressStatistics, segment *segment, op int) error {

	var err error
	if op == opInsert {
		if err = addressStatisticsRepo.InsertManyToNew(addressStatList); err != nil {
			logrus.Errorf("task %s addressStatisticsRepo.InsertManyToNew chain: %s err, %v", w.taskName, chain, err)
		}
	} else {
		if err = addressStatisticsRepo.BatchSwap(chain, segment.StartTime, segment.EndTime, addressStatList); err != nil {
			logrus.Errorf("task %s addressStatisticsRepo.BatchSwap chain: %s err, %v", w.taskName, chain, err)
		}
	}

	return err
}
