package task

import (
	"fmt"
	"sync"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/entity"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/utils"
	"github.com/sirupsen/logrus"
)

type FixDenomTraceDataTask struct {
	fixDenomTraceDataTrait
}

var _ OneOffTask = new(FixDenomTraceDataTask)

func (t *FixDenomTraceDataTask) Name() string {
	return "fix_denom_trace_data_task"
}

func (t *FixDenomTraceDataTask) Switch() bool {
	return global.Config.Task.SwitchFixDenomTraceDataTask
}

func (t *FixDenomTraceDataTask) Run() int {
	if !t.Switch() {
		logrus.Infof("task %s closed", t.Name())
		return 1
	}
	// init
	t.target = ibcTxTargetLatest
	t.startTime = global.Config.Task.FixDenomTraceDataStartTime
	t.endTime = global.Config.Task.FixDenomTraceDataEndTime

	if t.startTime < 0 || t.endTime < t.startTime {
		logrus.Errorf("task %s start/end time config error, start time: %d, end time: %d", t.Name(), t.startTime, t.endTime)
		return -1
	}

	// start worker
	if err := t.startWorker(t.Name()); err != nil {
		logrus.Errorf("task %s start worker failed, %v", t.Name(), err)
		return -1
	}

	return 1
}

// FixDenomTraceHistoryDataTask fix history tx
type FixDenomTraceHistoryDataTask struct {
	fixDenomTraceDataTrait
}

var _ OneOffTask = new(FixDenomTraceHistoryDataTask)

func (t *FixDenomTraceHistoryDataTask) Name() string {
	return "fix_denom_trace_history_data_task"
}

func (t *FixDenomTraceHistoryDataTask) Switch() bool {
	return global.Config.Task.SwitchFixDenomTraceHistoryDataTask
}

func (t *FixDenomTraceHistoryDataTask) Run() int {
	if !t.Switch() {
		logrus.Infof("task %s closed", t.Name())
		return 1
	}
	// init
	t.target = ibcTxTargetHistory
	t.startTime = global.Config.Task.FixDenomTraceHistoryDataStartTime
	t.endTime = global.Config.Task.FixDenomTraceHistoryDataEndTime

	if t.startTime < 0 || t.endTime < t.startTime {
		logrus.Errorf("task %s start/end time config error, start time: %d, end time: %d", t.Name(), t.startTime, t.endTime)
		return -1
	}

	// start worker
	if err := t.startWorker(t.Name()); err != nil {
		logrus.Errorf("task %s start worker failed, %v", t.Name(), err)
		return -1
	}

	return 1
}

// ================================================================================
// ================================================================================
// trait
type fixDenomTraceDataTrait struct {
	target    string
	startTime int64
	endTime   int64
}

func (trait *fixDenomTraceDataTrait) startWorker(taskName string) error {
	chainMap, err := getAllChainMap()
	if err != nil {
		logrus.Errorf("task %s getAllChainMap error, %v", taskName, err)
		return err
	}

	denomList, err := denomRepo.FindAll()
	if err != nil {
		logrus.Errorf("task %s denomRepo.FindAll error, %v", taskName, err)
		return err
	}

	denomMap := denomList.ConvertToMap()
	const step = 3600
	segments := trait.getSegment(trait.startTime, trait.endTime, step)
	var waitGroup sync.WaitGroup
	waitGroup.Add(fixDenomTraceDataTaskWorkerQuantity)
	for i := 0; i < fixDenomTraceDataTaskWorkerQuantity; i++ { // i must start from 0
		workName := fmt.Sprintf("worker-%d", i)
		workloadIndex := i
		go func(wn string, wi int) {
			newFixDenomTraceDataWorker(taskName, wn, trait.target, chainMap, denomMap, segments, wi).exec()
			waitGroup.Done()
		}(workName, workloadIndex)
	}
	waitGroup.Wait()

	return nil
}

func (trait *fixDenomTraceDataTrait) getSegment(startTime, endTime, step int64) []segment {
	var segments []segment
	for cursorTime := startTime; cursorTime <= endTime; cursorTime += step {
		if cursorTime+step > endTime {
			segments = append(segments, segment{
				StartTime: cursorTime,
				EndTime:   endTime,
			})
		} else {
			segments = append(segments, segment{
				StartTime: cursorTime,
				EndTime:   cursorTime + step - 1,
			})
		}
	}

	return segments
}

// ================================================================================
// ================================================================================
// worker
func newFixDenomTraceDataWorker(taskName string, workerName string, target string, chainMap map[string]*entity.ChainConfig,
	denomMap map[string]*entity.IBCDenom, workloadSegments []segment, workloadIndex int) *fixDenomTraceDataWorker {
	return &fixDenomTraceDataWorker{
		taskName:         taskName,
		workerName:       workerName,
		target:           target,
		chainMap:         chainMap,
		denomMap:         denomMap,
		workloadSegments: workloadSegments,
		workloadIndex:    workloadIndex,
	}
}

type fixDenomTraceDataWorker struct {
	taskName         string
	workerName       string
	target           string
	chainMap         map[string]*entity.ChainConfig
	denomMap         map[string]*entity.IBCDenom
	newDenomMap      map[string]*entity.IBCDenom
	workloadSegments []segment
	workloadIndex    int
}

func (w *fixDenomTraceDataWorker) exec() {
	logrus.Infof("task %s worker %s start", w.taskName, w.workerName)
	w.newDenomMap = make(map[string]*entity.IBCDenom)
	workloadAmount := 0
	for index, seg := range w.workloadSegments {
		if index%fixDenomTraceDataTaskWorkerQuantity != w.workloadIndex {
			continue
		}

		workloadAmount++
		w.fixData(seg.StartTime, seg.EndTime)
	}

	logrus.Infof("task %s worker %s end, workloadAmount: %d", w.taskName, w.workerName, workloadAmount)
}

func (w *fixDenomTraceDataWorker) fixData(startTime, endTime int64) {
	logrus.Debugf("task %s worker %s start fix data, setgment: %d, %d", w.taskName, w.workerName, startTime, endTime)
	const limit = 1000
	var skip int64 = 0
	var newDenomList entity.IBCDenomList

	for {
		txList, err := w.getTxs(startTime, endTime, skip, limit)
		if err != nil {
			logrus.Errorf("task %s worker %s fix data error, setgment: %d, %d", w.taskName, w.workerName, startTime, endTime)
			return
		}

		for _, tx := range txList {
			scDenom, dcDenom := w.parseDenom(tx)
			if scDenom != nil && w.newDenomMap[fmt.Sprintf("%s%s", scDenom.ChainId, scDenom.Denom)] == nil {
				newDenomList = append(newDenomList, scDenom)
				w.newDenomMap[fmt.Sprintf("%s%s", scDenom.ChainId, scDenom.Denom)] = scDenom
			}
			if dcDenom != nil && w.newDenomMap[fmt.Sprintf("%s%s", dcDenom.ChainId, dcDenom.Denom)] == nil {
				newDenomList = append(newDenomList, dcDenom)
				w.newDenomMap[fmt.Sprintf("%s%s", dcDenom.ChainId, dcDenom.Denom)] = dcDenom
			}

			// update ibc tx
			if scDenom != nil {
				tx.BaseDenom = scDenom.BaseDenom
				tx.BaseDenomChainId = scDenom.BaseDenomChainId
			}
			if tx.CreateAt < fixCreateAtErrTime {
				tx.CreateAt = tx.TxTime
				tx.UpdateAt = tx.TxTime
			}

			originRecordId := tx.RecordId
			tx.RecordId = utils.Md5(tx.RecordId)
			if err = w.UpdateTx(originRecordId, tx); err != nil {
				logrus.Errorf("task %s worker %s update tx(%s) error, %v", w.taskName, w.workerName, tx.ScTxInfo.Hash, err)
			}
		}

		if len(txList) < limit {
			break
		}
		skip += limit
	}

	if len(newDenomList) > 0 {
		if err := denomRepo.InsertBatchToNew(newDenomList); err != nil {
			logrus.Errorf("task %s worker %s insert new denoms error, %v", w.taskName, w.workerName, err)
		}
	}
}

func (w *fixDenomTraceDataWorker) getTxs(startTime, endTime, skip, limit int64) ([]*entity.ExIbcTx, error) {
	if w.target == ibcTxTargetHistory {
		return ibcTxRepo.FindHistoryByTxTime(startTime, endTime, skip, limit)
	}

	return ibcTxRepo.FindByTxTime(startTime, endTime, skip, limit)
}

func (w *fixDenomTraceDataWorker) UpdateTx(originRecordId string, tx *entity.ExIbcTx) error {
	if w.target == ibcTxTargetHistory {
		return ibcTxRepo.UpdateDenomTraceHistory(originRecordId, tx)
	}

	return ibcTxRepo.UpdateDenomTrace(originRecordId, tx)
}

func (w *fixDenomTraceDataWorker) parseDenom(tx *entity.ExIbcTx) (*entity.IBCDenom, *entity.IBCDenom) {
	// parse sc denom
	var scDenomEntityNew, dcDenomEntityNew *entity.IBCDenom
	scDenom := tx.Denoms.ScDenom
	scDenomEntity, ok := w.denomMap[fmt.Sprintf("%s%s", tx.ScChainId, scDenom)]
	if ok {
		var scDenomFullPath string
		if scDenomEntity.DenomPath != "" {
			scDenomFullPath = fmt.Sprintf("%s/%s", scDenomEntity.DenomPath, scDenomEntity.BaseDenom)
		} else {
			scDenomFullPath = scDenomEntity.BaseDenom
		}
		scDenomEntityNew = traceDenom(scDenomFullPath, tx.ScChainId, w.chainMap)
	}

	// parse dc denom
	if scDenomEntityNew == nil {
		_ = storageCache.AddMissDenom(utils.Md5(tx.RecordId), tx.ScChainId, scDenom)
		return nil, nil
	} else {
		scDenomEntityNew.CreateAt = scDenomEntity.CreateAt
		scDenomEntityNew.UpdateAt = scDenomEntity.UpdateAt
	}

	if tx.Status != entity.IbcTxStatusSuccess || tx.DcTxInfo == nil || tx.DcTxInfo.Msg == nil {
		return scDenomEntityNew, nil
	}

	_, isCrossBack := calculateNextDenomPath(tx.DcTxInfo.Msg.RecvPacketMsg().Packet)
	if isCrossBack { // transfer to origin chain
		return scDenomEntityNew, nil
	}

	dcDenom := tx.Denoms.DcDenom
	DcDenomEntity, ok := w.denomMap[fmt.Sprintf("%s%s", tx.DcChainId, dcDenom)]
	if ok {
		dcDenomEntityNew = &entity.IBCDenom{
			Symbol:           "",
			ChainId:          tx.DcChainId,
			Denom:            dcDenom,
			PrevDenom:        tx.Denoms.ScDenom,
			PrevChainId:      tx.ScChainId,
			BaseDenom:        scDenomEntityNew.BaseDenom,
			BaseDenomChainId: scDenomEntityNew.BaseDenomChainId,
			DenomPath:        DcDenomEntity.DenomPath,
			RootDenom:        scDenomEntityNew.RootDenom,
			IsBaseDenom:      false,
			CreateAt:         DcDenomEntity.CreateAt,
			UpdateAt:         DcDenomEntity.UpdateAt,
		}
	} else {
		_ = storageCache.AddMissDenom(tx.RecordId, tx.ScChainId, scDenom)
		return scDenomEntityNew, nil
	}

	return scDenomEntityNew, dcDenomEntityNew
}
